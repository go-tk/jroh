package apicommon

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

type ServerOptions struct {
	RPCFilters       map[int][]IncomingRPCHandler
	TraceIDGenerator TraceIDGenerator
}

type TraceIDGenerator func() (traceID string)

func (so *ServerOptions) Sanitize() {
	if so.TraceIDGenerator == nil {
		so.TraceIDGenerator = defaultTraceIDGenerator
	}
}

func defaultTraceIDGenerator() string {
	var buffer [16]byte
	rand.Read(buffer[:])
	traceID := base64.RawURLEncoding.EncodeToString(buffer[:])
	return traceID
}

func (so *ServerOptions) AddCommonRPCFilters(rpcFilters ...IncomingRPCHandler) {
	so.AddRPCFilters(-1, rpcFilters...)
}

func (so *ServerOptions) AddRPCFilters(methodIndex int, rpcFilters ...IncomingRPCHandler) {
	if so.RPCFilters == nil {
		so.RPCFilters = make(map[int][]IncomingRPCHandler)
	}
	so.RPCFilters[methodIndex] = append(so.RPCFilters[methodIndex], rpcFilters...)
}

func FillIncomingRPCFiltersTable(rpcFiltersTable [][]IncomingRPCHandler, rpcFilters map[int][]IncomingRPCHandler) {
	commonRPCFilters := rpcFilters[-1]
	for i := range rpcFiltersTable {
		oldRPCFilters := rpcFilters[i]
		if len(oldRPCFilters) == 0 {
			rpcFiltersTable[i] = commonRPCFilters
			continue
		}
		if len(commonRPCFilters) == 0 {
			rpcFiltersTable[i] = oldRPCFilters
			continue
		}
		newRPCFilters := make([]IncomingRPCHandler, len(commonRPCFilters)+len(oldRPCFilters))
		copy(newRPCFilters, commonRPCFilters)
		copy(newRPCFilters[len(commonRPCFilters):], oldRPCFilters)
		rpcFiltersTable[i] = newRPCFilters
	}
}

func HandleRequest(
	request *http.Request,
	incomingRPC *IncomingRPC,
	traceIDGenerator TraceIDGenerator,
	responseWriter http.ResponseWriter,
) {
	incomingRPC.RemoteIP = getRemoteIP(request)
	incomingRPC.InboundHeader = request.Header
	incomingRPC.OutboundHeader = responseWriter.Header()
	traceID, ok := getTraceID(incomingRPC.InboundHeader)
	if !ok {
		traceID = traceIDGenerator()
		setTraceID(traceID, incomingRPC.OutboundHeader)
	}
	incomingRPC.TraceID = traceID
	incomingRPC.StatusCode = http.StatusOK
	incomingRPC.SetReader(request.Body)
	if err := incomingRPC.Do(request.Context()); err != nil {
		error, ok := err.(*Error)
		if !ok {
			error = &Error{
				Code:       -1,
				StatusCode: incomingRPC.StatusCode,
				Message:    statusCode2Message[incomingRPC.StatusCode],
			}
			if error.StatusCode/100 != 5 {
				error.Details = err.Error()
			}
		}
		rawError, err := encodeError(error)
		if err != nil {
			panic(err)
		}
		setContentTypeAsJSON(incomingRPC.OutboundHeader)
		setErrorCode(error.Code, incomingRPC.OutboundHeader)
		responseWriter.WriteHeader(incomingRPC.StatusCode)
		responseWriter.Write(rawError)
		return
	}
	incomingRPC.EncodeResults()
	setContentTypeAsJSON(incomingRPC.OutboundHeader)
	responseWriter.WriteHeader(incomingRPC.StatusCode)
	responseWriter.Write(incomingRPC.RawResults)
}

func getRemoteIP(request *http.Request) string {
	if values := request.Header["X-Forwarded-For"]; len(values) >= 1 {
		ips := values[0]
		i := strings.IndexByte(ips, ',')
		if i < 0 {
			i = len(ips)
		}
		if i >= 1 {
			return ips[:i]
		}
	}
	remoteAddr := request.RemoteAddr
	i := strings.LastIndexByte(remoteAddr, ':')
	if i < 0 {
		i = len(remoteAddr)
	}
	return remoteAddr[:i]
}

var statusCode2Message = map[int]string{
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
	400: "bad request",
	401: "unauthorized",
	402: "payment required",
	403: "forbidden",
	404: "not found",
	405: "method not allowed",
	406: "not acceptable",
	407: "proxy authentication required",
	408: "request timeout",
	409: "conflict",
	410: "gone",
	411: "length required",
	412: "precondition failed",
	413: "payload too large",
	414: "uri too long",
	415: "unsupported media type",
	416: "range not satisfiable",
	417: "expectation failed",
	418: "i'm a teapot",
	422: "unprocessable entity",
	425: "too early",
	426: "upgrade required",
	428: "precondition required",
	429: "too many requests",
	431: "request header fields too large",
	451: "unavailable for legal reasons",
	500: "internal server error",
	501: "not implemented",
	502: "bad gateway",
	503: "service unavailable",
	504: "gateway timeout",
	505: "http version not supported",
	506: "variant also negotiates",
	507: "insufficient storage",
	508: "loop detected",
	510: "not extended",
	511: "network authentication required",
}

func encodeError(error *Error) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if DebugMode {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(error); err != nil {
		return nil, fmt.Errorf("error encoding failed: %w", err)
	}
	rawError := buffer.Bytes()
	if !DebugMode {
		// Remove trailing '\n'
		rawError = rawError[:len(rawError)-1]
	}
	return rawError, nil
}
