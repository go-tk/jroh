package apicommon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"unsafe"
)

func (r *RPC) IncomingRPC() *IncomingRPC { return (*IncomingRPC)(unsafe.Pointer(r)) }

type IncomingRPC struct {
	RPC

	RawParams    []byte
	Error        Error
	InternalErr  error
	StackTrace   string
	WriteRespErr error
}

type Validator interface {
	Validate() (err error)
}

func (ir *IncomingRPC) readParams(reader io.Reader) bool {
	params := ir.Params()
	if params == nil {
		return true
	}
	if err := json.Unmarshal(ir.RawParams, params); err != nil {
		switch err.(type) {
		case *json.SyntaxError:
			ir.Error = *ErrParse
			ir.Error.Details = err.Error()
		case *json.UnmarshalTypeError:
			ir.Error = *ErrInvalidParams
			ir.Error.Details = err.Error()
		default:
			ir.saveErr(err)
		}
		return false
	}
	if validator, ok := params.(Validator); ok {
		if err := validator.Validate(); err != nil {
			ir.Error = *ErrInvalidParams
			ir.Error.Details = err.Error()
			return false
		}
	}
	return true
}

func (ir *IncomingRPC) saveErr(err error) {
	if error, ok := err.(*Error); ok {
		ir.Error = *error
		return
	}
	ir.Error = *ErrInternal
	if DebugMode {
		ir.Error.Details = err.Error()
	}
	ir.InternalErr = err
}

func (ir *IncomingRPC) savePanic(v interface{}) {
	ir.Error = *ErrInternal
	errStr := fmt.Sprintf("%v", v)
	buffer := make([]byte, 4096)
	n := runtime.Stack(buffer, false)
	stackTrace := string(buffer[:n])
	if DebugMode {
		ir.Error.Details = errStr
		ir.Error.Data.SetValue("stackTrace", stackTrace)
	}
	ir.InternalErr = errors.New(errStr)
	ir.StackTrace = stackTrace
}

func (ir *IncomingRPC) writeResp(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(responseWriter)
	if DebugMode {
		encoder.SetIndent("", "    ")
	}
	if ir.Error.Code == 0 {
		ir.WriteRespErr = encoder.Encode(struct {
			TraceID string      `json:"traceID"`
			Results interface{} `json:"results"`
		}{
			TraceID: ir.TraceID(),
			Results: ir.Results(),
		})
	} else {
		ir.WriteRespErr = encoder.Encode(struct {
			TraceID string `json:"traceID"`
			Error   *Error `json:"error"`
		}{
			TraceID: ir.TraceID(),
			Error:   &ir.Error,
		})
	}
}

func handleIncomingRPC(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	incomingRPC := MustGetRPCFromContext(ctx).IncomingRPC()
	defer func() {
		if v := recover(); v != nil {
			incomingRPC.savePanic(v)
		}
		incomingRPC.writeResp(w)
	}()
	if !incomingRPC.readParams(r.Body) {
		return
	}
	if err := incomingRPC.Do(ctx); err != nil {
		incomingRPC.saveErr(err)
	}
}

var _ http.HandlerFunc = handleIncomingRPC
