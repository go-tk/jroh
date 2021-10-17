// Code generated by jrohc. DO NOT EDIT.

package petstoreapi

import (
	context "context"
	fmt "fmt"
	apicommon "github.com/go-tk/jroh/go/apicommon"
	http "net/http"
)

type UserClient interface {
	CreateUser(ctx context.Context, params *CreateUserParams) (err error)
	GetUser(ctx context.Context, params *GetUserParams) (results *GetUserResults, err error)
	GetUsers(ctx context.Context, params *GetUsersParams) (results *GetUsersResults, err error)
	UpdateUser(ctx context.Context, params *UpdateUserParams) (err error)
}

type userClient struct {
	apicommon.Client

	rpcFiltersTable [4][]apicommon.RPCHandler
	transportTable  [4]http.RoundTripper
}

func NewUserClient(rpcBaseURL string, options apicommon.ClientOptions) UserClient {
	options.Sanitize()
	var c userClient
	c.Init(options.Timeout, rpcBaseURL)
	apicommon.FillRPCFiltersTable(c.rpcFiltersTable[:], options.RPCFilters)
	apicommon.FillTransportTable(c.transportTable[:], options.Transport, options.Middlewares)
	return &c
}

func (c *userClient) CreateUser(ctx context.Context, params *CreateUserParams) error {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      CreateUserParams
	}
	s.Params = *params
	rpcFilters := c.rpcFiltersTable[User_CreateUser]
	s.OutgoingRPC.Init("Petstore", "User", "CreateUser", "Petstore.User.CreateUser", User_CreateUser, &s.Params, nil, rpcFilters)
	transport := c.transportTable[User_CreateUser]
	if err := c.DoRPC(ctx, &s.OutgoingRPC, transport, "/rpc/Petstore.User.CreateUser"); err != nil {
		return fmt.Errorf("rpc failed; fullMethodName=\"Petstore.User.CreateUser\" traceID=%q: %w",
			s.OutgoingRPC.TraceID(), err)
	}
	return nil
}

func (c *userClient) GetUser(ctx context.Context, params *GetUserParams) (*GetUserResults, error) {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      GetUserParams
		Results     GetUserResults
	}
	s.Params = *params
	rpcFilters := c.rpcFiltersTable[User_GetUser]
	s.OutgoingRPC.Init("Petstore", "User", "GetUser", "Petstore.User.GetUser", User_GetUser, &s.Params, &s.Results, rpcFilters)
	transport := c.transportTable[User_GetUser]
	if err := c.DoRPC(ctx, &s.OutgoingRPC, transport, "/rpc/Petstore.User.GetUser"); err != nil {
		return nil, fmt.Errorf("rpc failed; fullMethodName=\"Petstore.User.GetUser\" traceID=%q: %w",
			s.OutgoingRPC.TraceID(), err)
	}
	return &s.Results, nil
}

func (c *userClient) GetUsers(ctx context.Context, params *GetUsersParams) (*GetUsersResults, error) {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      GetUsersParams
		Results     GetUsersResults
	}
	s.Params = *params
	rpcFilters := c.rpcFiltersTable[User_GetUsers]
	s.OutgoingRPC.Init("Petstore", "User", "GetUsers", "Petstore.User.GetUsers", User_GetUsers, &s.Params, &s.Results, rpcFilters)
	transport := c.transportTable[User_GetUsers]
	if err := c.DoRPC(ctx, &s.OutgoingRPC, transport, "/rpc/Petstore.User.GetUsers"); err != nil {
		return nil, fmt.Errorf("rpc failed; fullMethodName=\"Petstore.User.GetUsers\" traceID=%q: %w",
			s.OutgoingRPC.TraceID(), err)
	}
	return &s.Results, nil
}

func (c *userClient) UpdateUser(ctx context.Context, params *UpdateUserParams) error {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      UpdateUserParams
	}
	s.Params = *params
	rpcFilters := c.rpcFiltersTable[User_UpdateUser]
	s.OutgoingRPC.Init("Petstore", "User", "UpdateUser", "Petstore.User.UpdateUser", User_UpdateUser, &s.Params, nil, rpcFilters)
	transport := c.transportTable[User_UpdateUser]
	if err := c.DoRPC(ctx, &s.OutgoingRPC, transport, "/rpc/Petstore.User.UpdateUser"); err != nil {
		return fmt.Errorf("rpc failed; fullMethodName=\"Petstore.User.UpdateUser\" traceID=%q: %w",
			s.OutgoingRPC.TraceID(), err)
	}
	return nil
}

type UserClientFuncs struct {
	CreateUserFunc func(context.Context, *CreateUserParams) error
	GetUserFunc    func(context.Context, *GetUserParams) (*GetUserResults, error)
	GetUsersFunc   func(context.Context, *GetUsersParams) (*GetUsersResults, error)
	UpdateUserFunc func(context.Context, *UpdateUserParams) error
}

var _ UserClient = (*UserClientFuncs)(nil)

func (cf *UserClientFuncs) CreateUser(ctx context.Context, params *CreateUserParams) error {
	f := cf.CreateUserFunc
	if f == nil {
		return apicommon.ErrNotImplemented
	}
	return f(ctx, params)
}

func (cf *UserClientFuncs) GetUser(ctx context.Context, params *GetUserParams) (*GetUserResults, error) {
	f := cf.GetUserFunc
	if f == nil {
		return nil, apicommon.ErrNotImplemented
	}
	return f(ctx, params)
}

func (cf *UserClientFuncs) GetUsers(ctx context.Context, params *GetUsersParams) (*GetUsersResults, error) {
	f := cf.GetUsersFunc
	if f == nil {
		return nil, apicommon.ErrNotImplemented
	}
	return f(ctx, params)
}

func (cf *UserClientFuncs) UpdateUser(ctx context.Context, params *UpdateUserParams) error {
	f := cf.UpdateUserFunc
	if f == nil {
		return apicommon.ErrNotImplemented
	}
	return f(ctx, params)
}
