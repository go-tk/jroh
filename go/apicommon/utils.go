package apicommon

import (
	"encoding/base64"
	"math/rand"
	"net/http"
)

type RPCFilters map[MethodIndex][]RPCHandler

func (rf *RPCFilters) Add(methodIndex MethodIndex, items ...RPCHandler) {
	if *rf == nil {
		*rf = make(map[MethodIndex][]RPCHandler)
	}
	(*rf)[methodIndex] = append((*rf)[methodIndex], items...)
}

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

func generateTraceID() string {
	var buffer [16]byte
	rand.Read(buffer[:])
	traceID := base64.RawURLEncoding.EncodeToString(buffer[:])
	return traceID
}

const traceIDHeaderKey = "X-JROH-Trace-ID"

func injectTraceID(traceID string, header http.Header) { header[traceIDHeaderKey] = []string{traceID} }

func extractTraceID(header http.Header) string {
	if values := header[traceIDHeaderKey]; len(values) >= 1 {
		return values[0]
	}
	return ""
}
