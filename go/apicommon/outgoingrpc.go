package apicommon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

var ErrUnexpectedStatus = errors.New("api: unexpected status")

func HandleRPC(ctx context.Context, rpc *RPC) error {
	outgoingRPC := rpc.OutgoingRPC()
	if err := outgoingRPC.encodeParams(); err != nil {
		return err
	}
	requestBody := bytes.NewReader(outgoingRPC.rawParams)
	request, err := http.NewRequestWithContext(ctx, "POST", outgoingRPC.url, requestBody)
	if err != nil {
		return err
	}
	if rpc.traceID != "" {
		injectTraceID(rpc.traceID, request.Header)
	}
	response, err := outgoingRPC.client.Do(request)
	if err != nil {
		return err
	}
	outgoingRPC.statusCode = response.StatusCode
	if outgoingRPC.statusCode != http.StatusOK {
		return fmt.Errorf("%w: statusCode=%v", ErrUnexpectedStatus, outgoingRPC.statusCode)
	}
	if err := outgoingRPC.readResp(response.Body); err != nil {
		return err
	}
	if outgoingRPC.error.Code != 0 {
		return &outgoingRPC.error
	}
	return nil
}

var _ RPCHandler = HandleRPC
