package tracerrf_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/go-tk/jroh/go/apicommon/testdata/fooapi"
	"github.com/go-tk/jroh/go/middleware/tracermw"
	. "github.com/go-tk/jroh/go/rpcfilter/tracerrf"
	"github.com/go-tk/testcase"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func TestForClient(t *testing.T) {
	type Input struct {
		TestServiceFuncs fooapi.TestServiceFuncs
		TraceIDGenerator apicommon.TraceIDGenerator
		Ctx              context.Context
		Params           fooapi.DoSomething2Params
	}
	type Output struct {
		SSs tracetest.SpanStubs
	}
	type Workspace struct {
		testcase.WorkspaceBase

		IME *tracetest.InMemoryExporter
		Tr  trace.Tracer
		TC  fooapi.TestClient
		SP  trace.Span

		Input  Input
		Output Output
	}
	type ReqKey struct{}
	mocktracer.New()
	tc := testcase.New().
		AddTask(10, func(w *Workspace) {
			r := apicommon.NewRouter()
			w.IME = tracetest.NewInMemoryExporter()
			tp := tracesdk.NewTracerProvider(tracesdk.WithSyncer(w.IME))
			w.Tr = tp.Tracer("testing")
			so := apicommon.ServerOptions{
				Middlewares: apicommon.ServerMiddlewares{
					apicommon.AnyMethod: {
						func(handler http.Handler) http.Handler {
							return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
								ctx := r.Context()
								ctx = context.WithValue(ctx, ReqKey{}, r)
								r = r.WithContext(ctx)
								handler.ServeHTTP(w, r)
							})
						},
					},
				},
				TraceIDGenerator: w.Input.TraceIDGenerator,
			}
			fooapi.RegisterTestService(&w.Input.TestServiceFuncs, r, so)
			co := apicommon.ClientOptions{
				RPCFilters: apicommon.RPCFilters{
					apicommon.AnyMethod: {
						NewForClient(w.Tr),
					},
				},
				Middlewares: apicommon.ClientMiddlewares{
					apicommon.AnyMethod: {
						tracermw.NewForClient(),
					},
				},
				Transport: apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
					responseRecorder := httptest.NewRecorder()
					r.ServeHTTP(responseRecorder, request.WithContext(context.Background()))
					response := responseRecorder.Result()
					return response, nil
				}),
			}
			w.TC = fooapi.NewTestClient("http://127.0.0.1", co)
		}).
		AddTask(20, func(w *Workspace) {
			ctx := w.Input.Ctx
			if ctx == nil {
				ctx = context.Background()
			}
			w.TC.DoSomething2(ctx, &w.Input.Params)
			w.Output.SSs = w.IME.GetSpans()
		})
	testcase.RunList(t,
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServiceFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					results.MyOnOff = true
					return nil
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
			}).
			AddTask(21, func(w *Workspace) {
				sss := w.Output.SSs
				if !assert.Len(w.T(), sss, 1) {
					w.T().FailNow()
				}
				ss := sss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ss.Name)
				assert.Equal(w.T(), trace.SpanKind(trace.SpanKindClient), ss.SpanKind)
				assert.Equal(w.T(), []attribute.KeyValue{
					semconv.RPCSystemKey.String("JROH"),
					semconv.RPCServiceKey.String("Foo.Test"),
					semconv.RPCMethodKey.String("DoSomething2"),
					attribute.Bool("rpc.jroh.is_requested", true),
					semconv.HTTPStatusCodeKey.Int(200),
					attribute.Int("rpc.jroh.error_code", int(0)),
				}, ss.Attributes)
				for i := range ss.Events {
					ss.Events[i].Time = time.Time{}
				}
				assert.Equal(w.T(), []tracesdk.Event{
					{
						Name: "outgoing rpc",
						Attributes: []attribute.KeyValue{
							attribute.String("trace_id", "tid1"),
							attribute.String("url", "http://127.0.0.1/rpc/Foo.Test.DoSomething2"),
						},
					},
				}, ss.Events)
				assert.Equal(w.T(), codes.Ok, ss.Status.Code)
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServiceFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					return errors.New("just wrong")
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
			}).
			AddTask(21, func(w *Workspace) {
				sss := w.Output.SSs
				if !assert.Len(w.T(), sss, 1) {
					w.T().FailNow()
				}
				ss := sss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ss.Name)
				assert.Equal(w.T(), trace.SpanKind(trace.SpanKindClient), ss.SpanKind)
				assert.Equal(w.T(), []attribute.KeyValue{
					semconv.RPCSystemKey.String("JROH"),
					semconv.RPCServiceKey.String("Foo.Test"),
					semconv.RPCMethodKey.String("DoSomething2"),
					attribute.Bool("rpc.jroh.is_requested", true),
					semconv.HTTPStatusCodeKey.Int(200),
					attribute.Int("rpc.jroh.error_code", int(-32603)),
				}, ss.Attributes)
				for i := range ss.Events {
					ss.Events[i].Time = time.Time{}
				}
				assert.Equal(w.T(), []tracesdk.Event{
					{
						Name: "outgoing rpc",
						Attributes: []attribute.KeyValue{
							attribute.String("trace_id", "tid1"),
							attribute.String("url", "http://127.0.0.1/rpc/Foo.Test.DoSomething2"),
							attribute.String("api_error", "internal error"),
						},
					},
				}, ss.Events)
				assert.Equal(w.T(), codes.Ok, ss.Status.Code)
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServiceFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					results.MyOnOff = true
					return nil
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
				tmp := fooapi.MyStructString{}
				tmp.TheStringA = "taboo"
				w.Input.Params.MyStructString = &tmp
			}).
			AddTask(21, func(w *Workspace) {
				sss := w.Output.SSs
				if !assert.Len(w.T(), sss, 1) {
					w.T().FailNow()
				}
				ss := sss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ss.Name)
				assert.Equal(w.T(), trace.SpanKind(trace.SpanKindClient), ss.SpanKind)
				assert.Equal(w.T(), []attribute.KeyValue{
					semconv.RPCSystemKey.String("JROH"),
					semconv.RPCServiceKey.String("Foo.Test"),
					semconv.RPCMethodKey.String("DoSomething2"),
					attribute.Bool("rpc.jroh.is_requested", false),
				}, ss.Attributes)
				for i := range ss.Events {
					ss.Events[i].Time = time.Time{}
				}
				assert.Equal(w.T(), []tracesdk.Event{
					{
						Name: "outgoing rpc",
						Attributes: []attribute.KeyValue{
							attribute.String("url", "http://127.0.0.1/rpc/Foo.Test.DoSomething2"),
							attribute.String("native_error", "params encoding failed: json: error calling MarshalJSON for type *fooapi.MyStructString: bad word"),
						},
					},
				}, ss.Events)
				assert.Equal(w.T(), codes.Error, ss.Status.Code)
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				apicommon.DebugMode = true
				w.AddCleanup(func() {
					apicommon.DebugMode = false
				})
				w.Input.TestServiceFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					results.MyOnOff = true
					return nil
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
			}).
			AddTask(21, func(w *Workspace) {
				sss := w.Output.SSs
				if !assert.Len(w.T(), sss, 1) {
					w.T().FailNow()
				}
				ss := sss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ss.Name)
				assert.Equal(w.T(), trace.SpanKind(trace.SpanKindClient), ss.SpanKind)
				assert.Equal(w.T(), []attribute.KeyValue{
					semconv.RPCSystemKey.String("JROH"),
					semconv.RPCServiceKey.String("Foo.Test"),
					semconv.RPCMethodKey.String("DoSomething2"),
					attribute.Bool("rpc.jroh.is_requested", true),
					semconv.HTTPStatusCodeKey.Int(200),
					attribute.Int("rpc.jroh.error_code", int(0)),
				}, ss.Attributes)
				for i := range ss.Events {
					ss.Events[i].Time = time.Time{}
				}
				assert.Equal(w.T(), []tracesdk.Event{
					{
						Name: "outgoing rpc",
						Attributes: []attribute.KeyValue{
							attribute.String("trace_id", "tid1"),
							attribute.String("url", "http://127.0.0.1/rpc/Foo.Test.DoSomething2"),
							attribute.String("params", "{\n  \"myOnOff\": false\n}\n"),
							attribute.String("resp", "{\n  \"results\": {\n    \"myOnOff\": true\n  }\n}\n"),
						},
					},
				}, ss.Events)
				assert.Equal(w.T(), codes.Ok, ss.Status.Code)
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServiceFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					r := ctx.Value(ReqKey{}).(*http.Request)
					ctx = propagation.TraceContext{}.Extract(ctx, propagation.HeaderCarrier(r.Header))
					sc := trace.SpanContextFromContext(ctx)
					assert.Equal(w.T(), w.SP.SpanContext().TraceID(), sc.TraceID())
					results.MyOnOff = true
					return nil
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
			}).
			AddTask(19, func(w *Workspace) {
				ctx, spanParent := w.Tr.Start(context.Background(), "temp")
				w.SP = spanParent
				w.Input.Ctx = ctx
			}).
			AddTask(21, func(w *Workspace) {
				sss := w.Output.SSs
				if !assert.Len(w.T(), sss, 1) {
					w.T().FailNow()
				}
				ss := sss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ss.Name)
				assert.Equal(w.T(), w.SP.SpanContext().SpanID(), ss.Parent.SpanID())
				assert.Equal(w.T(), trace.SpanKind(trace.SpanKindClient), ss.SpanKind)
				assert.Equal(w.T(), []attribute.KeyValue{
					semconv.RPCSystemKey.String("JROH"),
					semconv.RPCServiceKey.String("Foo.Test"),
					semconv.RPCMethodKey.String("DoSomething2"),
					attribute.Bool("rpc.jroh.is_requested", true),
					semconv.HTTPStatusCodeKey.Int(200),
					attribute.Int("rpc.jroh.error_code", int(0)),
				}, ss.Attributes)
				for i := range ss.Events {
					ss.Events[i].Time = time.Time{}
				}
				assert.Equal(w.T(), []tracesdk.Event{
					{
						Name: "outgoing rpc",
						Attributes: []attribute.KeyValue{
							attribute.String("trace_id", "tid1"),
							attribute.String("url", "http://127.0.0.1/rpc/Foo.Test.DoSomething2"),
						},
					},
				}, ss.Events)
				assert.Equal(w.T(), codes.Ok, ss.Status.Code)
			}),
	)
}
