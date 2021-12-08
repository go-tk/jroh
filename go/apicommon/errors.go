package apicommon

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Code       ErrorCode `json:"-"`
	StatusCode int       `json:"-"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	Data       ErrorData `json:"data,omitempty"`
}

var _ error = (*Error)(nil)

func (e *Error) Error() string {
	if e.Code < 0 {
		if e.Details == "" {
			return fmt.Sprintf("http: %v", e.Message)
		}
		return fmt.Sprintf("http: %v: %v", e.Message, e.Details)
	}
	if e.Details == "" {
		return fmt.Sprintf("api: %v", e.Message)
	}
	return fmt.Sprintf("api: %v: %v", e.Message, e.Details)
}

func (e *Error) Temporary() bool { return e.StatusCode == http.StatusServiceUnavailable }

type ErrorCode int32

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

const ErrorNotImplemented ErrorCode = 1

func NewNotImplementedError() *Error {
	return &Error{
		Code:       ErrorNotImplemented,
		StatusCode: http.StatusNotImplemented,
		Message:    "not implemented",
	}
}

const ErrorInvalidParams ErrorCode = 2

func NewInvalidParamsError() *Error {
	return &Error{
		Code:       ErrorInvalidParams,
		StatusCode: http.StatusUnprocessableEntity,
		Message:    "invalid params",
	}
}

var (
	ErrInvalidResults       = errors.New("invalid results")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)
