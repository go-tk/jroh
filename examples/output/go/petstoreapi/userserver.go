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

func RegisterUserServer(server UserServer, router *apicommon.Router, options apicommon.ServerOptions) {
	options.Sanitize()
	var rpcFiltersTable [NumberOfUserMethods][]apicommon.IncomingRPCHandler
	apicommon.FillIncomingRPCFiltersTable(rpcFiltersTable[:], options.RPCFilters)
	{
		rpcFilters := rpcFiltersTable[User_CreateUser]
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var s struct {
				rpc     apicommon.IncomingRPC
				params  CreateUserParams
				results apicommon.DummyModel
			}
			s.rpc.Namespace = "Petstore"
			s.rpc.ServiceName = "User"
			s.rpc.MethodName = "CreateUser"
			s.rpc.FullMethodName = "Petstore.User.CreateUser"
			s.rpc.MethodIndex = User_CreateUser
			s.rpc.Params = &s.params
			s.rpc.Results = &s.results
			s.rpc.SetHandler(func(ctx context.Context, rpc *apicommon.IncomingRPC) error {
				return server.CreateUser(ctx, rpc.Params.(*CreateUserParams))
			})
			s.rpc.SetFilters(rpcFilters)
			apicommon.HandleRequest(r, &s.rpc, options.TraceIDGenerator, w)
		})
		router.AddRoute("/rpc/Petstore.User.CreateUser", handler, "Petstore.User.CreateUser", rpcFilters)
	}
	{
		rpcFilters := rpcFiltersTable[User_GetUser]
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var s struct {
				rpc     apicommon.IncomingRPC
				params  GetUserParams
				results GetUserResults
			}
			s.rpc.Namespace = "Petstore"
			s.rpc.ServiceName = "User"
			s.rpc.MethodName = "GetUser"
			s.rpc.FullMethodName = "Petstore.User.GetUser"
			s.rpc.MethodIndex = User_GetUser
			s.rpc.Params = &s.params
			s.rpc.Results = &s.results
			s.rpc.SetHandler(func(ctx context.Context, rpc *apicommon.IncomingRPC) error {
				return server.GetUser(ctx, rpc.Params.(*GetUserParams), rpc.Results.(*GetUserResults))
			})
			s.rpc.SetFilters(rpcFilters)
			apicommon.HandleRequest(r, &s.rpc, options.TraceIDGenerator, w)
		})
		router.AddRoute("/rpc/Petstore.User.GetUser", handler, "Petstore.User.GetUser", rpcFilters)
	}
	{
		rpcFilters := rpcFiltersTable[User_GetUsers]
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var s struct {
				rpc     apicommon.IncomingRPC
				params  GetUsersParams
				results GetUsersResults
			}
			s.rpc.Namespace = "Petstore"
			s.rpc.ServiceName = "User"
			s.rpc.MethodName = "GetUsers"
			s.rpc.FullMethodName = "Petstore.User.GetUsers"
			s.rpc.MethodIndex = User_GetUsers
			s.rpc.Params = &s.params
			s.rpc.Results = &s.results
			s.rpc.SetHandler(func(ctx context.Context, rpc *apicommon.IncomingRPC) error {
				return server.GetUsers(ctx, rpc.Params.(*GetUsersParams), rpc.Results.(*GetUsersResults))
			})
			s.rpc.SetFilters(rpcFilters)
			apicommon.HandleRequest(r, &s.rpc, options.TraceIDGenerator, w)
		})
		router.AddRoute("/rpc/Petstore.User.GetUsers", handler, "Petstore.User.GetUsers", rpcFilters)
	}
	{
		rpcFilters := rpcFiltersTable[User_UpdateUser]
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var s struct {
				rpc     apicommon.IncomingRPC
				params  UpdateUserParams
				results apicommon.DummyModel
			}
			s.rpc.Namespace = "Petstore"
			s.rpc.ServiceName = "User"
			s.rpc.MethodName = "UpdateUser"
			s.rpc.FullMethodName = "Petstore.User.UpdateUser"
			s.rpc.MethodIndex = User_UpdateUser
			s.rpc.Params = &s.params
			s.rpc.Results = &s.results
			s.rpc.SetHandler(func(ctx context.Context, rpc *apicommon.IncomingRPC) error {
				return server.UpdateUser(ctx, rpc.Params.(*UpdateUserParams))
			})
			s.rpc.SetFilters(rpcFilters)
			apicommon.HandleRequest(r, &s.rpc, options.TraceIDGenerator, w)
		})
		router.AddRoute("/rpc/Petstore.User.UpdateUser", handler, "Petstore.User.UpdateUser", rpcFilters)
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
	return apicommon.NewNotImplementedError()
}

func (sf *UserServerFuncs) GetUser(ctx context.Context, params *GetUserParams, results *GetUserResults) error {
	if f := sf.GetUserFunc; f != nil {
		return f(ctx, params, results)
	}
	return apicommon.NewNotImplementedError()
}

func (sf *UserServerFuncs) GetUsers(ctx context.Context, params *GetUsersParams, results *GetUsersResults) error {
	if f := sf.GetUsersFunc; f != nil {
		return f(ctx, params, results)
	}
	return apicommon.NewNotImplementedError()
}

func (sf *UserServerFuncs) UpdateUser(ctx context.Context, params *UpdateUserParams) error {
	if f := sf.UpdateUserFunc; f != nil {
		return f(ctx, params)
	}
	return apicommon.NewNotImplementedError()
}
