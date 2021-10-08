package opentracingmw

import (
	"net/http"

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
			const errorCodeKey = "jroh.error_code"
			if errorCode := incomingRPC.Error().Code; errorCode == 0 {
				span.SetTag(errorCodeKey, "-")
			} else {
				span.SetTag(errorCodeKey, errorCode)
			}
			if incomingRPC.Error().Code != 0 {
				span.LogFields(log.Event("rpc error"), log.Message(incomingRPC.Error().Message))
			}
			if internalErr := incomingRPC.InternalErr(); internalErr != nil {
				span.LogFields(log.Event("internal error"), log.Message(internalErr.Error()))
			}
			span.Finish()
		})
	}
}

func NewForClient(tracer opentracing.Tracer) apicommon.ClientMiddleware {
	return func(transport http.RoundTripper) http.RoundTripper {
		return apicommon.TransportFunc(func(request *http.Request) (returnedResponse *http.Response, returnedErr error) {
			span := opentracing.SpanFromContext(request.Context())
			if span != nil {
				if err := tracer.Inject(span.Context(), opentracing.HTTPHeaders,
					opentracing.HTTPHeadersCarrier(request.Header)); err != nil {
					panic(err)
				}
			}
			returnedResponse, returnedErr = transport.RoundTrip(request)
			return
		})
	}
}
