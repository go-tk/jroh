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
type IncomingRPCFactory func() (incomingRPC *IncomingRPC)
type TraceIDGenerator func() (traceID string)

func MakeHandler(
	serverMiddlewares map[MethodIndex][]ServerMiddleware,
	methodIndex MethodIndex,
	incomingRPCFactory IncomingRPCFactory,
	traceIDGenerator TraceIDGenerator,
) http.Handler {
	handler := wrapHandler(http.HandlerFunc(handleHTTP), serverMiddlewares, methodIndex)
	handler = func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var buffer bytes.Buffer
			if _, err := buffer.ReadFrom(r.Body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(nil)
				return
			}
			incomingRPC := incomingRPCFactory()
			if traceID := extractTraceID(r.Header); traceID == "" {
				incomingRPC.traceID = traceIDGenerator()
			} else {
				incomingRPC.traceID = traceID
				incomingRPC.traceIDIsReceived = true
			}
			if buffer.Len() >= 1 {
				incomingRPC.rawParams = buffer.Bytes()
			}
			ctx := makeContextWithRPC(r.Context(), &incomingRPC.RPC)
			handler.ServeHTTP(w, r.WithContext(ctx))
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
		incomingRPC.encodeAndWriteResp(w)
	}()
	if !incomingRPC.decodeParams(ctx) {
		return
	}
	if err := incomingRPC.Do(ctx); err != nil {
		incomingRPC.saveErr(err)
	}
}

var _ http.HandlerFunc = handleHTTP
