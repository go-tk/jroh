package zerologmw

type OptionsBuilder func(options *options)

type options struct {
	MaxParamsSize  int
	MaxResultsSize int
}

func (o *options) Init() *options {
	o.MaxParamsSize = 500
	o.MaxResultsSize = 500
	return o
}

func MaxParamsSize(value int) OptionsBuilder {
	return func(options *options) { options.MaxParamsSize = value }
}

func MaxResultsSize(value int) OptionsBuilder {
	return func(options *options) { options.MaxResultsSize = value }
}
