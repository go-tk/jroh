package opentracingrf

import (
	"context"
	"unsafe"

	"github.com/go-tk/jroh/go/apicommon"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func NewForClient(tracer opentracing.Tracer) apicommon.RPCHandler {
	return func(ctx context.Context, rpc *apicommon.RPC) (returnedErr error) {
		outgoingRPC := rpc.OutgoingRPC()
		var spanParentContext opentracing.SpanContext
		if spanParent := opentracing.SpanFromContext(ctx); spanParent != nil {
			tracer = spanParent.Tracer()
			spanParentContext = spanParent.Context()
		}
		span := tracer.StartSpan(
			outgoingRPC.FullMethodName(),
			opentracing.ChildOf(spanParentContext),
			ext.SpanKindRPCClient,
		)
		ext.Component.Set(span, "JROH")
		ctx = opentracing.ContextWithSpan(ctx, span)
		// Before
		returnedErr = outgoingRPC.Do(ctx)
		// After
		logFields := []log.Field{log.Event("outgoing rpc")}
		if traceID := outgoingRPC.TraceID(); traceID != "" {
			logFields = append(logFields, log.String("trace_id", traceID))
		}
		if apicommon.DebugMode {
			if rawParams := outgoingRPC.RawParams(); rawParams != nil {
				logFields = append(logFields, log.String("params", bytesToString(rawParams)))
			}
		}
		span.SetTag("jroh.is_requested", outgoingRPC.IsRequested())
		if outgoingRPC.IsRequested() {
			ext.HTTPStatusCode.Set(span, uint16(outgoingRPC.StatusCode()))
			span.SetTag("jroh.error_code", int32(outgoingRPC.Error().Code))
			if outgoingRPC.Error().Code != 0 {
				logFields = append(logFields, log.String("api_error", outgoingRPC.Error().Message))
			}
			if apicommon.DebugMode {
				if rawResp := outgoingRPC.RawResp(); rawResp != nil {
					logFields = append(logFields, log.String("resp", bytesToString(rawResp)))
				}
			}
		} else {
			if returnedErr != nil {
				ext.Error.Set(span, true)
				logFields = append(logFields, log.String("pre_request_error", returnedErr.Error()))
			}
		}
		span.FinishWithOptions(opentracing.FinishOptions{
			LogRecords: []opentracing.LogRecord{
				{Fields: logFields},
			},
		})
		return
	}
}

func bytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
