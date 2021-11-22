// Code generated by jrohc. DO NOT EDIT.

package fooapi

import (
	context "context"
	apicommon "github.com/go-tk/jroh/go/apicommon"
)

type TestService interface {
	DoSomething(ctx context.Context, params *DoSomethingParams) (err error)
	DoSomething2(ctx context.Context, params *DoSomething2Params, results *DoSomething2Results) (err error)
	DoSomething3(ctx context.Context) (err error)
}

func RegisterTestService(service TestService, router *apicommon.Router, serverOptions apicommon.ServerOptions) {
	serverOptions.Sanitize()
	var serverMiddlewareTable [NumberOfTestMethods][]apicommon.ServerMiddleware
	apicommon.FillServerMiddlewareTable(serverMiddlewareTable[:], serverOptions.Middlewares)
	var rpcFiltersTable [NumberOfTestMethods][]apicommon.RPCHandler
	apicommon.FillRPCFiltersTable(rpcFiltersTable[:], serverOptions.RPCFilters)
	{
		serverMiddlewares := serverMiddlewareTable[Test_DoSomething]
		rpcFilters := rpcFiltersTable[Test_DoSomething]
		incomingRPCFactory := func() *apicommon.IncomingRPC {
			var s struct {
				IncomingRPC apicommon.IncomingRPC
				Params      DoSomethingParams
			}
			rpcHandler := func(ctx context.Context, rpc *apicommon.RPC) error {
				return service.DoSomething(ctx, rpc.Params().(*DoSomethingParams))
			}
			s.IncomingRPC.Init("Foo", "Test", "DoSomething", "Foo.Test.DoSomething", Test_DoSomething, &s.Params, nil, rpcHandler, rpcFilters)
			return &s.IncomingRPC
		}
		handler := apicommon.MakeHandler(serverMiddlewares, incomingRPCFactory, serverOptions.TraceIDGenerator)
		router.AddRoute("/rpc/Foo.Test.DoSomething", handler, "Foo.Test.DoSomething", serverMiddlewares, rpcFilters)
	}
	{
		serverMiddlewares := serverMiddlewareTable[Test_DoSomething2]
		rpcFilters := rpcFiltersTable[Test_DoSomething2]
		incomingRPCFactory := func() *apicommon.IncomingRPC {
			var s struct {
				IncomingRPC apicommon.IncomingRPC
				Params      DoSomething2Params
				Results     DoSomething2Results
			}
			rpcHandler := func(ctx context.Context, rpc *apicommon.RPC) error {
				return service.DoSomething2(ctx, rpc.Params().(*DoSomething2Params), rpc.Results().(*DoSomething2Results))
			}
			s.IncomingRPC.Init("Foo", "Test", "DoSomething2", "Foo.Test.DoSomething2", Test_DoSomething2, &s.Params, &s.Results, rpcHandler, rpcFilters)
			return &s.IncomingRPC
		}
		handler := apicommon.MakeHandler(serverMiddlewares, incomingRPCFactory, serverOptions.TraceIDGenerator)
		router.AddRoute("/rpc/Foo.Test.DoSomething2", handler, "Foo.Test.DoSomething2", serverMiddlewares, rpcFilters)
	}
	{
		serverMiddlewares := serverMiddlewareTable[Test_DoSomething3]
		rpcFilters := rpcFiltersTable[Test_DoSomething3]
		incomingRPCFactory := func() *apicommon.IncomingRPC {
			var s struct {
				IncomingRPC apicommon.IncomingRPC
			}
			rpcHandler := func(ctx context.Context, rpc *apicommon.RPC) error {
				return service.DoSomething3(ctx)
			}
			s.IncomingRPC.Init("Foo", "Test", "DoSomething3", "Foo.Test.DoSomething3", Test_DoSomething3, nil, nil, rpcHandler, rpcFilters)
			return &s.IncomingRPC
		}
		handler := apicommon.MakeHandler(serverMiddlewares, incomingRPCFactory, serverOptions.TraceIDGenerator)
		router.AddRoute("/rpc/Foo.Test.DoSomething3", handler, "Foo.Test.DoSomething3", serverMiddlewares, rpcFilters)
	}
}

type TestServiceFuncs struct {
	DoSomethingFunc  func(context.Context, *DoSomethingParams) error
	DoSomething2Func func(context.Context, *DoSomething2Params, *DoSomething2Results) error
	DoSomething3Func func(context.Context) error
}

var _ TestService = (*TestServiceFuncs)(nil)

func (sf *TestServiceFuncs) DoSomething(ctx context.Context, params *DoSomethingParams) error {
	if f := sf.DoSomethingFunc; f != nil {
		return f(ctx, params)
	}
	return apicommon.ErrNotImplemented
}

func (sf *TestServiceFuncs) DoSomething2(ctx context.Context, params *DoSomething2Params, results *DoSomething2Results) error {
	if f := sf.DoSomething2Func; f != nil {
		return f(ctx, params, results)
	}
	return apicommon.ErrNotImplemented
}

func (sf *TestServiceFuncs) DoSomething3(ctx context.Context) error {
	if f := sf.DoSomething3Func; f != nil {
		return f(ctx)
	}
	return apicommon.ErrNotImplemented
}