// Code generated by jrohc. DO NOT EDIT.

package helloworldapi

import (
	context "context"
	apicommon "github.com/go-tk/jroh/go/apicommon"
	http "net/http"
)

type GreeterActor interface {
	SayHello(ctx context.Context, params *SayHelloParams, results *SayHelloResults) (err error)
}

func RegisterGreeterActor(a GreeterActor, router *apicommon.Router, options apicommon.ActorOptions) {
	options.Sanitize()
	var rpcFiltersTable [NumberOfGreeterMethods][]apicommon.IncomingRPCHandler
	apicommon.FillIncomingRPCFiltersTable(rpcFiltersTable[:], options.RPCFilters)
	{
		rpcFilters := rpcFiltersTable[Greeter_SayHello]
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var s struct {
				rpc     apicommon.IncomingRPC
				params  SayHelloParams
				results SayHelloResults
			}
			s.rpc.Namespace = "HelloWorld"
			s.rpc.ServiceName = "Greeter"
			s.rpc.MethodName = "SayHello"
			s.rpc.FullMethodName = "HelloWorld.Greeter.SayHello"
			s.rpc.MethodIndex = Greeter_SayHello
			s.rpc.Params = &s.params
			s.rpc.Results = &s.results
			s.rpc.SetHandler(func(ctx context.Context, rpc *apicommon.IncomingRPC) error {
				return a.SayHello(ctx, rpc.Params.(*SayHelloParams), rpc.Results.(*SayHelloResults))
			})
			s.rpc.SetFilters(rpcFilters)
			apicommon.HandleRequest(r, &s.rpc, options.TraceIDGenerator, w)
		})
		router.AddRoute("/rpc/HelloWorld.Greeter.SayHello", handler, "HelloWorld.Greeter.SayHello", rpcFilters)
	}
}

type GreeterActorFuncs struct {
	SayHelloFunc func(context.Context, *SayHelloParams, *SayHelloResults) error
}

var _ GreeterActor = (*GreeterActorFuncs)(nil)

func (sf *GreeterActorFuncs) SayHello(ctx context.Context, params *SayHelloParams, results *SayHelloResults) error {
	if f := sf.SayHelloFunc; f != nil {
		return f(ctx, params, results)
	}
	return apicommon.NewNotImplementedError()
}
