// Code generated by jrohc. DO NOT EDIT.

package fooapi

import (
	apicommon "github.com/go-tk/jroh/go/apicommon"
)

const ErrorSomethingWrong apicommon.ErrorCode = 1

func NewSomethingWrongError() *apicommon.Error {
	return &apicommon.Error{
		Code:    ErrorSomethingWrong,
		Message: "something wrong",
	}
}

var errSomethingWrong *apicommon.Error = NewSomethingWrongError()
var ErrSomethingWrong error = errSomethingWrong