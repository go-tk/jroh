// Code generated by jrohc. DO NOT EDIT.

package petstoreapi

import (
	context "context"
	fmt "fmt"
	apicommon "github.com/go-tk/jroh/go/apicommon"
	http "net/http"
)

type StoreClient interface {
	CreateOrder(ctx context.Context, params *CreateOrderParams) (results *CreateOrderResults, err error)
	GetOrder(ctx context.Context, params *GetOrderParams) (results *GetOrderResults, err error)
}

type storeClient struct {
	apicommon.Client

	rpcFiltersTable [NumberOfStoreMethods][]apicommon.RPCHandler
	transportTable  [NumberOfStoreMethods]http.RoundTripper
}

func NewStoreClient(rpcBaseURL string, options apicommon.ClientOptions) StoreClient {
	options.Sanitize()
	var c storeClient
	c.Init(options.Timeout, rpcBaseURL)
	apicommon.FillRPCFiltersTable(c.rpcFiltersTable[:], options.RPCFilters)
	apicommon.FillTransportTable(c.transportTable[:], options.Transport, options.Middlewares)
	return &c
}

func (c *storeClient) CreateOrder(ctx context.Context, params *CreateOrderParams) (*CreateOrderResults, error) {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      CreateOrderParams
		Results     CreateOrderResults
	}
	s.Params = *params
	rpcFilters := c.rpcFiltersTable[Store_CreateOrder]
	s.OutgoingRPC.Init("Petstore", "Store", "CreateOrder", "Petstore.Store.CreateOrder", Store_CreateOrder, &s.Params, &s.Results, rpcFilters)
	transport := c.transportTable[Store_CreateOrder]
	if err := c.DoRPC(ctx, &s.OutgoingRPC, transport, "/rpc/Petstore.Store.CreateOrder"); err != nil {
		return nil, fmt.Errorf("rpc failed; fullMethodName=\"Petstore.Store.CreateOrder\" traceID=%q: %w",
			s.OutgoingRPC.TraceID(), err)
	}
	return &s.Results, nil
}

func (c *storeClient) GetOrder(ctx context.Context, params *GetOrderParams) (*GetOrderResults, error) {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      GetOrderParams
		Results     GetOrderResults
	}
	s.Params = *params
	rpcFilters := c.rpcFiltersTable[Store_GetOrder]
	s.OutgoingRPC.Init("Petstore", "Store", "GetOrder", "Petstore.Store.GetOrder", Store_GetOrder, &s.Params, &s.Results, rpcFilters)
	transport := c.transportTable[Store_GetOrder]
	if err := c.DoRPC(ctx, &s.OutgoingRPC, transport, "/rpc/Petstore.Store.GetOrder"); err != nil {
		return nil, fmt.Errorf("rpc failed; fullMethodName=\"Petstore.Store.GetOrder\" traceID=%q: %w",
			s.OutgoingRPC.TraceID(), err)
	}
	return &s.Results, nil
}

type StoreClientFuncs struct {
	CreateOrderFunc func(context.Context, *CreateOrderParams) (*CreateOrderResults, error)
	GetOrderFunc    func(context.Context, *GetOrderParams) (*GetOrderResults, error)
}

var _ StoreClient = (*StoreClientFuncs)(nil)

func (cf *StoreClientFuncs) CreateOrder(ctx context.Context, params *CreateOrderParams) (*CreateOrderResults, error) {
	if f := cf.CreateOrderFunc; f != nil {
		return f(ctx, params)
	}
	return nil, apicommon.ErrNotImplemented
}

func (cf *StoreClientFuncs) GetOrder(ctx context.Context, params *GetOrderParams) (*GetOrderResults, error) {
	if f := cf.GetOrderFunc; f != nil {
		return f(ctx, params)
	}
	return nil, apicommon.ErrNotImplemented
}
