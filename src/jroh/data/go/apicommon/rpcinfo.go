package apicommon

import "context"

type RPCInfo struct {
	namespace   string
	serviceName string
	methodName  string
	id          string

	rawParams   []byte
	params      interface{}
	error       *Error
	internalErr error
	stackTrace  string
	results     interface{}

	respWriteErr error
}

func NewRPCInfo(namespace string, serviceName string, methodName string, id string) *RPCInfo {
	return &RPCInfo{
		namespace:   namespace,
		serviceName: serviceName,
		methodName:  methodName,
		id:          id,
	}
}

func (ri *RPCInfo) Namespace() string   { return ri.namespace }
func (ri *RPCInfo) ServiceName() string { return ri.serviceName }
func (ri *RPCInfo) MethodName() string  { return ri.methodName }
func (ri *RPCInfo) ID() string          { return ri.id }

func (ri *RPCInfo) RawParams() []byte             { return ri.rawParams }
func (ri *RPCInfo) SetRawParams(rawParams []byte) { ri.rawParams = rawParams }

func (ri *RPCInfo) Params() interface{}          { return ri.params }
func (ri *RPCInfo) SetParams(params interface{}) { ri.params = params }

func (ri *RPCInfo) Error() *Error         { return ri.error }
func (ri *RPCInfo) SetError(error *Error) { ri.error = error }

func (ri *RPCInfo) InternalErr() error               { return ri.internalErr }
func (ri *RPCInfo) SetInternalErr(internalErr error) { ri.internalErr = internalErr }

func (ri *RPCInfo) StackTrace() string              { return ri.stackTrace }
func (ri *RPCInfo) SetStackTrace(stackTrace string) { ri.stackTrace = stackTrace }

func (ri *RPCInfo) Results() interface{}           { return ri.results }
func (ri *RPCInfo) SetResults(results interface{}) { ri.results = results }

func (ri *RPCInfo) RespWriteErr() error { return ri.respWriteErr }
func (ri *RPCInfo) SetRespWriteErr(respWriteErr error) {
	ri.respWriteErr = respWriteErr
}

type RPCInfoFactory func(id string) *RPCInfo

type contextValueRPCInfo struct{}

func ContextWithRPCInfo(ctx context.Context, rpcInfo *RPCInfo) context.Context {
	return context.WithValue(ctx, contextValueRPCInfo{}, rpcInfo)
}

func RPCInfoFromContext(ctx context.Context) *RPCInfo {
	return ctx.Value(contextValueRPCInfo{}).(*RPCInfo)
}
