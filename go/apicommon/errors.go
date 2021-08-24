package apicommon

import (
	"strconv"
)

const ErrorParse ErrorCode = -32700

func NewParseError() *Error {
	return &Error{
		Code:    ErrorParse,
		Message: "parse error",
	}
}

var errParse *Error = NewParseError()
var ErrParse error = errParse

const ErrorInvalidParams ErrorCode = -32602

func NewInvalidParamsError() *Error {
	return &Error{
		Code:    ErrorInvalidParams,
		Message: "invalid params",
	}
}

var errInvalidParams *Error = NewInvalidParamsError()
var ErrInvalidParams error = errInvalidParams

const ErrorInternal ErrorCode = -32603

func NewInternalError() *Error {
	return &Error{
		Code:    ErrorInternal,
		Message: "internal error",
	}
}

var errInternal *Error = NewInternalError()
var ErrInternal error = errInternal

type ErrorCode int32

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details"`
	Data    ErrorData `json:"data,omitempty"`
}

func (e *Error) Error() string {
	if e.Details == "" {
		return "api: " + e.Message + " (" + strconv.FormatInt(int64(e.Code), 10) + ")"
	}
	return "api: " + e.Message + " (" + strconv.FormatInt(int64(e.Code), 10) + "): " + e.Details
}

func (e *Error) Is(err error) bool {
	if e2, ok := err.(*Error); ok && e.Code == e2.Code {
		return true
	}
	return false
}

type ErrorData map[string]interface{}

func (ed *ErrorData) SetValue(key string, value interface{}) {
	if *ed == nil {
		*ed = map[string]interface{}{key: value}
	} else {
		(*ed)[key] = value
	}
}

func (ed ErrorData) GetValue(key string) (interface{}, bool) { value, ok := ed[key]; return value, ok }
func (ed ErrorData) ClearValue(key string)                   { delete(ed, key) }
