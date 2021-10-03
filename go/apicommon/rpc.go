package apicommon

import "context"

type RPC struct {
	mark byte

	namespace       string
	serviceName     string
	methodName      string
	fullMethodName  string
	params          Model
	results         Model
	traceID         string
	handler         RPCHandler
	filters         []RPCHandler
	nextFilterIndex int
}

func (r *RPC) init(
	namespace string,
	serviceName string,
	methodName string,
	fullMethodName string,
	params Model,
	results Model,
	handler RPCHandler,
	filters []RPCHandler,
) {
	r.namespace = namespace
	r.serviceName = serviceName
	r.methodName = methodName
	r.fullMethodName = fullMethodName
	r.params = params
	r.results = results
	r.handler = handler
	r.filters = filters
}

func (r *RPC) Namespace() string      { return r.namespace }
func (r *RPC) ServiceName() string    { return r.serviceName }
func (r *RPC) MethodName() string     { return r.methodName }
func (r *RPC) FullMethodName() string { return r.fullMethodName }
func (r *RPC) Params() Model          { return r.params }
func (r *RPC) Results() Model         { return r.results }
func (r *RPC) TraceID() string        { return r.traceID }

func (r *RPC) Do(ctx context.Context) error {
	if i := r.nextFilterIndex; i < len(r.filters) {
		r.nextFilterIndex++
		return r.filters[i](ctx, r)
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
