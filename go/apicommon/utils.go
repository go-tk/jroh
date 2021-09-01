package apicommon

import (
	"encoding/base64"
	"math/rand"
	"net/http"
)

func FillMiddlewareTable(middlewareTable [][]Middleware, middlewares map[MethodIndex][]Middleware) {
	commonMiddlewares := middlewares[AnyMethod]
	for i := range middlewareTable {
		methodIndex := MethodIndex(i)
		oldRPCInterceptors := middlewares[methodIndex]
		if len(oldRPCInterceptors) == 0 {
			middlewareTable[i] = commonMiddlewares
			continue
		}
		if len(commonMiddlewares) == 0 {
			middlewareTable[i] = oldRPCInterceptors
			continue
		}
		newMiddlewares := make([]Middleware, len(commonMiddlewares)+len(oldRPCInterceptors))
		copy(newMiddlewares, commonMiddlewares)
		copy(newMiddlewares[len(commonMiddlewares):], oldRPCInterceptors)
		middlewareTable[i] = newMiddlewares
	}
}

func FillRPCInterceptorTable(rpcInterceptorTable [][]RPCHandler, rpcInterceptors map[MethodIndex][]RPCHandler) {
	commonRPCInterceptors := rpcInterceptors[AnyMethod]
	for i := range rpcInterceptorTable {
		methodIndex := MethodIndex(i)
		oldRPCInterceptors := rpcInterceptors[methodIndex]
		if len(oldRPCInterceptors) == 0 {
			rpcInterceptorTable[i] = commonRPCInterceptors
			continue
		}
		if len(commonRPCInterceptors) == 0 {
			rpcInterceptorTable[i] = oldRPCInterceptors
			continue
		}
		newRPCInterceptors := make([]RPCHandler, len(commonRPCInterceptors)+len(oldRPCInterceptors))
		copy(newRPCInterceptors, commonRPCInterceptors)
		copy(newRPCInterceptors[len(commonRPCInterceptors):], oldRPCInterceptors)
		rpcInterceptorTable[i] = newRPCInterceptors
	}
}

func generateTraceID() string {
	var buffer [16]byte
	rand.Read(buffer[:])
	traceID := base64.RawURLEncoding.EncodeToString(buffer[:])
	return traceID
}

const traceIDHeaderKey = "X-Trace-ID"

func injectTraceID(traceID string, header http.Header) { header.Set(traceIDHeaderKey, traceID) }
func extractTraceID(header http.Header) string         { return header.Get(traceIDHeaderKey) }
