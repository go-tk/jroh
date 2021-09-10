// Code generated by jrohc. DO NOT EDIT.

package fooapi

import (
	context "context"
	apicommon "github.com/go-tk/jroh/go/apicommon"
)

type TestClient interface {
	DoSomething(ctx context.Context, params *DoSomethingParams) (err error)
	DoSomething2(ctx context.Context, params *DoSomething2Params) (results *DoSomething2Results, err error)
	DoSomething3(ctx context.Context) (err error)
}

type testClient struct {
	apicommon.Client

	rpcInterceptorTable [3][]apicommon.RPCHandler
}

func NewTestClient(rpcBaseURL string, options apicommon.ClientOptions) TestClient {
	var c testClient
	c.Init(rpcBaseURL, options)
	apicommon.FillRPCInterceptorTable(c.rpcInterceptorTable[:], options.RPCInterceptors)
	return &c
}

func (c *testClient) DoSomething(ctx context.Context, params *DoSomethingParams) error {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      DoSomethingParams
	}
	s.Params = *params
	rpcInterceptors := c.rpcInterceptorTable[Test_DoSomething]
	s.OutgoingRPC.Init("Foo", "Test", "DoSomething", &s.Params, nil, apicommon.HandleRPC, rpcInterceptors)
	return c.DoRPC(ctx, &s.OutgoingRPC, "/rpc/Foo.Test.DoSomething")
}

func (c *testClient) DoSomething2(ctx context.Context, params *DoSomething2Params) (*DoSomething2Results, error) {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      DoSomething2Params
		Results     DoSomething2Results
	}
	s.Params = *params
	rpcInterceptors := c.rpcInterceptorTable[Test_DoSomething2]
	s.OutgoingRPC.Init("Foo", "Test", "DoSomething2", &s.Params, &s.Results, apicommon.HandleRPC, rpcInterceptors)
	if err := c.DoRPC(ctx, &s.OutgoingRPC, "/rpc/Foo.Test.DoSomething2"); err != nil {
		return nil, err
	}
	return &s.Results, nil
}

func (c *testClient) DoSomething3(ctx context.Context) error {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
	}
	rpcInterceptors := c.rpcInterceptorTable[Test_DoSomething3]
	s.OutgoingRPC.Init("Foo", "Test", "DoSomething3", nil, nil, apicommon.HandleRPC, rpcInterceptors)
	return c.DoRPC(ctx, &s.OutgoingRPC, "/rpc/Foo.Test.DoSomething3")
}

type TestClientFuncs struct {
	DoSomethingFunc  func(context.Context, *DoSomethingParams) error
	DoSomething2Func func(context.Context, *DoSomething2Params) (*DoSomething2Results, error)
	DoSomething3Func func(context.Context) error
}

var _ TestClient = (*TestClientFuncs)(nil)

func (cf *TestClientFuncs) DoSomething(ctx context.Context, params *DoSomethingParams) error {
	if f := cf.DoSomethingFunc; f != nil {
		return f(ctx, params)
	}
	return nil
}

func (cf *TestClientFuncs) DoSomething2(ctx context.Context, params *DoSomething2Params) (*DoSomething2Results, error) {
	if f := cf.DoSomething2Func; f != nil {
		return f(ctx, params)
	}
	return &DoSomething2Results{}, nil
}

func (cf *TestClientFuncs) DoSomething3(ctx context.Context) error {
	if f := cf.DoSomething3Func; f != nil {
		return f(ctx)
	}
	return nil
}
