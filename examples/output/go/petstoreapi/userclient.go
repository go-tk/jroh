// Code generated by jrohc. DO NOT EDIT.

package petstoreapi

import (
	context "context"
	apicommon "github.com/go-tk/jroh/go/apicommon"
)

type UserClient interface {
	CreateUser(ctx context.Context, params *CreateUserParams) (err error)
	GetUser(ctx context.Context, params *GetUserParams) (results *GetUserResults, err error)
	GetUsers(ctx context.Context, params *GetUsersParams) (results *GetUsersResults, err error)
	UpdateUser(ctx context.Context, params *UpdateUserParams) (err error)
}

type userClient struct {
	apicommon.Client

	rpcInterceptorTable [4][]apicommon.RPCHandler
}

func NewUserClient(rpcBaseURL string, options apicommon.ClientOptions) UserClient {
	var c userClient
	c.Init(rpcBaseURL, options)
	apicommon.FillRPCInterceptorTable(c.rpcInterceptorTable[:], options.RPCInterceptors)
	return &c
}

func (c *userClient) CreateUser(ctx context.Context, params *CreateUserParams) error {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      CreateUserParams
	}
	s.Params = *params
	rpcInterceptors := c.rpcInterceptorTable[User_CreateUser]
	s.OutgoingRPC.Init("Petstore", "User", "CreateUser", &s.Params, nil, apicommon.HandleRPC, rpcInterceptors)
	return c.DoRPC(ctx, &s.OutgoingRPC, "/rpc/Petstore.User.CreateUser")
}

func (c *userClient) GetUser(ctx context.Context, params *GetUserParams) (*GetUserResults, error) {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      GetUserParams
		Results     GetUserResults
	}
	s.Params = *params
	rpcInterceptors := c.rpcInterceptorTable[User_GetUser]
	s.OutgoingRPC.Init("Petstore", "User", "GetUser", &s.Params, &s.Results, apicommon.HandleRPC, rpcInterceptors)
	if err := c.DoRPC(ctx, &s.OutgoingRPC, "/rpc/Petstore.User.GetUser"); err != nil {
		return nil, err
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
	rpcInterceptors := c.rpcInterceptorTable[User_GetUsers]
	s.OutgoingRPC.Init("Petstore", "User", "GetUsers", &s.Params, &s.Results, apicommon.HandleRPC, rpcInterceptors)
	if err := c.DoRPC(ctx, &s.OutgoingRPC, "/rpc/Petstore.User.GetUsers"); err != nil {
		return nil, err
	}
	return &s.Results, nil
}

func (c *userClient) UpdateUser(ctx context.Context, params *UpdateUserParams) error {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      UpdateUserParams
	}
	s.Params = *params
	rpcInterceptors := c.rpcInterceptorTable[User_UpdateUser]
	s.OutgoingRPC.Init("Petstore", "User", "UpdateUser", &s.Params, nil, apicommon.HandleRPC, rpcInterceptors)
	return c.DoRPC(ctx, &s.OutgoingRPC, "/rpc/Petstore.User.UpdateUser")
}

type UserClientFuncs struct {
	CreateUserFunc func(context.Context, *CreateUserParams) error
	GetUserFunc    func(context.Context, *GetUserParams) (*GetUserResults, error)
	GetUsersFunc   func(context.Context, *GetUsersParams) (*GetUsersResults, error)
	UpdateUserFunc func(context.Context, *UpdateUserParams) error
}

var _ UserClient = (*UserClientFuncs)(nil)

func (cf *UserClientFuncs) CreateUser(ctx context.Context, params *CreateUserParams) error {
	if f := cf.CreateUserFunc; f != nil {
		return f(ctx, params)
	}
	return nil
}

func (cf *UserClientFuncs) GetUser(ctx context.Context, params *GetUserParams) (*GetUserResults, error) {
	if f := cf.GetUserFunc; f != nil {
		return f(ctx, params)
	}
	return &GetUserResults{}, nil
}

func (cf *UserClientFuncs) GetUsers(ctx context.Context, params *GetUsersParams) (*GetUsersResults, error) {
	if f := cf.GetUsersFunc; f != nil {
		return f(ctx, params)
	}
	return &GetUsersResults{}, nil
}

func (cf *UserClientFuncs) UpdateUser(ctx context.Context, params *UpdateUserParams) error {
	if f := cf.UpdateUserFunc; f != nil {
		return f(ctx, params)
	}
	return nil
}