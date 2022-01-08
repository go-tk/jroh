package opentelemetrymw

import (
	"context"

	"github.com/go-tk/jroh/go/apicommon"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func NewIncomingRPCFilter(tracer trace.Tracer) apicommon.IncomingRPCHandler {
	return func(ctx context.Context, incomingRPC *apicommon.IncomingRPC) (returnedErr error) {
		ctx = propagation.TraceContext{}.Extract(ctx, propagation.HeaderCarrier(incomingRPC.InboundHeader))
		ctx, span := tracer.Start(ctx, incomingRPC.FullMethodName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		returnedErr = incomingRPC.Do(ctx)

		if returnedErr == nil {
			incomingRPC.EncodeResults()
		}
		span.SetAttributes(
			semconv.RPCSystemKey.String("JROH"),
			semconv.RPCServiceKey.String(rpcServiceKey(incomingRPC.FullMethodName, incomingRPC.MethodName)),
			semconv.RPCMethodKey.String(incomingRPC.MethodName),
			semconv.HTTPStatusCodeKey.Int(incomingRPC.StatusCode),
			attribute.Int("rpc.jroh.error_code", int(incomingRPC.ErrorCode)),
		)
		if apicommon.ServerShouldReportError(returnedErr, incomingRPC.StatusCode) {
			span.SetStatus(codes.Error, "")
		} else {
			span.SetStatus(codes.Ok, "")
		}
		kvs := []attribute.KeyValue{
			attribute.String("trace_id", incomingRPC.TraceID),
			attribute.String("remote_ip", incomingRPC.RemoteIP),
		}
		if returnedErr != nil {
			kvs = append(kvs, attribute.String("err", returnedErr.Error()))
		}
		if apicommon.DebugMode {
			kvs = append(kvs,
				attribute.String("params", apicommon.BytesToString(incomingRPC.RawParams)),
				attribute.String("results", apicommon.BytesToString(incomingRPC.RawResults)),
			)
		}
		span.AddEvent("incoming rpc", trace.WithAttributes(kvs...))
		return
	}
}
