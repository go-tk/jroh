package opentracingmw

import (
	"net/http"
	"unsafe"

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
			ctx = opentracing.ContextWithSpan(ctx, span)
			r = r.WithContext(ctx)
			// Before
			handler.ServeHTTP(w, r)
			// After
			ext.Component.Set(span, "JROH")
			ext.HTTPStatusCode.Set(span, uint16(incomingRPC.StatusCode()))
			span.SetTag("jroh.error_code", int32(incomingRPC.Error().Code))
			if incomingRPC.StatusCode()/100 == 5 || incomingRPC.Error().Code == apicommon.ErrorInternal {
				ext.Error.Set(span, true)
			}
			logFields := []log.Field{
				log.Event("incoming rpc"),
				log.String("trace_id", incomingRPC.TraceID()),
			}
			if apicommon.DebugMode {
				if rawParams := incomingRPC.RawParams(); rawParams != nil {
					logFields = append(logFields, log.String("params", bytesToString(rawParams)))
				}
			}
			if incomingRPC.Error().Code != 0 {
				logFields = append(logFields, log.String("api_error", incomingRPC.Error().Message))
			}
			if err := incomingRPC.Err(); err != nil {
				logFields = append(logFields, log.String("native_error", err.Error()))
			}
			if apicommon.DebugMode {
				if rawResp := incomingRPC.RawResp(); rawResp != nil {
					logFields = append(logFields, log.String("resp", bytesToString(rawResp)))
				}
			}
			span.LogFields(logFields...)
			span.Finish()
		})
	}
}

func bytesToString(bytes []byte) string { return *(*string)(unsafe.Pointer(&bytes)) }

func NewForClient() apicommon.ClientMiddleware {
	return func(transport http.RoundTripper) http.RoundTripper {
		return apicommon.TransportFunc(func(request *http.Request) (returnedResponse *http.Response, returnedErr error) {
			span := opentracing.SpanFromContext(request.Context())
			if span != nil {
				if err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders,
					opentracing.HTTPHeadersCarrier(request.Header)); err != nil {
					panic(err)
				}
			}
			returnedResponse, returnedErr = transport.RoundTrip(request)
			return
		})
	}
}
