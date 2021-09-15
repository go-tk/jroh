package fooapi

import "github.com/go-tk/jroh/go/apicommon"

var _ apicommon.FurtherValidator = XString("")

func (m XString) FurtherValidate(validationContext *apicommon.ValidationContext) bool {
	if m == "taboo" {
		validationContext.SetErrorDetails("this is taboo!")
		return false
	}
	return true
}

var _ apicommon.FurtherValidator = (*MyStructInt32)(nil)

func (m *MyStructInt32) FurtherValidate(validationContext *apicommon.ValidationContext) bool {
	if m.TheInt32A == 666666 {
		validationContext.SetErrorDetails("theInt32A is evil!")
		return false
	}
	return true
}
