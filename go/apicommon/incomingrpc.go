package apicommon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"unsafe"
)

func (r *RPC) IncomingRPC() *IncomingRPC {
	if r.mark != 'i' {
		panic("not incoming rpc")
	}
	return (*IncomingRPC)(unsafe.Pointer(r))
}

type IncomingRPC struct {
	RPC

	traceIDIsReceived bool
	rawParams         []byte
	error             Error
	internalErr       error
	stackTrace        string
	encodeRespErr     error
	rawResp           []byte
	writeRespErr      error
}

func (ir *IncomingRPC) Init(
	namespace string,
	serviceName string,
	methodName string,
	params Model,
	results Model,
	handler RPCHandler,
	filters []RPCHandler,
) {
	ir.mark = 'i'
	ir.init(namespace, serviceName, methodName, params, results, handler, filters)
}

func (ir *IncomingRPC) RawParams() []byte    { return ir.rawParams }
func (ir *IncomingRPC) Error() Error         { return ir.error }
func (ir *IncomingRPC) InternalErr() error   { return ir.internalErr }
func (ir *IncomingRPC) StackTrace() string   { return ir.stackTrace }
func (ir *IncomingRPC) EncodeRespErr() error { return ir.encodeRespErr }
func (ir *IncomingRPC) RawResp() []byte      { return ir.rawResp }
func (ir *IncomingRPC) WriteRespErr() error  { return ir.writeRespErr }

func (ir *IncomingRPC) decodeParams(ctx context.Context) bool {
	if ir.params == nil {
		return true
	}
	if err := json.Unmarshal(ir.rawParams, ir.params); err != nil {
		switch err.(type) {
		case *json.SyntaxError:
			ir.error = *errParse
			ir.error.Details = err.Error()
		case *json.UnmarshalTypeError:
			ir.error = *errInvalidParams
			ir.error.Details = err.Error()
		default:
			ir.saveErr(err)
		}
		return false
	}
	validationContext := NewValidationContext(ctx)
	if !ir.params.Validate(validationContext) {
		ir.error = *errInvalidParams
		ir.error.Details = validationContext.ErrorDetails()
		return false
	}
	return true
}

func (ir *IncomingRPC) saveErr(err error) {
	if error, ok := err.(*Error); ok {
		ir.error = *error
		return
	}
	ir.error = *errInternal
	if DebugMode {
		ir.error.Details = err.Error()
	}
	ir.internalErr = err
}

func (ir *IncomingRPC) savePanic(v interface{}) {
	ir.error = *errInternal
	errStr := fmt.Sprintf("%v", v)
	buffer := make([]byte, 4096)
	n := runtime.Stack(buffer, false)
	stackTrace := string(buffer[:n])
	if DebugMode {
		ir.error.Details = errStr
		ir.error.Data.SetValue("stackTrace", stackTrace)
	}
	ir.internalErr = errors.New(errStr)
	ir.stackTrace = stackTrace
}

func (ir *IncomingRPC) encodeAndWriteResp(responseWriter http.ResponseWriter) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if DebugMode {
		encoder.SetIndent("", "  ")
	}
	resp := struct {
		TraceID *string     `json:"traceID,omitempty"`
		Error   *Error      `json:"error,omitempty"`
		Results interface{} `json:"results,omitempty"`
	}{}
	if !ir.traceIDIsReceived {
		resp.TraceID = &ir.traceID
	}
	if ir.error.Code != 0 {
		resp.Error = &ir.error
	} else {
		resp.Results = ir.results
	}
	if err := encoder.Encode(resp); err != nil {
		ir.encodeRespErr = err
		return
	}
	ir.rawResp = buffer.Bytes()
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_, ir.writeRespErr = responseWriter.Write(ir.rawResp)
}
