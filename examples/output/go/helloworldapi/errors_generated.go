// Code generated by jrohc. DO NOT EDIT.

package helloworldapi

import (
	apicommon "github.com/go-tk/jroh/go/apicommon"
)

const ErrorUserNotAllowed apicommon.ErrorCode = 1001

func NewUserNotAllowedError() *apicommon.Error {
	return &apicommon.Error{
		Code:       ErrorUserNotAllowed,
		StatusCode: 403,
		Message:    "user not allowed",
	}
}