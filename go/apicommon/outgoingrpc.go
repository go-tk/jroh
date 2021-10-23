package apicommon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unsafe"
)

func (r *RPC) OutgoingRPC() *OutgoingRPC {
	if r.mark != 'o' {
		panic("not outgoing rpc")
	}
	return (*OutgoingRPC)(unsafe.Pointer(r))
}

type OutgoingRPC struct {
	RPC

	traceIDIsReceived bool
	client            *http.Client
	transport         http.RoundTripper

	url         string
	rawParams   []byte
	isRequested bool
	statusCode  int
	rawResp     []byte
	error       Error
}

func (or *OutgoingRPC) Init(
	namespace string,
	serviceName string,
	methodName string,
	fullMethodName string,
	methodIndex MethodIndex,
	params Model,
	results Model,
	filters []RPCHandler,
) {
	or.mark = 'o'
	or.init(namespace, serviceName, methodName, fullMethodName, methodIndex, params, results, handleRPC, filters)
}

func (or *OutgoingRPC) URL() string                  { return or.url }
func (or *OutgoingRPC) RawParams() []byte            { return or.rawParams }
func (or *OutgoingRPC) IsRequested() bool            { return or.isRequested }
func (or *OutgoingRPC) StatusCode() int              { return or.statusCode }
func (or *OutgoingRPC) RawResp() []byte              { return or.rawResp }
func (or *OutgoingRPC) UpdateRawResp(rawResp []byte) { or.rawResp = rawResp }
func (or *OutgoingRPC) Error() Error                 { return or.error }

func (or *OutgoingRPC) encodeParams() error {
	if or.params == nil {
		return nil
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if DebugMode {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(or.params); err != nil {
		return fmt.Errorf("params encoding failed: %v", err)
	}
	or.rawParams = buffer.Bytes()
	if !DebugMode {
		// Remove '\n'
		or.rawParams = or.rawParams[:len(or.rawParams)-1]
	}
	return nil
}

type UnexpectedStatusCodeError struct {
	StatusCode int
}

func (usce *UnexpectedStatusCodeError) Error() string {
	return "unexpected status code - " + strconv.Itoa(usce.StatusCode)
}

func (or *OutgoingRPC) requestHTTP(ctx context.Context) error {
	requestBody := bytes.NewReader(or.rawParams)
	request, err := http.NewRequestWithContext(ctx, "POST", or.url, requestBody)
	if err != nil {
		return fmt.Errorf("http request failed (1): %v", err)
	}
	if or.traceIDIsReceived {
		injectTraceID(or.traceID, request.Header)
	}
	if _, err := or.client.Do(request); err != nil {
		return fmt.Errorf("http request failed (2): %w", err)
	}
	if or.statusCode != http.StatusOK {
		return fmt.Errorf("http request failed (3): %w", &UnexpectedStatusCodeError{or.statusCode})
	}
	return nil
}

func (or *OutgoingRPC) decodeResp(ctx context.Context) error {
	if or.results == nil {
		resp := struct {
			TraceID *string `json:"traceID"`
			Error   *Error  `json:"error"`
		}{
			TraceID: &or.traceID,
			Error:   &or.error,
		}
		if err := json.Unmarshal(or.rawResp, &resp); err != nil {
			return fmt.Errorf("resp decoding failed (1): %v", err)
		}
		return nil
	}
	resp := struct {
		TraceID *string     `json:"traceID"`
		Error   *Error      `json:"error"`
		Results interface{} `json:"results"`
	}{
		TraceID: &or.traceID,
		Error:   &or.error,
		Results: or.results,
	}
	if err := json.Unmarshal(or.rawResp, &resp); err != nil {
		return fmt.Errorf("resp decoding failed (2): %v", err)
	}
	if or.error.Code == 0 {
		validationContext := NewValidationContext(ctx)
		if !or.results.Validate(validationContext) {
			return fmt.Errorf("resp decoding failed (3): invalid results: %s", validationContext.ErrorDetails())
		}
	}
	return nil
}

func handleRPC(ctx context.Context, rpc *RPC) error {
	outgoingRPC := rpc.OutgoingRPC()
	if err := outgoingRPC.encodeParams(); err != nil {
		return err
	}
	if err := outgoingRPC.requestHTTP(ctx); err != nil {
		return err
	}
	if err := outgoingRPC.decodeResp(ctx); err != nil {
		return err
	}
	if outgoingRPC.error.Code != 0 {
		return &outgoingRPC.error
	}
	return nil
}

var _ RPCHandler = handleRPC
