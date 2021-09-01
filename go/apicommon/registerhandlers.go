package apicommon

import (
	"bytes"
	"net/http"
)

type RegisterHandlersOptions struct {
	Middlewares      map[MethodIndex][]Middleware
	RPCInterceptors  map[MethodIndex][]RPCHandler
	TraceIDGenerator TraceIDGenerator
}

func (rho *RegisterHandlersOptions) Sanitize() {
	if rho.TraceIDGenerator == nil {
		rho.TraceIDGenerator = generateTraceID
	}
}

type MethodIndex int

const AnyMethod MethodIndex = -1

type Middleware func(oldHandler http.Handler) (newHandler http.Handler)
type TraceIDGenerator func() (traceID string)
type IncomingRPCFactory func() (incomingRPC *IncomingRPC)

func MakeHandler(
	middlewares []Middleware,
	incomingRPCFactory IncomingRPCFactory,
	traceIDGenerator TraceIDGenerator,
) http.Handler {
	handler := http.Handler(http.HandlerFunc(handleIncomingHTTP))
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		handler = middleware(handler)
	}
	handler = func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var buffer bytes.Buffer
			if _, err := buffer.ReadFrom(r.Body); err != nil {
				return
			}
			incomingRPC := incomingRPCFactory()
			if traceID := extractTraceID(r.Header); traceID == "" {
				incomingRPC.traceID = traceIDGenerator()
			} else {
				incomingRPC.traceID = traceID
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
