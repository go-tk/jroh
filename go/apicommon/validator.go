package apicommon

import (
	"bytes"
	"unsafe"
)

type Validator interface {
	Validate(validationContext *ValidationContext) (ok bool)
}

type ValidationContext struct {
	path         []string
	errorDetails string
}

func NewValidationContext() *ValidationContext { return new(ValidationContext) }

func (vc *ValidationContext) Enter(pathComponent string) {
	vc.path = append(vc.path, pathComponent)
}

func (vc *ValidationContext) Leave() {
	vc.path = vc.path[:len(vc.path)-1]
}

func (vc *ValidationContext) SetErrorDetails(errorDetails string) {
	n := len(vc.path)
	if n == 0 {
		vc.errorDetails = errorDetails
		return
	}
	if n == 1 {
		vc.errorDetails = vc.path[0] + ": " + errorDetails
		return
	}
	var buffer bytes.Buffer
	buffer.WriteString(vc.path[0])
	for i := 1; i < n; i++ {
		buffer.WriteByte('.')
		buffer.WriteString(vc.path[i])
	}
	buffer.WriteString(": ")
	buffer.WriteString(errorDetails)
	data := buffer.Bytes()
	vc.errorDetails = *(*string)(unsafe.Pointer(&data))
}

func (vc *ValidationContext) ErrorDetails() string { return vc.errorDetails }
