package apicommon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type ClientOptions struct {
	RPCFilters map[int][]OutgoingRPCHandler
	Timeout    time.Duration
	Transport  http.RoundTripper
}

func (co *ClientOptions) Sanitize() {
	if co.Transport == nil {
		co.Transport = http.DefaultTransport
	}
}

func (co *ClientOptions) AddCommonRPCFilters(rpcFilters ...OutgoingRPCHandler) {
	co.AddRPCFilters(-1, rpcFilters...)
}

func (co *ClientOptions) AddRPCFilters(methodIndex int, rpcFilters ...OutgoingRPCHandler) {
	if co.RPCFilters == nil {
		co.RPCFilters = make(map[int][]OutgoingRPCHandler)
	}
	co.RPCFilters[methodIndex] = append(co.RPCFilters[methodIndex], rpcFilters...)
}

func FillOutgoingRPCFiltersTable(rpcFiltersTable [][]OutgoingRPCHandler, rpcFilters map[int][]OutgoingRPCHandler) {
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
		newRPCFilters := make([]OutgoingRPCHandler, len(commonRPCFilters)+len(oldRPCFilters))
		copy(newRPCFilters, commonRPCFilters)
		copy(newRPCFilters[len(commonRPCFilters):], oldRPCFilters)
		rpcFiltersTable[i] = newRPCFilters
	}
}

var outboundHeaderPool = sync.Pool{New: func() interface{} { return make(http.Header) }}

func HandleOutgoingRPC(ctx context.Context, outgoingRPC *OutgoingRPC) error {
	if err := outgoingRPC.EncodeParams(); err != nil {
		return err
	}
	requestBody := bytes.NewReader(outgoingRPC.RawParams)
	request, err := http.NewRequestWithContext(ctx, "POST", outgoingRPC.URL, requestBody)
	if err != nil {
		return err
	}
	outboundHeaderPool.Put(request.Header)
	request.Header = outgoingRPC.OutboundHeader
	setContentTypeAsJSON(outgoingRPC.OutboundHeader)
	if incomingRPC, ok := GetIncomingRPCFromContext(ctx); ok {
		setTraceID(incomingRPC.TraceID, outgoingRPC.OutboundHeader)
		outgoingRPC.TraceID = incomingRPC.TraceID
	}
	response, err := outgoingRPC.Transport.RoundTrip(request)
	if err != nil {
		return err
	}
	outgoingRPC.StatusCode = response.StatusCode
	outgoingRPC.InboundHeader = response.Header
	if traceID, ok := getTraceID(outgoingRPC.InboundHeader); ok {
		outgoingRPC.TraceID = traceID
	}
	if errorCode, ok := getErrorCode(outgoingRPC.InboundHeader); ok {
		outgoingRPC.ErrorCode = errorCode
		error := Error{
			Code:       outgoingRPC.ErrorCode,
			StatusCode: outgoingRPC.StatusCode,
		}
		err := loadError(&error, response.Body)
		response.Body.Close()
		if err != nil {
			return err
		}
		return &error
	}
	if outgoingRPC.StatusCode != http.StatusOK {
		response.Body.Close()
		return fmt.Errorf("%w - %v", ErrUnexpectedStatusCode, outgoingRPC.StatusCode)
	}
	outgoingRPC.SetReadCloser(response.Body)
	return nil
}

var _ = OutgoingRPCHandler(HandleOutgoingRPC)

func loadError(error *Error, reader io.Reader) error {
	rawError, err := readRawError(reader)
	if err != nil {
		return err
	}
	if err := decodeRawError(rawError, error); err != nil {
		return err
	}
	return nil
}

func readRawError(reader io.Reader) ([]byte, error) {
	var buffer bytes.Buffer
	n, err := buffer.ReadFrom(reader)
	if err != nil {
		return nil, fmt.Errorf("raw error read failed: %w", err)
	}
	if n == 0 {
		return nil, nil
	}
	rawError := buffer.Bytes()
	return rawError, nil
}

func decodeRawError(rawError []byte, error *Error) error {
	if err := json.Unmarshal(rawError, error); err != nil {
		return fmt.Errorf("raw error decoding failed: %w", err)
	}
	return nil
}
