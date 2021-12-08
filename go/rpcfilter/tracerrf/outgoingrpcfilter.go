package tracerrf

import (
	"context"

	"github.com/go-tk/jroh/go/apicommon"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func NewOutgoingRPCFilter(tracer trace.Tracer) apicommon.OutgoingRPCHandler {
	return func(ctx context.Context, outgoingRPC *apicommon.OutgoingRPC) (returnedErr error) {
		ctx, span := tracer.Start(ctx, outgoingRPC.FullMethodName, trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()
		propagation.TraceContext{}.Inject(ctx, propagation.HeaderCarrier(outgoingRPC.OutboundHeader))

		returnedErr = outgoingRPC.Do(ctx)

		if returnedErr == nil {
			outgoingRPC.ReadRawResults()
		}
		span.SetAttributes(
			semconv.RPCSystemKey.String("JROH"),
			semconv.RPCServiceKey.String(rpcServiceKey(outgoingRPC.FullMethodName, outgoingRPC.MethodName)),
			semconv.RPCMethodKey.String(outgoingRPC.MethodName),
			semconv.HTTPStatusCodeKey.Int(outgoingRPC.StatusCode),
			attribute.Int("rpc.jroh.error_code", int(outgoingRPC.ErrorCode)),
		)
		if returnedErr != nil && !apicommon.ErrIsTemporary(returnedErr) {
			span.SetStatus(codes.Error, "")
		} else {
			span.SetStatus(codes.Ok, "")
		}
		kvs := []attribute.KeyValue{
			attribute.String("trace_id", outgoingRPC.TraceID),
			attribute.String("url", outgoingRPC.URL),
		}
		if returnedErr != nil {
			kvs = append(kvs, attribute.String("err", returnedErr.Error()))
		}
		if apicommon.DebugMode {
			kvs = append(kvs,
				attribute.String("params", apicommon.BytesToString(outgoingRPC.RawParams)),
				attribute.String("results", apicommon.BytesToString(outgoingRPC.RawResults)),
			)
		}
		span.AddEvent("outgoing rpc", trace.WithAttributes(kvs...))
		return
	}
}
