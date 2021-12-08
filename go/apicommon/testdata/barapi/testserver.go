// Code generated by jrohc. DO NOT EDIT.

package barapi

import (
	context "context"
	apicommon "github.com/go-tk/jroh/go/apicommon"
	http "net/http"
)

type TestServer interface {
	DoSomething(ctx context.Context) (err error)
}

func RegisterTestServer(server TestServer, router *apicommon.Router, options apicommon.ServerOptions) {
	options.Sanitize()
	var rpcFiltersTable [NumberOfTestMethods][]apicommon.IncomingRPCHandler
	apicommon.FillIncomingRPCFiltersTable(rpcFiltersTable[:], options.RPCFilters)
	{
		rpcFilters := rpcFiltersTable[Test_DoSomething]
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var s struct {
				rpc     apicommon.IncomingRPC
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
			s.rpc.SetHandler(func(ctx context.Context, rpc *apicommon.IncomingRPC) error {
				return server.DoSomething(ctx)
			})
			s.rpc.SetFilters(rpcFilters)
			apicommon.HandleRequest(r, &s.rpc, options.TraceIDGenerator, w)
		})
		router.AddRoute("/rpc/Bar.Test.DoSomething", handler, "Bar.Test.DoSomething", rpcFilters)
	}
}

type TestServerFuncs struct {
	DoSomethingFunc func(context.Context) error
}

var _ TestServer = (*TestServerFuncs)(nil)

func (sf *TestServerFuncs) DoSomething(ctx context.Context) error {
	if f := sf.DoSomethingFunc; f != nil {
		return f(ctx)
	}
	return apicommon.NewNotImplementedError()
}
