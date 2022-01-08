package opentelemetrymw

func rpcServiceKey(fullMethodName string, methodName string) string {
	return fullMethodName[:len(fullMethodName)-len(methodName)-1]
}
