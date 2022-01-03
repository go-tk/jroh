// Code generated by jrohc. DO NOT EDIT.

package fooapi

import (
	context "context"
	fmt "fmt"
	apicommon "github.com/go-tk/jroh/go/apicommon"
)

type TestClient interface {
	DoSomething(ctx context.Context) (err error)
	DoSomething1(ctx context.Context, params *DoSomething1Params) (err error)
	DoSomething2(ctx context.Context) (results *DoSomething2Results, err error)
	DoSomething3(ctx context.Context, params *DoSomething3Params) (results *DoSomething3Results, err error)
}

type testClient struct {
	rpcBaseURL      string
	options         apicommon.ClientOptions
	rpcFiltersTable [NumberOfTestMethods][]apicommon.OutgoingRPCHandler
}

func NewTestClient(rpcBaseURL string, options apicommon.ClientOptions) TestClient {
	var c testClient
	c.rpcBaseURL = rpcBaseURL
	c.options = options
	c.options.Sanitize()
	apicommon.FillOutgoingRPCFiltersTable(c.rpcFiltersTable[:], options.RPCFilters)
	return &c
}

func (c *testClient) DoSomething(ctx context.Context) error {
	var s struct {
		rpc     apicommon.OutgoingRPC
		params  apicommon.DummyModel
		results apicommon.DummyModel
	}
	s.rpc.Namespace = "Foo"
	s.rpc.ServiceName = "Test"
	s.rpc.MethodName = "DoSomething"
	s.rpc.FullMethodName = "Foo.Test.DoSomething"
	s.rpc.MethodIndex = Test_DoSomething
	s.rpc.Params = &s.params
	s.rpc.Results = &s.results
	if err := c.doRPC(ctx, &s.rpc, "/rpc/Foo.Test.DoSomething"); err != nil {
		return err
	}
	return nil
}

func (c *testClient) DoSomething1(ctx context.Context, params *DoSomething1Params) error {
	var s struct {
		rpc     apicommon.OutgoingRPC
		params  DoSomething1Params
		results apicommon.DummyModel
	}
	s.rpc.Namespace = "Foo"
	s.rpc.ServiceName = "Test"
	s.rpc.MethodName = "DoSomething1"
	s.rpc.FullMethodName = "Foo.Test.DoSomething1"
	s.rpc.MethodIndex = Test_DoSomething1
	s.params = *params
	s.rpc.Params = &s.params
	s.rpc.Results = &s.results
	if err := c.doRPC(ctx, &s.rpc, "/rpc/Foo.Test.DoSomething1"); err != nil {
		return err
	}
	return nil
}

func (c *testClient) DoSomething2(ctx context.Context) (*DoSomething2Results, error) {
	var s struct {
		rpc     apicommon.OutgoingRPC
		params  apicommon.DummyModel
		results DoSomething2Results
	}
	s.rpc.Namespace = "Foo"
	s.rpc.ServiceName = "Test"
	s.rpc.MethodName = "DoSomething2"
	s.rpc.FullMethodName = "Foo.Test.DoSomething2"
	s.rpc.MethodIndex = Test_DoSomething2
	s.rpc.Params = &s.params
	s.rpc.Results = &s.results
	if err := c.doRPC(ctx, &s.rpc, "/rpc/Foo.Test.DoSomething2"); err != nil {
		return nil, err
	}
	return &s.results, nil
}

func (c *testClient) DoSomething3(ctx context.Context, params *DoSomething3Params) (*DoSomething3Results, error) {
	var s struct {
		rpc     apicommon.OutgoingRPC
		params  DoSomething3Params
		results DoSomething3Results
	}
	s.rpc.Namespace = "Foo"
	s.rpc.ServiceName = "Test"
	s.rpc.MethodName = "DoSomething3"
	s.rpc.FullMethodName = "Foo.Test.DoSomething3"
	s.rpc.MethodIndex = Test_DoSomething3
	s.params = *params
	s.rpc.Params = &s.params
	s.rpc.Results = &s.results
	if err := c.doRPC(ctx, &s.rpc, "/rpc/Foo.Test.DoSomething3"); err != nil {
		return nil, err
	}
	return &s.results, nil
}

func (c *testClient) doRPC(ctx context.Context, rpc *apicommon.OutgoingRPC, rpcPath string) error {
	if timeout := c.options.Timeout; timeout >= 1 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	rpc.Transport = c.options.Transport
	rpc.URL = c.rpcBaseURL + rpcPath
	rpc.SetHandler(apicommon.HandleOutgoingRPC)
	rpc.SetFilters(c.rpcFiltersTable[rpc.MethodIndex])
	if err := rpc.Do(ctx); err != nil {
		return fmt.Errorf("rpc failed; fullMethodName=%q traceID=%q: %w", rpc.FullMethodName, rpc.TraceID, err)
	}
	return nil
}

type TestClientFuncs struct {
	DoSomethingFunc  func(context.Context) error
	DoSomething1Func func(context.Context, *DoSomething1Params) error
	DoSomething2Func func(context.Context) (*DoSomething2Results, error)
	DoSomething3Func func(context.Context, *DoSomething3Params) (*DoSomething3Results, error)
}

var _ TestClient = (*TestClientFuncs)(nil)

func (cf *TestClientFuncs) DoSomething(ctx context.Context) error {
	if f := cf.DoSomethingFunc; f != nil {
		return f(ctx)
	}
	return apicommon.NewNotImplementedError()
}

func (cf *TestClientFuncs) DoSomething1(ctx context.Context, params *DoSomething1Params) error {
	if f := cf.DoSomething1Func; f != nil {
		return f(ctx, params)
	}
	return apicommon.NewNotImplementedError()
}

func (cf *TestClientFuncs) DoSomething2(ctx context.Context) (*DoSomething2Results, error) {
	if f := cf.DoSomething2Func; f != nil {
		return f(ctx)
	}
	return nil, apicommon.NewNotImplementedError()
}

func (cf *TestClientFuncs) DoSomething3(ctx context.Context, params *DoSomething3Params) (*DoSomething3Results, error) {
	if f := cf.DoSomething3Func; f != nil {
		return f(ctx, params)
	}
	return nil, apicommon.NewNotImplementedError()
}
