package apicommon

import (
	"net/http"
	"reflect"
	"runtime"
)

type Router struct {
	serveMux   *http.ServeMux
	routeInfos []RouteInfo
}

func NewRouter(serveMux *http.ServeMux) *Router {
	var r Router
	if serveMux == nil {
		serveMux = http.NewServeMux()
	}
	r.serveMux = serveMux
	return &r
}

func (r *Router) AddRoute(
	rpcPath string,
	handler http.Handler,
	fullMethodName string,
	serverMiddlewares []ServerMiddleware,
	rpcFilters []RPCHandler,
) {
	r.serveMux.Handle(rpcPath, handler)
	r.addRouteInfo(rpcPath, fullMethodName, serverMiddlewares, rpcFilters)
}

func (r *Router) addRouteInfo(
	rpcPath string,
	fullMethodName string,
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
	routeInfo := RouteInfo{
		RPCPath:           rpcPath,
		FullMethodName:    fullMethodName,
		ServerMiddlewares: serverMiddlewares2,
		RPCFilters:        rpcFilters2,
	}
	r.routeInfos = append(r.routeInfos, routeInfo)
}

func (r *Router) ServeMux() *http.ServeMux { return r.serveMux }
func (r *Router) RouteInfos() []RouteInfo  { return r.routeInfos }

type RouteInfo struct {
	RPCPath           string
	FullMethodName    string
	ServerMiddlewares []string
	RPCFilters        []string
}
