// Code generated by jrohc. DO NOT EDIT.

package helloworldapi

import (
	context "context"
	fmt "fmt"
	apicommon "github.com/go-tk/jroh/go/apicommon"
	http "net/http"
)

type GreeterClient interface {
	SayHello(ctx context.Context, params *SayHelloParams) (results *SayHelloResults, err error)
}

type greeterClient struct {
	apicommon.Client

	rpcFiltersTable [1][]apicommon.RPCHandler
	transportTable  [1]http.RoundTripper
}

func NewGreeterClient(rpcBaseURL string, options apicommon.ClientOptions) GreeterClient {
	options.Sanitize()
	var c greeterClient
	c.Init(rpcBaseURL, options.Timeout)
	apicommon.FillRPCFiltersTable(c.rpcFiltersTable[:], options.RPCFilters)
	apicommon.FillTransportTable(c.transportTable[:], options.Transport, options.Middlewares)
	return &c
}

func (c *greeterClient) SayHello(ctx context.Context, params *SayHelloParams) (*SayHelloResults, error) {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      SayHelloParams
		Results     SayHelloResults
	}
	s.Params = *params
	rpcFilters := c.rpcFiltersTable[Greeter_SayHello]
	s.OutgoingRPC.Init("HelloWorld", "Greeter", "SayHello", &s.Params, &s.Results, apicommon.HandleRPC, rpcFilters)
	transport := c.transportTable[Greeter_SayHello]
	if err := c.DoRPC(ctx, &s.OutgoingRPC, transport, "/rpc/HelloWorld.Greeter.SayHello"); err != nil {
		return nil, fmt.Errorf("rpc failed; namespace=\"HelloWorld\" serviceName=\"Greeter\" methodName=\"SayHello\" traceID=%q: %w",
			s.OutgoingRPC.TraceID(), err)
	}
	return &s.Results, nil
}

type GreeterClientFuncs struct {
	SayHelloFunc func(context.Context, *SayHelloParams) (*SayHelloResults, error)
}

var _ GreeterClient = (*GreeterClientFuncs)(nil)

func (cf *GreeterClientFuncs) SayHello(ctx context.Context, params *SayHelloParams) (*SayHelloResults, error) {
	f := cf.SayHelloFunc
	if f == nil {
		return nil, apicommon.ErrNotImplemented
	}
	return f(ctx, params)
}
