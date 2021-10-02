package apicommon

import (
	"bytes"
	"encoding/base64"
	"math/rand"
	"net/http"
)

type MethodIndex int

const AnyMethod MethodIndex = -1

func FillRPCFiltersTable(rpcFiltersTable [][]RPCHandler, rpcFilters map[MethodIndex][]RPCHandler) {
	commonRPCFilters := rpcFilters[AnyMethod]
	for i := range rpcFiltersTable {
		methodIndex := MethodIndex(i)
		oldRPCFilters := rpcFilters[methodIndex]
		if len(oldRPCFilters) == 0 {
			rpcFiltersTable[i] = commonRPCFilters
			continue
		}
		if len(commonRPCFilters) == 0 {
			rpcFiltersTable[i] = oldRPCFilters
			continue
		}
		newRPCFilters := make([]RPCHandler, len(commonRPCFilters)+len(oldRPCFilters))
		copy(newRPCFilters, commonRPCFilters)
		copy(newRPCFilters[len(commonRPCFilters):], oldRPCFilters)
		rpcFiltersTable[i] = newRPCFilters
	}
}

func wrapHandler(handler http.Handler, serverMiddlewares map[MethodIndex][]ServerMiddleware, methodIndex MethodIndex) http.Handler {
	serverMiddlewares2 := serverMiddlewares[methodIndex]
	for i := len(serverMiddlewares2) - 1; i >= 0; i-- {
		serverMiddleware := serverMiddlewares2[i]
		handler = serverMiddleware(handler)
	}
	serverMiddlewares3 := serverMiddlewares[AnyMethod]
	for i := len(serverMiddlewares3) - 1; i >= 0; i-- {
		serverMiddleware := serverMiddlewares3[i]
		handler = serverMiddleware(handler)
	}
	return handler
}

func FillTransportTable(transportTable []http.RoundTripper, transport http.RoundTripper, clientMiddlewares map[MethodIndex][]ClientMiddleware) {
	transportFunc := TransportFunc(func(request *http.Request) (*http.Response, error) {
		outgoingRPC := MustGetRPCFromContext(request.Context()).OutgoingRPC()
		outgoingRPC.requestIsSent = true
		response, err := transport.RoundTrip(request)
		if err != nil {
			return nil, err
		}
		outgoingRPC.statusCode = response.StatusCode
		if response.StatusCode != http.StatusOK {
			response.Body.Close()
			return response, nil
		}
		var buffer bytes.Buffer
		_, err = buffer.ReadFrom(response.Body)
		response.Body.Close()
		if err != nil {
			return nil, err
		}
		if buffer.Len() >= 1 {
			outgoingRPC.rawResp = buffer.Bytes()
		}
		return response, nil
	})
	for i := range transportTable {
		methodIndex := MethodIndex(i)
		transportTable[methodIndex] = wrapTransport(transportFunc, clientMiddlewares, methodIndex)
	}
}

func wrapTransport(transport http.RoundTripper, clientMiddlewares map[MethodIndex][]ClientMiddleware, methodIndex MethodIndex) http.RoundTripper {
	clientMiddlewares2 := clientMiddlewares[methodIndex]
	for i := len(clientMiddlewares2) - 1; i >= 0; i-- {
		clientMiddleware := clientMiddlewares2[i]
		transport = clientMiddleware(transport)
	}
	clientMiddlewares3 := clientMiddlewares[AnyMethod]
	for i := len(clientMiddlewares3) - 1; i >= 0; i-- {
		clientMiddleware := clientMiddlewares3[i]
		transport = clientMiddleware(transport)
	}
	return transport
}

func generateTraceID() string {
	var buffer [16]byte
	rand.Read(buffer[:])
	traceID := base64.RawURLEncoding.EncodeToString(buffer[:])
	return traceID
}

const traceIDHeaderKey = "X-JROH-Trace-ID"

func injectTraceID(traceID string, header http.Header) { header.Set(traceIDHeaderKey, traceID) }
func extractTraceID(header http.Header) string         { return header.Get(traceIDHeaderKey) }
