// Code generated by jrohc. DO NOT EDIT.

package barapi

import (
	context "context"
	fmt "fmt"
	apicommon "github.com/go-tk/jroh/go/apicommon"
)

type TestClient interface {
	DoSomething(ctx context.Context) (err error)
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
	s.rpc.Namespace = "Bar"
	s.rpc.ServiceName = "Test"
	s.rpc.MethodName = "DoSomething"
	s.rpc.FullMethodName = "Bar.Test.DoSomething"
	s.rpc.MethodIndex = Test_DoSomething
	s.rpc.Params = &s.params
	s.rpc.Results = &s.results
	if err := c.doRPC(ctx, &s.rpc, "/rpc/Bar.Test.DoSomething"); err != nil {
		return err
	}
	return nil
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
	DoSomethingFunc func(context.Context) error
}

var _ TestClient = (*TestClientFuncs)(nil)

func (cf *TestClientFuncs) DoSomething(ctx context.Context) error {
	if f := cf.DoSomethingFunc; f != nil {
		return f(ctx)
	}
	return apicommon.NewNotImplementedError()
}
