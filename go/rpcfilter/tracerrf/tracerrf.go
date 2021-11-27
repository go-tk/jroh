package tracerrf

import (
	"context"
	"unsafe"

	"github.com/go-tk/jroh/go/apicommon"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func NewForClient(tracer trace.Tracer) apicommon.RPCHandler {
	return func(ctx context.Context, rpc *apicommon.RPC) (returnedErr error) {
		outgoingRPC := rpc.OutgoingRPC()
		ctx, span := tracer.Start(ctx, outgoingRPC.FullMethodName(), trace.WithSpanKind(trace.SpanKindClient))
		// Before
		returnedErr = outgoingRPC.Do(ctx)
		// After
		span.SetAttributes(
			semconv.RPCSystemKey.String("JROH"),
			semconv.RPCServiceKey.String(rpcServiceKey(&outgoingRPC.RPC)),
			semconv.RPCMethodKey.String(outgoingRPC.MethodName()),
			attribute.Bool("rpc.jroh.is_requested", outgoingRPC.IsRequested()),
		)
		if outgoingRPC.IsRequested() {
			span.SetAttributes(
				semconv.HTTPStatusCodeKey.Int(outgoingRPC.StatusCode()),
				attribute.Int("rpc.jroh.error_code", int(outgoingRPC.Error().Code)),
			)
		}
		var err error
		if _, ok := returnedErr.(*apicommon.Error); !ok {
			err = returnedErr
		}
		if err != nil {
			span.SetStatus(codes.Error, "")
		} else {
			span.SetStatus(codes.Ok, "")
		}
		var kvs []attribute.KeyValue
		if traceID := outgoingRPC.TraceID(); traceID != "" {
			kvs = append(kvs, attribute.String("trace_id", traceID))
		}
		kvs = append(kvs, attribute.String("url", outgoingRPC.URL()))
		if apicommon.DebugMode {
			if rawParams := outgoingRPC.RawParams(); rawParams != nil {
				kvs = append(kvs, attribute.String("params", bytesToString(rawParams)))
			}
		}
		if outgoingRPC.IsRequested() {
			if apicommon.DebugMode {
				if rawResp := outgoingRPC.RawResp(); rawResp != nil {
					kvs = append(kvs, attribute.String("resp", bytesToString(rawResp)))
				}
			}
			if outgoingRPC.Error().Code != 0 {
				kvs = append(kvs, attribute.String("api_error", outgoingRPC.Error().Message))
			}
		}
		if err != nil {
			kvs = append(kvs, attribute.String("native_error", err.Error()))
		}
		span.AddEvent("outgoing rpc", trace.WithAttributes(kvs...))
		span.End()
		return
	}
}

func rpcServiceKey(rpc *apicommon.RPC) string {
	fullMethodName := rpc.FullMethodName()
	return fullMethodName[:len(fullMethodName)-len(rpc.MethodName())-1]
}

func bytesToString(bytes []byte) string { return *(*string)(unsafe.Pointer(&bytes)) }
