package apicommon

import "context"

type Model interface {
	Validate(validationContext *ValidationContext) (ok bool)
}

type FurtherValidator interface {
	FurtherValidate(validationContext *ValidationContext) (ok bool)
}

type ValidationContext struct {
	Ctx context.Context

	path         []string
	errorDetails string
}

func NewValidationContext(ctx context.Context) *ValidationContext {
	return &ValidationContext{Ctx: ctx}
}

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
	i := len(vc.path[0])
	for j := 1; j < n; j++ {
		i += 1
		i += len(vc.path[j])
	}
	i += 2
	i += len(errorDetails)
	data := make([]byte, i)
	i = copy(data, vc.path[0])
	for j := 1; j < n; j++ {
		i += copy(data[i:], ".")
		i += copy(data[i:], vc.path[j])
	}
	i += copy(data[i:], ": ")
	i += copy(data[i:], errorDetails)
	vc.errorDetails = BytesToString(data)
}

func (vc *ValidationContext) ErrorDetails() string { return vc.errorDetails }

type DummyModel struct{}

var _ Model = DummyModel{}

func (DummyModel) Validate(*ValidationContext) bool { return true }

type DummyFurtherValidator struct{ dummyFurtherValidator }

var _ FurtherValidator = DummyFurtherValidator{}

type dummyFurtherValidator struct{}

func (dummyFurtherValidator) FurtherValidate(*ValidationContext) bool { return true }
