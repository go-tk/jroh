package apicommon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	error             Error

	url        string
	rawParams  []byte
	statusCode int
	rawResp    []byte
}

func (or *OutgoingRPC) Init(
	namespace string,
	serviceName string,
	methodName string,
	params Model,
	results Model,
	handler RPCHandler,
	filters []RPCHandler,
) {
	or.mark = 'o'
	or.init(namespace, serviceName, methodName, params, results, handler, filters)
}

func (or *OutgoingRPC) URL() string                  { return or.url }
func (or *OutgoingRPC) RawParams() []byte            { return or.rawParams }
func (or *OutgoingRPC) StatusCode() int              { return or.statusCode }
func (or *OutgoingRPC) RawResp() []byte              { return or.rawResp }
func (or *OutgoingRPC) UpdateRawResp(rawResp []byte) { or.rawResp = rawResp }

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
		return err
	}
	or.rawParams = buffer.Bytes()
	if !DebugMode {
		// Remove '\n'
		or.rawParams = or.rawParams[:len(or.rawParams)-1]
	}
	return nil
}

func (or *OutgoingRPC) requestHTTP(ctx context.Context) error {
	requestBody := bytes.NewReader(or.rawParams)
	request, err := http.NewRequestWithContext(ctx, "POST", or.url, requestBody)
	if err != nil {
		return err
	}
	if or.traceIDIsReceived {
		injectTraceID(or.traceID, request.Header)
	}
	if _, err := or.client.Do(request); err != nil {
		return err
	}
	if or.statusCode != http.StatusOK {
		return &UnexpectedStatusCodeError{or.statusCode}
	}
	return err
}

type UnexpectedStatusCodeError struct {
	StatusCode int
}

func (usce *UnexpectedStatusCodeError) Error() string {
	return "unexpected status code - %v" + strconv.Itoa(usce.StatusCode)
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
		return json.Unmarshal(or.rawResp, &resp)
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
		return err
	}
	if or.error.Code == 0 {
		validationContext := NewValidationContext(ctx)
		if !or.results.Validate(validationContext) {
			return errors.New("invalid results: " + validationContext.ErrorDetails())
		}
	}
	return nil
}
