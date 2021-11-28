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
		ctx = opentracing.ContextWithSpan(ctx, span)
		// Before
		returnedErr = outgoingRPC.Do(ctx)
		// After
		var err error
		switch returnedErr.(type) {
		case *apicommon.UnexpectedStatusCodeError, *apicommon.Error:
		default:
			err = returnedErr
		}
		ext.Component.Set(span, "JROH")
		span.SetTag("jroh.is_requested", outgoingRPC.IsRequested())
		if outgoingRPC.IsRequested() {
			ext.HTTPStatusCode.Set(span, uint16(outgoingRPC.StatusCode()))
			span.SetTag("jroh.error_code", int32(outgoingRPC.Error().Code))
		}
		if err != nil {
			ext.Error.Set(span, true)
		}
		logFields := []log.Field{log.Event("outgoing rpc")}
		if traceID := outgoingRPC.TraceID(); traceID != "" {
			logFields = append(logFields, log.String("trace_id", traceID))
		}
		logFields = append(logFields, log.String("url", outgoingRPC.URL()))
		if apicommon.DebugMode {
			if rawParams := outgoingRPC.RawParams(); rawParams != nil {
				logFields = append(logFields, log.String("params", bytesToString(rawParams)))
			}
		}
		if outgoingRPC.IsRequested() {
			if apicommon.DebugMode {
				if rawResp := outgoingRPC.RawResp(); rawResp != nil {
					logFields = append(logFields, log.String("resp", bytesToString(rawResp)))
				}
			}
			if outgoingRPC.Error().Code != 0 {
				logFields = append(logFields, log.String("api_error", outgoingRPC.Error().Message))
			}
		}
		if err != nil {
			logFields = append(logFields, log.String("native_error", err.Error()))
		}
		span.LogFields(logFields...)
		span.Finish()
		return
	}
}

func bytesToString(bytes []byte) string { return *(*string)(unsafe.Pointer(&bytes)) }
