// Code generated by jrohc. DO NOT EDIT.

package petstoreapi

import (
	context "context"
	apicommon "github.com/go-tk/jroh/go/apicommon"
	http "net/http"
)

type UserServer interface {
	CreateUser(ctx context.Context, params *CreateUserParams) (err error)
	GetUser(ctx context.Context, params *GetUserParams, results *GetUserResults) (err error)
	GetUsers(ctx context.Context, params *GetUsersParams, results *GetUsersResults) (err error)
	UpdateUser(ctx context.Context, params *UpdateUserParams) (err error)
}

func RegisterUserServer(server UserServer, serveMux *http.ServeMux, serverOptions apicommon.ServerOptions) {
	serverOptions.Sanitize()
	var middlewareTable [4][]apicommon.Middleware
	apicommon.FillMiddlewareTable(middlewareTable[:], serverOptions.Middlewares)
	var rpcInterceptorTable [4][]apicommon.RPCHandler
	apicommon.FillRPCInterceptorTable(rpcInterceptorTable[:], serverOptions.RPCInterceptors)
	{
		middlewares := middlewareTable[User_CreateUser]
		rpcInterceptors := rpcInterceptorTable[User_CreateUser]
		incomingRPCFactory := func() *apicommon.IncomingRPC {
			var s struct {
				IncomingRPC apicommon.IncomingRPC
				Params      CreateUserParams
			}
			rpcHandler := func(ctx context.Context, rpc *apicommon.RPC) error {
				return server.CreateUser(ctx, rpc.Params().(*CreateUserParams))
			}
			s.IncomingRPC.Init("Petstore", "User", "CreateUser", &s.Params, nil, rpcHandler, rpcInterceptors)
			return &s.IncomingRPC
		}
		handler := apicommon.MakeHandler(middlewares, incomingRPCFactory, serverOptions.TraceIDGenerator)
		serveMux.Handle("/rpc/Petstore.User.CreateUser", handler)
	}
	{
		middlewares := middlewareTable[User_GetUser]
		rpcInterceptors := rpcInterceptorTable[User_GetUser]
		incomingRPCFactory := func() *apicommon.IncomingRPC {
			var s struct {
				IncomingRPC apicommon.IncomingRPC
				Params      GetUserParams
				Results     GetUserResults
			}
			rpcHandler := func(ctx context.Context, rpc *apicommon.RPC) error {
				return server.GetUser(ctx, rpc.Params().(*GetUserParams), rpc.Results().(*GetUserResults))
			}
			s.IncomingRPC.Init("Petstore", "User", "GetUser", &s.Params, &s.Results, rpcHandler, rpcInterceptors)
			return &s.IncomingRPC
		}
		handler := apicommon.MakeHandler(middlewares, incomingRPCFactory, serverOptions.TraceIDGenerator)
		serveMux.Handle("/rpc/Petstore.User.GetUser", handler)
	}
	{
		middlewares := middlewareTable[User_GetUsers]
		rpcInterceptors := rpcInterceptorTable[User_GetUsers]
		incomingRPCFactory := func() *apicommon.IncomingRPC {
			var s struct {
				IncomingRPC apicommon.IncomingRPC
				Params      GetUsersParams
				Results     GetUsersResults
			}
			rpcHandler := func(ctx context.Context, rpc *apicommon.RPC) error {
				return server.GetUsers(ctx, rpc.Params().(*GetUsersParams), rpc.Results().(*GetUsersResults))
			}
			s.IncomingRPC.Init("Petstore", "User", "GetUsers", &s.Params, &s.Results, rpcHandler, rpcInterceptors)
			return &s.IncomingRPC
		}
		handler := apicommon.MakeHandler(middlewares, incomingRPCFactory, serverOptions.TraceIDGenerator)
		serveMux.Handle("/rpc/Petstore.User.GetUsers", handler)
	}
	{
		middlewares := middlewareTable[User_UpdateUser]
		rpcInterceptors := rpcInterceptorTable[User_UpdateUser]
		incomingRPCFactory := func() *apicommon.IncomingRPC {
			var s struct {
				IncomingRPC apicommon.IncomingRPC
				Params      UpdateUserParams
			}
			rpcHandler := func(ctx context.Context, rpc *apicommon.RPC) error {
				return server.UpdateUser(ctx, rpc.Params().(*UpdateUserParams))
			}
			s.IncomingRPC.Init("Petstore", "User", "UpdateUser", &s.Params, nil, rpcHandler, rpcInterceptors)
			return &s.IncomingRPC
		}
		handler := apicommon.MakeHandler(middlewares, incomingRPCFactory, serverOptions.TraceIDGenerator)
		serveMux.Handle("/rpc/Petstore.User.UpdateUser", handler)
	}
}

type UserServerFuncs struct {
	CreateUserFunc func(context.Context, *CreateUserParams) error
	GetUserFunc    func(context.Context, *GetUserParams, *GetUserResults) error
	GetUsersFunc   func(context.Context, *GetUsersParams, *GetUsersResults) error
	UpdateUserFunc func(context.Context, *UpdateUserParams) error
}

var _ UserServer = (*UserServerFuncs)(nil)

func (sf *UserServerFuncs) CreateUser(ctx context.Context, params *CreateUserParams) error {
	if f := sf.CreateUserFunc; f != nil {
		return f(ctx, params)
	}
	return nil
}

func (sf *UserServerFuncs) GetUser(ctx context.Context, params *GetUserParams, results *GetUserResults) error {
	if f := sf.GetUserFunc; f != nil {
		return f(ctx, params, results)
	}
	return nil
}

func (sf *UserServerFuncs) GetUsers(ctx context.Context, params *GetUsersParams, results *GetUsersResults) error {
	if f := sf.GetUsersFunc; f != nil {
		return f(ctx, params, results)
	}
	return nil
}

func (sf *UserServerFuncs) UpdateUser(ctx context.Context, params *UpdateUserParams) error {
	if f := sf.UpdateUserFunc; f != nil {
		return f(ctx, params)
	}
	return nil
}