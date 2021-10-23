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

	rawParams        []byte
	error            Error
	err              error
	stackTrace       string
	rawResp          []byte
	statusCode       int
	responseWriteErr error
}

func (ir *IncomingRPC) Init(
	namespace string,
	serviceName string,
	methodName string,
	fullMethodName string,
	methodIndex MethodIndex,
	params Model,
	results Model,
	handler RPCHandler,
	filters []RPCHandler,
) {
	ir.mark = 'i'
	ir.init(namespace, serviceName, methodName, fullMethodName, methodIndex, params, results, handler, filters)
	ir.statusCode = http.StatusOK
}

func (ir *IncomingRPC) RawParams() []byte                { return ir.rawParams }
func (ir *IncomingRPC) UpdateRawParams(rawParams []byte) { ir.rawParams = rawParams }
func (ir *IncomingRPC) Error() Error                     { return ir.error }
func (ir *IncomingRPC) Err() error                       { return ir.err }
func (ir *IncomingRPC) StackTrace() string               { return ir.stackTrace }
func (ir *IncomingRPC) RawResp() []byte                  { return ir.rawResp }
func (ir *IncomingRPC) StatusCode() int                  { return ir.statusCode }
func (ir *IncomingRPC) ResponseWriteErr() error          { return ir.responseWriteErr }

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
	ir.err = err
}

func (ir *IncomingRPC) savePanic(v interface{}) {
	ir.error = *errInternal
	errText := fmt.Sprintf("%v", v)
	buffer := make([]byte, 4096)
	n := runtime.Stack(buffer, false)
	stackTrace := string(buffer[:n])
	if DebugMode {
		ir.error.Details = errText
		ir.error.Data.SetValue("stackTrace", stackTrace)
	}
	ir.err = errors.New(errText)
	ir.stackTrace = stackTrace
}

func (ir *IncomingRPC) encodeResp(responseWriter http.ResponseWriter) bool {
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
	if ir.error.Code == 0 {
		resp.Results = ir.results
	} else {
		resp.Error = &ir.error
	}
	if err := encoder.Encode(resp); err != nil {
		err = fmt.Errorf("resp encoding failed: %v", err)
		ir.RespondHTTPWithErr(responseWriter, http.StatusInternalServerError, err, "")
		return false
	}
	ir.rawResp = buffer.Bytes()
	if !DebugMode {
		// Remove '\n'
		ir.rawResp = ir.rawResp[:len(ir.rawResp)-1]
	}
	return true
}

func (ir *IncomingRPC) RespondHTTPWithErr(responseWriter http.ResponseWriter, statusCode int, err error, stackTrace string) {
	ir.error = Error{}
	ir.err = err
	ir.stackTrace = stackTrace
	ir.rawResp = nil
	if !DebugMode {
		responseWriter.WriteHeader(statusCode)
		return
	}
	responseWriter.Header().Set("Content-Type", "text/plain; charest=utf-8")
	responseWriter.WriteHeader(statusCode)
	var buffer bytes.Buffer
	buffer.WriteString(err.Error())
	buffer.WriteByte('\n')
	if stackTrace != "" {
		buffer.WriteString("stack trace:\n")
		buffer.WriteString(stackTrace)
		buffer.WriteByte('\n')
	}
	responseWriter.Write(buffer.Bytes())
}
