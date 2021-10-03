package opentracingmw

import (
	"net/http"
	"strconv"

	"github.com/go-tk/jroh/go/apicommon"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func NewForServer(tracer opentracing.Tracer) apicommon.ServerMiddleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			incomingRPC := apicommon.MustGetRPCFromContext(ctx).IncomingRPC()
			spanContext, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
			if err != nil && err != opentracing.ErrSpanContextNotFound {
				panic(err)
			}
			span := tracer.StartSpan(
				incomingRPC.FullMethodName(),
				ext.RPCServerOption(spanContext),
			)
			ext.Component.Set(span, "JROH")
			ctx = opentracing.ContextWithSpan(ctx, span)
			r = r.WithContext(ctx)
			// Before
			handler.ServeHTTP(w, r)
			// After
			if incomingRPC.StatusCode()/100 == 5 || incomingRPC.InternalErr() != nil {
				ext.Error.Set(span, true)
			}
			ext.HTTPStatusCode.Set(span, uint16(incomingRPC.StatusCode()))
			var errorCodeStr string
			if errorCode := incomingRPC.Error().Code; errorCode == 0 {
				errorCodeStr = "-"
			} else {
				errorCodeStr = strconv.FormatInt(int64(errorCode), 10)
			}
			span.SetTag("jroh.error_code", errorCodeStr)
			if incomingRPC.Error().Code != 0 {
				span.LogFields(log.Event("error"), log.Message(incomingRPC.Error().Message))
			}
			if internalErr := incomingRPC.InternalErr(); internalErr != nil {
				span.LogFields(log.Event("internal error"), log.Message(internalErr.Error()))
			}
			span.Finish()
		})
	}
}
