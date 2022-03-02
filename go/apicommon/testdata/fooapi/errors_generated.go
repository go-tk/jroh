// Code generated by jrohc. DO NOT EDIT.

package fooapi

import (
	apicommon "github.com/go-tk/jroh/go/apicommon"
)

const ErrorSomethingWrong apicommon.ErrorCode = 1000

func NewSomethingWrongError() *apicommon.Error {
	return &apicommon.Error{
		Code:       ErrorSomethingWrong,
		StatusCode: 500,
		Message:    "something wrong",
	}
}