package apicommon

import (
	"net/http"
	"reflect"
	"runtime"
)

type RPCRouter struct {
	serveMux      *http.ServeMux
	rpcRouteInfos []RPCRouteInfo
}

func NewRPCRouter(serveMux *http.ServeMux) *RPCRouter {
	var rr RPCRouter
	if serveMux == nil {
		serveMux = http.NewServeMux()
	}
	rr.serveMux = serveMux
	return &rr
}

func (rr *RPCRouter) AddRPCRoute(
	rpcPath string,
	handler http.Handler,
	fullMethodName string,
	serverMiddlewares []ServerMiddleware,
	rpcFilters []RPCHandler,
) {
	rr.serveMux.Handle(rpcPath, handler)
	rr.addRPCRouteInfo(fullMethodName, rpcPath, serverMiddlewares, rpcFilters)
}

func (rr *RPCRouter) addRPCRouteInfo(
	fullMethodName string,
	rpcPath string,
	serverMiddlewares []ServerMiddleware,
	rpcFilters []RPCHandler,
) {
	serverMiddlewares2 := make([]string, len(serverMiddlewares))
	for i, serverMiddleware := range serverMiddlewares {
		serverMiddlewares2[i] = runtime.FuncForPC(reflect.ValueOf(serverMiddleware).Pointer()).Name()
	}
	rpcFilters2 := make([]string, len(rpcFilters))
	for i, rpcFilter := range rpcFilters {
		rpcFilters2[i] = runtime.FuncForPC(reflect.ValueOf(rpcFilter).Pointer()).Name()
	}
	rpcRouteInfo := RPCRouteInfo{
		FullMethodName:    fullMethodName,
		RPCPath:           rpcPath,
		ServerMiddlewares: serverMiddlewares2,
		RPCFilters:        rpcFilters2,
	}
	rr.rpcRouteInfos = append(rr.rpcRouteInfos, rpcRouteInfo)
}

func (rr *RPCRouter) ServeMux() *http.ServeMux      { return rr.serveMux }
func (rr *RPCRouter) RPCRouteInfos() []RPCRouteInfo { return rr.rpcRouteInfos }

type RPCRouteInfo struct {
	FullMethodName    string
	RPCPath           string
	ServerMiddlewares []string
	RPCFilters        []string
}
