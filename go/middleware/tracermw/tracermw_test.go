package tracermw_test

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
	. "github.com/go-tk/jroh/go/middleware/tracermw"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func TestForServer(t *testing.T) {
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
	tc := testcase.New().
		AddTask(10, func(w *Workspace) {
			r := apicommon.NewRouter()
			w.IME = tracetest.NewInMemoryExporter()
			tp := tracesdk.NewTracerProvider(tracesdk.WithSyncer(w.IME))
			w.Tr = tp.Tracer("testing")
			so := apicommon.ServerOptions{
				TraceIDGenerator: w.Input.TraceIDGenerator,
				Middlewares: apicommon.ServerMiddlewares{
					apicommon.AnyMethod: {
						NewForServer(w.Tr),
					},
				},
			}
			fooapi.RegisterTestService(&w.Input.TestServiceFuncs, r, so)
			co := apicommon.ClientOptions{
				Middlewares: apicommon.ClientMiddlewares{
					apicommon.AnyMethod: {
						NewForClient(),
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
				assert.Equal(w.T(), trace.SpanKind(trace.SpanKindServer), ss.SpanKind)
				assert.Equal(w.T(), []attribute.KeyValue{
					semconv.RPCSystemKey.String("JROH"),
					semconv.RPCServiceKey.String("Foo.Test"),
					semconv.RPCMethodKey.String("DoSomething2"),
					semconv.HTTPStatusCodeKey.Int(200),
					attribute.Int("rpc.jroh.error_code", int(0)),
				}, ss.Attributes)
				for i := range ss.Events {
					ss.Events[i].Time = time.Time{}
				}
				assert.Equal(w.T(), []tracesdk.Event{
					{
						Name: "incoming rpc",
						Attributes: []attribute.KeyValue{
							attribute.String("trace_id", "tid1"),
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
				assert.Equal(w.T(), trace.SpanKind(trace.SpanKindServer), ss.SpanKind)
				assert.Equal(w.T(), []attribute.KeyValue{
					semconv.RPCSystemKey.String("JROH"),
					semconv.RPCServiceKey.String("Foo.Test"),
					semconv.RPCMethodKey.String("DoSomething2"),
					semconv.HTTPStatusCodeKey.Int(200),
					attribute.Int("rpc.jroh.error_code", -32603),
				}, ss.Attributes)
				for i := range ss.Events {
					ss.Events[i].Time = time.Time{}
				}
				assert.Equal(w.T(), []tracesdk.Event{
					{
						Name: "incoming rpc",
						Attributes: []attribute.KeyValue{
							attribute.String("trace_id", "tid1"),
							attribute.String("api_error", "internal error"),
							attribute.String("native_error", "just wrong"),
						},
					},
				}, ss.Events)
				assert.Equal(w.T(), codes.Error, ss.Status.Code)
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServiceFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					tmp := fooapi.MyStructString{}
					tmp.TheStringA = "taboo"
					results.MyStructString = &tmp
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
				assert.Equal(w.T(), trace.SpanKind(trace.SpanKindServer), ss.SpanKind)
				assert.Equal(w.T(), []attribute.KeyValue{
					semconv.RPCSystemKey.String("JROH"),
					semconv.RPCServiceKey.String("Foo.Test"),
					semconv.RPCMethodKey.String("DoSomething2"),
					semconv.HTTPStatusCodeKey.Int(500),
					attribute.Int("rpc.jroh.error_code", 0),
				}, ss.Attributes)
				for i := range ss.Events {
					ss.Events[i].Time = time.Time{}
				}
				assert.Equal(w.T(), []tracesdk.Event{
					{
						Name: "incoming rpc",
						Attributes: []attribute.KeyValue{
							attribute.String("trace_id", "tid1"),
							attribute.String("native_error", "resp encoding failed: json: error calling MarshalJSON for type *fooapi.MyStructString: bad word"),
						},
					},
				}, ss.Events)
				assert.Equal(w.T(), codes.Error, ss.Status.Code)
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
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
				assert.Equal(w.T(), trace.SpanKind(trace.SpanKindServer), ss.SpanKind)
				assert.Equal(w.T(), []attribute.KeyValue{
					semconv.RPCSystemKey.String("JROH"),
					semconv.RPCServiceKey.String("Foo.Test"),
					semconv.RPCMethodKey.String("DoSomething2"),
					semconv.HTTPStatusCodeKey.Int(200),
					attribute.Int("rpc.jroh.error_code", -32000),
				}, ss.Attributes)
				for i := range ss.Events {
					ss.Events[i].Time = time.Time{}
				}
				assert.Equal(w.T(), []tracesdk.Event{
					{
						Name: "incoming rpc",
						Attributes: []attribute.KeyValue{
							attribute.String("trace_id", "tid1"),
							attribute.String("api_error", "not implemented"),
						},
					},
				}, ss.Events)
				assert.Equal(w.T(), codes.Ok, ss.Status.Code)
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
				assert.Equal(w.T(), trace.SpanKind(trace.SpanKindServer), ss.SpanKind)
				assert.Equal(w.T(), []attribute.KeyValue{
					semconv.RPCSystemKey.String("JROH"),
					semconv.RPCServiceKey.String("Foo.Test"),
					semconv.RPCMethodKey.String("DoSomething2"),
					semconv.HTTPStatusCodeKey.Int(200),
					attribute.Int("rpc.jroh.error_code", 0),
				}, ss.Attributes)
				for i := range ss.Events {
					ss.Events[i].Time = time.Time{}
				}
				assert.Equal(w.T(), []tracesdk.Event{
					{
						Name: "incoming rpc",
						Attributes: []attribute.KeyValue{
							attribute.String("trace_id", "tid1"),
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
				assert.Equal(w.T(), trace.SpanKind(trace.SpanKindServer), ss.SpanKind)
				assert.Equal(w.T(), []attribute.KeyValue{
					semconv.RPCSystemKey.String("JROH"),
					semconv.RPCServiceKey.String("Foo.Test"),
					semconv.RPCMethodKey.String("DoSomething2"),
					semconv.HTTPStatusCodeKey.Int(200),
					attribute.Int("rpc.jroh.error_code", int(0)),
				}, ss.Attributes)
				for i := range ss.Events {
					ss.Events[i].Time = time.Time{}
				}
				assert.Equal(w.T(), []tracesdk.Event{
					{
						Name: "incoming rpc",
						Attributes: []attribute.KeyValue{
							attribute.String("trace_id", "tid1"),
						},
					},
				}, ss.Events)
				assert.Equal(w.T(), codes.Ok, ss.Status.Code)
			}),
	)
}
