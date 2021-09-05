// Code generated by jrohc. DO NOT EDIT.

package barapi

import (
	context "context"
	apicommon "github.com/go-tk/jroh/go/apicommon"
	http "net/http"
)

type TestServer interface {
	DoSomething(ctx context.Context, results *DoSomethingResults) (err error)
}

func RegisterTestServer(server TestServer, serveMux *http.ServeMux, serverOptions apicommon.ServerOptions) {
	serverOptions.Sanitize()
	var middlewareTable [1][]apicommon.Middleware
	apicommon.FillMiddlewareTable(middlewareTable[:], serverOptions.Middlewares)
	var rpcInterceptorTable [1][]apicommon.RPCHandler
	apicommon.FillRPCInterceptorTable(rpcInterceptorTable[:], serverOptions.RPCInterceptors)
	{
		middlewares := middlewareTable[Test_DoSomething]
		rpcInterceptors := rpcInterceptorTable[Test_DoSomething]
		incomingRPCFactory := func() *apicommon.IncomingRPC {
			var s struct {
				IncomingRPC apicommon.IncomingRPC
				Results     DoSomethingResults
			}
			rpcHandler := func(ctx context.Context, rpc *apicommon.RPC) error {
				return server.DoSomething(ctx, rpc.Results().(*DoSomethingResults))
			}
			s.IncomingRPC.Init(
				"Bar",
				"Test",
				"DoSomething",
				nil,
				&s.Results,
				rpcHandler,
				rpcInterceptors,
			)
			return &s.IncomingRPC
		}
		handler := apicommon.MakeHandler(middlewares, incomingRPCFactory, serverOptions.TraceIDGenerator)
		serveMux.Handle("/rpc/Bar.Test.DoSomething", handler)
	}
}

type TestServerFuncs struct {
	DoSomethingFunc func(context.Context, *DoSomethingResults) error
}

var _ TestServer = (*TestServerFuncs)(nil)

func (sf *TestServerFuncs) DoSomething(ctx context.Context, results *DoSomethingResults) error {
	if f := sf.DoSomethingFunc; f != nil {
		return f(ctx, results)
	}
	return nil
}
