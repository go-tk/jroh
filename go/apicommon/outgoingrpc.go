package apicommon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"unsafe"
)

func (r *RPC) OutgoingRPC() *OutgoingRPC { return (*OutgoingRPC)(unsafe.Pointer(r)) }

type OutgoingRPC struct {
	RPC

	client     *http.Client
	url        string
	rawParams  []byte
	statusCode int
	rawResp    []byte
	error      Error
}

func (or *OutgoingRPC) URL() string       { return or.url }
func (or *OutgoingRPC) RawParams() []byte { return or.rawParams }
func (or *OutgoingRPC) StatusCode() int   { return or.statusCode }
func (or *OutgoingRPC) RawResp() []byte   { return or.rawResp }
func (or *OutgoingRPC) Error() Error      { return or.error }

func (or *OutgoingRPC) encodeParams() error {
	if or.params == nil {
		return nil
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(or.params); err != nil {
		return err
	}
	or.rawParams = buffer.Bytes()
	return nil
}

func (or *OutgoingRPC) requestHTTP(ctx context.Context) (*http.Response, error) {
	requestBody := bytes.NewReader(or.rawParams)
	request, err := http.NewRequestWithContext(ctx, "POST", or.url, requestBody)
	if err != nil {
		return nil, err
	}
	if or.traceID != "" {
		injectTraceID(or.traceID, request.Header)
	}
	return or.client.Do(request)
}

func (or *OutgoingRPC) readResp(responseBody io.ReadCloser) error {
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(responseBody)
	responseBody.Close()
	if err != nil {
		return err
	}
	or.rawResp = buffer.Bytes()
	if or.results == nil {
		if err := json.Unmarshal(or.rawResp, &struct {
			TraceID *string `json:"traceID"`
			Error   *Error  `json:"error"`
		}{
			TraceID: &or.traceID,
			Error:   &or.error,
		}); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(or.rawResp, &struct {
			TraceID *string     `json:"traceID"`
			Error   *Error      `json:"error"`
			Results interface{} `json:"results"`
		}{
			TraceID: &or.traceID,
			Error:   &or.error,
			Results: or.results,
		}); err != nil {
			return err
		}
	}
	return nil
}

type UnexpectedStatusCodeError struct {
	StatusCode int

	namespace   string
	serviceName string
	methodName  string
	traceID     string
}

func (use *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf("apicommon: unexpected status code; namespace=%v serviceName=%v methodName=%v traceID=%v statusCode=%v",
		use.namespace, use.serviceName, use.methodName, use.traceID, use.StatusCode)
}

func HandleRPC(ctx context.Context, rpc *RPC) error {
	outgoingRPC := rpc.OutgoingRPC()
	if err := outgoingRPC.encodeParams(); err != nil {
		return fmt.Errorf("apicommon: params encoding failed; namespace=%v serviceName=%v methodName=%v traceID=%v: %v",
			rpc.namespace, rpc.serviceName, rpc.methodName, rpc.traceID, err)
	}
	response, err := outgoingRPC.requestHTTP(ctx)
	if err != nil {
		return fmt.Errorf("apicommon: http request failed; namespace=%v serviceName=%v methodName=%v traceID=%v: %v",
			rpc.namespace, rpc.serviceName, rpc.methodName, rpc.traceID, err)
	}
	outgoingRPC.statusCode = response.StatusCode
	if outgoingRPC.statusCode != http.StatusOK {
		return &UnexpectedStatusCodeError{
			StatusCode: outgoingRPC.statusCode,

			namespace:   rpc.namespace,
			serviceName: rpc.serviceName,
			methodName:  rpc.methodName,
			traceID:     rpc.traceID,
		}
	}
	if err := outgoingRPC.readResp(response.Body); err != nil {
		return fmt.Errorf("apicommon: resp read failed; namespace=%v serviceName=%v methodName=%v traceID=%v: %v",
			rpc.namespace, rpc.serviceName, rpc.methodName, rpc.traceID, err)
	}
	if outgoingRPC.error.Code != 0 {
		return fmt.Errorf("apicommon: rpc failed; namespace=%v serviceName=%v methodName=%v traceID=%v: %w",
			rpc.namespace, rpc.serviceName, rpc.methodName, rpc.traceID, &outgoingRPC.error)
	}
	return nil
}

var _ RPCHandler = HandleRPC
