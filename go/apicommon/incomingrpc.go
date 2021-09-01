package apicommon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"unsafe"
)

func (r *RPC) IncomingRPC() *IncomingRPC { return (*IncomingRPC)(unsafe.Pointer(r)) }

type IncomingRPC struct {
	RPC

	rawParams    []byte
	error        Error
	internalErr  error
	stackTrace   string
	rawResp      []byte
	writeRespErr error
}

func (ir *IncomingRPC) RawParams() []byte   { return ir.rawParams }
func (ir *IncomingRPC) Error() Error        { return ir.error }
func (ir *IncomingRPC) InternalErr() error  { return ir.internalErr }
func (ir *IncomingRPC) StackTrace() string  { return ir.stackTrace }
func (ir *IncomingRPC) RawResp() []byte     { return ir.rawResp }
func (ir *IncomingRPC) WriteRespErr() error { return ir.writeRespErr }

func (ir *IncomingRPC) readParams() bool {
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
	validationContext := NewValidationContext()
	if !ir.params.(Validator).Validate(validationContext) {
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

func (ir *IncomingRPC) writeResp(responseWriter http.ResponseWriter) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if DebugMode {
		encoder.SetIndent("", "  ")
	}
	var error *Error
	if ir.error.Code != 0 {
		error = &ir.error
	}
	if ir.results == nil {
		if err := encoder.Encode(struct {
			TraceID string `json:"traceID"`
			Error   *Error `json:"error,omitempty"`
		}{
			TraceID: ir.traceID,
			Error:   error,
		}); err != nil {
			ir.writeRespErr = err
			return
		}
	} else {
		if err := encoder.Encode(struct {
			TraceID string      `json:"traceID"`
			Error   *Error      `json:"error,omitempty"`
			Results interface{} `json:"results,omitempty"`
		}{
			TraceID: ir.traceID,
			Error:   error,
			Results: ir.results,
		}); err != nil {
			ir.writeRespErr = err
			return
		}
	}
	ir.rawResp = buffer.Bytes()
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_, ir.writeRespErr = responseWriter.Write(ir.rawResp)
}

func handleIncomingHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	incomingRPC := MustGetRPCFromContext(ctx).IncomingRPC()
	defer func() {
		if v := recover(); v != nil {
			incomingRPC.savePanic(v)
		}
		incomingRPC.writeResp(w)
	}()
	if !incomingRPC.readParams() {
		return
	}
	if err := incomingRPC.Do(ctx); err != nil {
		incomingRPC.saveErr(err)
	}
}

var _ http.HandlerFunc = handleIncomingHTTP
