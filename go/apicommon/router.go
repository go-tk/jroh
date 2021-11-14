package apicommon

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

type Router struct {
	rpcPath2Handler map[string]http.Handler
	routeInfos      []RouteInfo

	NotFoundHandler http.Handler
}

var _ http.Handler = (*Router)(nil)

func NewRouter() *Router {
	var r Router
	r.NotFoundHandler = http.NotFoundHandler()
	return &r
}

func (r *Router) AddRoute(
	rpcPath string,
	handler http.Handler,
	fullMethodName string,
	serverMiddlewares []ServerMiddleware,
	rpcFilters []RPCHandler,
) {
	rpcPath2Handler := r.rpcPath2Handler
	if rpcPath2Handler == nil {
		rpcPath2Handler = make(map[string]http.Handler)
		r.rpcPath2Handler = rpcPath2Handler
	}
	if _, ok := rpcPath2Handler[rpcPath]; ok {
		panic(fmt.Sprintf("duplicate route; rpcPath=%q", rpcPath))
	}
	rpcPath2Handler[rpcPath] = handler
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

type RouteInfo struct {
	RPCPath           string
	FullMethodName    string
	ServerMiddlewares []string
	RPCFilters        []string
}

func (r *Router) RouteInfos() []RouteInfo { return r.routeInfos }

func (r *Router) ServeHTTP(w http.ResponseWriter, rr *http.Request) {
	if rr.Method == "POST" {
		if handler, ok := r.rpcPath2Handler[rr.URL.Path]; ok {
			handler.ServeHTTP(w, rr)
			return
		}
	}
	r.NotFoundHandler.ServeHTTP(w, rr)
}
