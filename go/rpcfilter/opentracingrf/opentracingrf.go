package opentracingrf

import (
	"context"

	"github.com/go-tk/jroh/go/apicommon"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func NewForClient(tracer opentracing.Tracer) apicommon.RPCHandler {
	return func(ctx context.Context, rpc *apicommon.RPC) (returnedErr error) {
		var spanParentContext opentracing.SpanContext
		if spanParent := opentracing.SpanFromContext(ctx); spanParent != nil {
			spanParentContext = spanParent.Context()
		}
		span := tracer.StartSpan(
			rpc.FullMethodName(),
			opentracing.ChildOf(spanParentContext),
			ext.SpanKindRPCClient,
		)
		ext.Component.Set(span, "JROH")
		// Before
		returnedErr = rpc.Do(ctx)
		// After
		outgoingRPC := rpc.OutgoingRPC()
		var preRequestErr error
		if returnedErr != nil && !outgoingRPC.IsRequested() {
			preRequestErr = returnedErr
		}
		if preRequestErr == nil {
			if statusCode := outgoingRPC.StatusCode(); statusCode == 0 {
				span.SetTag(string(ext.HTTPStatusCode), "-")
			} else {
				ext.HTTPStatusCode.Set(span, uint16(statusCode))
			}
			const errorCodeKey = "jroh.error_code"
			if errorCode := outgoingRPC.Error().Code; errorCode == 0 {
				span.SetTag(errorCodeKey, "-")
			} else {
				span.SetTag(errorCodeKey, errorCode)
			}
			if outgoingRPC.Error().Code != 0 {
				span.LogFields(log.Event("rpc error"), log.Message(outgoingRPC.Error().Message))
			}
		} else {
			ext.Error.Set(span, true)
			span.LogFields(log.Event("pre-request error"), log.Message(preRequestErr.Error()))
		}
		span.Finish()
		return
	}
}
