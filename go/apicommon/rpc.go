package apicommon

import "context"

type RPC struct {
	namespace            string
	serviceName          string
	methodName           string
	params               interface{}
	results              interface{}
	handler              RPCHandler
	interceptors         []RPCHandler
	nextInterceptorIndex int
	traceID              string
}

func (r *RPC) Init(
	namespace string,
	serviceName string,
	methodName string,
	params interface{},
	results interface{},
	handler RPCHandler,
	interceptors []RPCHandler,
) {
	r.namespace = namespace
	r.serviceName = serviceName
	r.methodName = methodName
	r.params = params
	r.results = results
	r.handler = handler
	r.interceptors = interceptors
}

func (r *RPC) Namespace() string    { return r.namespace }
func (r *RPC) ServiceName() string  { return r.serviceName }
func (r *RPC) MethodName() string   { return r.methodName }
func (r *RPC) Params() interface{}  { return r.params }
func (r *RPC) Results() interface{} { return r.results }
func (r *RPC) TraceID() string      { return r.traceID }

func (r *RPC) Do(ctx context.Context) error {
	if i := r.nextInterceptorIndex; i < len(r.interceptors) {
		r.nextInterceptorIndex++
		return r.interceptors[i](ctx, r)
	}
	return r.handler(ctx, r)
}

type RPCHandler func(ctx context.Context, rpc *RPC) (err error)

type contextValueRPC struct{}

func makeContextWithRPC(ctx context.Context, rpc *RPC) context.Context {
	return context.WithValue(ctx, contextValueRPC{}, rpc)
}

func MustGetRPCFromContext(ctx context.Context) *RPC {
	return ctx.Value(contextValueRPC{}).(*RPC)
}

func GetRPCFromContext(ctx context.Context) (*RPC, bool) {
	rpc, ok := ctx.Value(contextValueRPC{}).(*RPC)
	return rpc, ok
}
