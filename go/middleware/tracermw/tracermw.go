package tracermw

import (
	"net/http"
	"unsafe"

	"github.com/go-tk/jroh/go/apicommon"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func NewForServer(tracer trace.Tracer) apicommon.ServerMiddleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			incomingRPC := apicommon.MustGetRPCFromContext(ctx).IncomingRPC()
			ctx = propagation.TraceContext{}.Extract(ctx, propagation.HeaderCarrier(r.Header))
			ctx, span := tracer.Start(ctx, incomingRPC.FullMethodName(), trace.WithSpanKind(trace.SpanKindServer))
			r = r.WithContext(ctx)
			// Before
			handler.ServeHTTP(w, r)
			// After
			span.SetAttributes(
				semconv.RPCSystemKey.String("JROH"),
				semconv.RPCServiceKey.String(rpcServiceKey(&incomingRPC.RPC)),
				semconv.RPCMethodKey.String(incomingRPC.MethodName()),
				semconv.HTTPStatusCodeKey.Int(incomingRPC.StatusCode()),
				attribute.Int("rpc.jroh.error_code", int(incomingRPC.Error().Code)),
			)
			if incomingRPC.StatusCode()/100 == 5 || incomingRPC.Error().Code == apicommon.ErrorInternal {
				span.SetStatus(codes.Error, "")
			} else {
				span.SetStatus(codes.Ok, "")
			}
			kvs := []attribute.KeyValue{
				attribute.String("trace_id", incomingRPC.TraceID()),
			}
			if apicommon.DebugMode {
				if rawParams := incomingRPC.RawParams(); rawParams != nil {
					kvs = append(kvs, attribute.String("params", bytesToString(rawParams)))
				}
			}
			if incomingRPC.Error().Code != 0 {
				kvs = append(kvs, attribute.String("api_error", incomingRPC.Error().Message))
			}
			if err := incomingRPC.Err(); err != nil {
				kvs = append(kvs, attribute.String("native_error", err.Error()))
			}
			if apicommon.DebugMode {
				if rawResp := incomingRPC.RawResp(); rawResp != nil {
					kvs = append(kvs, attribute.String("resp", bytesToString(rawResp)))
				}
			}
			span.AddEvent("incoming rpc", trace.WithAttributes(kvs...))
			span.End()
		})
	}
}

func rpcServiceKey(rpc *apicommon.RPC) string {
	fullMethodName := rpc.FullMethodName()
	return fullMethodName[:len(fullMethodName)-len(rpc.MethodName())-1]
}

func bytesToString(bytes []byte) string { return *(*string)(unsafe.Pointer(&bytes)) }

func NewForClient() apicommon.ClientMiddleware {
	return func(transport http.RoundTripper) http.RoundTripper {
		return apicommon.TransportFunc(func(request *http.Request) (returnedResponse *http.Response, returnedErr error) {
			propagation.TraceContext{}.Inject(request.Context(), propagation.HeaderCarrier(request.Header))
			returnedResponse, returnedErr = transport.RoundTrip(request)
			return
		})
	}
}
