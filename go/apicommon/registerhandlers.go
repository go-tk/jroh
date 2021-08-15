package apicommon

import (
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"net/http"
)

type RegisterHandlersOptions struct {
	Middlewares      map[MethodIndex][]Middleware
	RPCInterceptors  map[MethodIndex][]RPCHandler
	TraceIDGenerator TraceIDGenerator
}

func (rho *RegisterHandlersOptions) Normalize(numberOfMethods int) {
	if commonMiddlewares := rho.Middlewares[AnyMethod]; len(commonMiddlewares) >= 1 {
		middlewares := make(map[MethodIndex][]Middleware, numberOfMethods)
		for methodIndex := MethodIndex(0); methodIndex < MethodIndex(numberOfMethods); methodIndex++ {
			oldMiddlewares := rho.Middlewares[methodIndex]
			if len(oldMiddlewares) == 0 {
				middlewares[methodIndex] = commonMiddlewares
				continue
			}
			newMiddlewares := make([]Middleware, len(commonMiddlewares)+len(oldMiddlewares))
			copy(newMiddlewares, commonMiddlewares)
			copy(newMiddlewares[len(commonMiddlewares):], oldMiddlewares)
			middlewares[methodIndex] = newMiddlewares
		}
		rho.Middlewares = middlewares
	}
	if commonRPCInterceptors := rho.RPCInterceptors[AnyMethod]; len(commonRPCInterceptors) >= 1 {
		rpcInterceptors := make(map[MethodIndex][]RPCHandler, numberOfMethods)
		for methodIndex := MethodIndex(0); methodIndex < MethodIndex(numberOfMethods); methodIndex++ {
			oldRPCInterceptors := rho.RPCInterceptors[methodIndex]
			if len(oldRPCInterceptors) == 0 {
				rpcInterceptors[methodIndex] = commonRPCInterceptors
			}
			newRPCInterceptors := make([]RPCHandler, len(commonRPCInterceptors)+len(oldRPCInterceptors))
			copy(newRPCInterceptors, commonRPCInterceptors)
			copy(newRPCInterceptors[len(commonRPCInterceptors):], oldRPCInterceptors)
			rpcInterceptors[methodIndex] = newRPCInterceptors
		}
		rho.RPCInterceptors = rpcInterceptors
	}
	if rho.TraceIDGenerator == nil {
		rho.TraceIDGenerator = func() string {
			var buffer [16]byte
			binary.BigEndian.PutUint64(buffer[:8], rand.Uint64())
			binary.BigEndian.PutUint64(buffer[8:], rand.Uint64())
			traceID := base64.RawURLEncoding.EncodeToString(buffer[:])
			return traceID
		}
	}
}

type MethodIndex int

const AnyMethod MethodIndex = -1

type Middleware func(oldHandler http.Handler) (newHandler http.Handler)
type TraceIDGenerator func() (traceID string)
type IncomingRPCFactory func(traceID string) (incomingRPC *IncomingRPC)

func MakeHandler(
	middlewares []Middleware,
	traceIDGenerator TraceIDGenerator,
	incomingRPCFactory IncomingRPCFactory,
) http.Handler {
	handler := http.Handler(http.HandlerFunc(handleIncomingRPC))
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		handler = middleware(handler)
	}
	handler = func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceID := r.Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = traceIDGenerator()
			}
			incomingRPC := incomingRPCFactory(traceID)
			ctx := makeContextWithRPC(r.Context(), &incomingRPC.RPC)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}(handler)
	return handler
}
