// Code generated by jrohc. DO NOT EDIT.

package helloworldapi

import (
	context "context"
	fmt "fmt"
	apicommon "github.com/go-tk/jroh/go/apicommon"
)

type GreeterClient interface {
	SayHello(ctx context.Context, params *SayHelloParams) (results *SayHelloResults, err error)
}

type greeterClient struct {
	rpcBaseURL      string
	options         apicommon.ClientOptions
	rpcFiltersTable [NumberOfGreeterMethods][]apicommon.OutgoingRPCHandler
}

func NewGreeterClient(rpcBaseURL string, options apicommon.ClientOptions) GreeterClient {
	var c greeterClient
	c.rpcBaseURL = rpcBaseURL
	c.options = options
	c.options.Sanitize()
	apicommon.FillOutgoingRPCFiltersTable(c.rpcFiltersTable[:], options.RPCFilters)
	return &c
}

func (c *greeterClient) SayHello(ctx context.Context, params *SayHelloParams) (*SayHelloResults, error) {
	var s struct {
		rpc     apicommon.OutgoingRPC
		params  SayHelloParams
		results SayHelloResults
	}
	s.rpc.Namespace = "HelloWorld"
	s.rpc.ServiceName = "Greeter"
	s.rpc.MethodName = "SayHello"
	s.rpc.FullMethodName = "HelloWorld.Greeter.SayHello"
	s.rpc.MethodIndex = Greeter_SayHello
	s.params = *params
	s.rpc.Params = &s.params
	s.rpc.Results = &s.results
	if err := c.doRPC(ctx, &s.rpc, "/rpc/HelloWorld.Greeter.SayHello"); err != nil {
		return nil, fmt.Errorf("do rpc; fullMethodName=\"HelloWorld.Greeter.SayHello\" traceID=%q: %w", s.rpc.TraceID, err)
	}
	return &s.results, nil
}

func (c *greeterClient) doRPC(ctx context.Context, rpc *apicommon.OutgoingRPC, rpcPath string) error {
	if timeout := c.options.Timeout; timeout >= 1 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	rpc.Transport = c.options.Transport
	rpc.URL = c.rpcBaseURL + rpcPath
	rpc.SetHandler(apicommon.HandleOutgoingRPC)
	rpc.SetFilters(c.rpcFiltersTable[rpc.MethodIndex])
	return rpc.Do(ctx)
}

type GreeterClientFuncs struct {
	SayHelloFunc func(context.Context, *SayHelloParams) (*SayHelloResults, error)
}

var _ GreeterClient = (*GreeterClientFuncs)(nil)

func (cf *GreeterClientFuncs) SayHello(ctx context.Context, params *SayHelloParams) (*SayHelloResults, error) {
	if f := cf.SayHelloFunc; f != nil {
		return f(ctx, params)
	}
	err := apicommon.NewNotImplementedError()
	return nil, fmt.Errorf("do rpc; fullMethodName=\"HelloWorld.Greeter.SayHello\": %w", err)
}
