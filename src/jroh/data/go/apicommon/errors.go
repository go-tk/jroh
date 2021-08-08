package apicommon

import "fmt"

const ErrorParse ErrorCode = -32700

var ErrParse = &Error{
	Code:    ErrorParse,
	Message: "parse error",
}

const ErrorInvalidParams ErrorCode = -32602

var ErrInvalidParams = &Error{
	Code:    ErrorInvalidParams,
	Message: "invalid params",
}

const ErrorInternal ErrorCode = -32603

var ErrInternal = &Error{
	Code:    ErrorInternal,
	Message: "internal error",
}

type ErrorCode int32

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Data    ErrorData `json:"data,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("api: %s (%d)", e.Message, e.Code)
}

func (e *Error) Is(err error) bool {
	if e2, ok := err.(*Error); ok && e.Code == e2.Code {
		return true
	}
	return false
}

func (e *Error) WithData(data ErrorData) *Error {
	e2 := *e
	e2.Data = data
	return &e2
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
