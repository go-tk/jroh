package apicommon

import (
	"bytes"
	"net/http"
)

type ServerOptions struct {
	Middlewares      map[MethodIndex][]ServerMiddleware
	RPCFilters       map[MethodIndex][]RPCHandler
	TraceIDGenerator TraceIDGenerator
}

func (so *ServerOptions) Sanitize() {
	if so.TraceIDGenerator == nil {
		so.TraceIDGenerator = generateTraceID
	}
}

type ServerMiddleware func(oldHandler http.Handler) (newHandler http.Handler)

func FillServerMiddlewareTable(serverMiddlewareTable [][]ServerMiddleware, serverMiddlewares map[MethodIndex][]ServerMiddleware) {
	commonServerMiddlewares := serverMiddlewares[AnyMethod]
	for i := range serverMiddlewareTable {
		methodIndex := MethodIndex(i)
		oldServerMiddlewares := serverMiddlewares[methodIndex]
		if len(oldServerMiddlewares) == 0 {
			serverMiddlewareTable[i] = commonServerMiddlewares
			continue
		}
		if len(commonServerMiddlewares) == 0 {
			serverMiddlewareTable[i] = oldServerMiddlewares
			continue
		}
		newServerMiddlewares := make([]ServerMiddleware, len(commonServerMiddlewares)+len(oldServerMiddlewares))
		copy(newServerMiddlewares, commonServerMiddlewares)
		copy(newServerMiddlewares[len(commonServerMiddlewares):], oldServerMiddlewares)
		serverMiddlewareTable[i] = newServerMiddlewares
	}
}

type IncomingRPCFactory func() (incomingRPC *IncomingRPC)
type TraceIDGenerator func() (traceID string)

func MakeHandler(
	serverMiddlewares []ServerMiddleware,
	incomingRPCFactory IncomingRPCFactory,
	traceIDGenerator TraceIDGenerator,
) http.Handler {
	handler := http.Handler(http.HandlerFunc(handleHTTP))
	for i := len(serverMiddlewares) - 1; i >= 0; i-- {
		serverMiddleware := serverMiddlewares[i]
		handler = serverMiddleware(handler)
	}
	handler = func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			var buffer bytes.Buffer
			if _, err := buffer.ReadFrom(r.Body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			incomingRPC := incomingRPCFactory()
			traceID := extractTraceID(r.Header)
			if traceID == "" {
				traceID = traceIDGenerator()
				injectTraceID(traceID, w.Header())
			}
			incomingRPC.traceID = traceID
			if buffer.Len() >= 1 {
				incomingRPC.rawParams = buffer.Bytes()
			}
			ctx := makeContextWithRPC(r.Context(), &incomingRPC.RPC)
			handler.ServeHTTP(responseWriterWrapper{w, incomingRPC}, r.WithContext(ctx))
		})
	}(handler)
	return handler
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	incomingRPC := MustGetRPCFromContext(ctx).IncomingRPC()
	defer func() {
		if v := recover(); v != nil {
			incomingRPC.savePanic(v)
		}
		if !incomingRPC.encodeResp(w) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(incomingRPC.rawResp)
	}()
	if !incomingRPC.decodeParams(ctx) {
		return
	}
	if err := incomingRPC.Do(ctx); err != nil {
		incomingRPC.saveErr(err)
	}
}

var _ http.HandlerFunc = handleHTTP

type responseWriterWrapper struct {
	http.ResponseWriter

	incomingRPC *IncomingRPC
}

var _ http.ResponseWriter = responseWriterWrapper{}

func (rww responseWriterWrapper) WriteHeader(statusCode int) {
	rww.ResponseWriter.WriteHeader(statusCode)
	rww.incomingRPC.statusCode = statusCode
}

func (rww responseWriterWrapper) Write(data []byte) (int, error) {
	n, err := rww.ResponseWriter.Write(data)
	rww.incomingRPC.responseWriteErr = err
	return n, err
}
