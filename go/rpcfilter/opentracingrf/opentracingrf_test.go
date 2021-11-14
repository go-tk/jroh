package opentracingrf_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/go-tk/jroh/go/apicommon/testdata/fooapi"
	"github.com/go-tk/jroh/go/middleware/opentracingmw"
	. "github.com/go-tk/jroh/go/rpcfilter/opentracingrf"
	"github.com/go-tk/testcase"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestOpenTracingMiddleware(t *testing.T) {
	type Input struct {
		TestServerFuncs  fooapi.TestServerFuncs
		TraceIDGenerator apicommon.TraceIDGenerator
		Ctx              context.Context
		Params           fooapi.DoSomething2Params
	}
	type Output struct {
		MSs []*mocktracer.MockSpan
	}
	type Workspace struct {
		testcase.WorkspaceBase

		MT *mocktracer.MockTracer
		TC fooapi.TestClient

		Input  Input
		Output Output
	}
	type ReqKey struct{}
	mocktracer.New()
	tc := testcase.New().
		AddTask(10, func(w *Workspace) {
			r := apicommon.NewRouter()
			mt := mocktracer.New()
			w.MT = mt
			so := apicommon.ServerOptions{
				Middlewares: map[apicommon.MethodIndex][]apicommon.ServerMiddleware{
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
			fooapi.RegisterTestServer(&w.Input.TestServerFuncs, r, so)
			co := apicommon.ClientOptions{
				RPCFilters: map[apicommon.MethodIndex][]apicommon.RPCHandler{
					apicommon.AnyMethod: {
						NewForClient(mt),
					},
				},
				Middlewares: map[apicommon.MethodIndex][]apicommon.ClientMiddleware{
					apicommon.AnyMethod: {
						opentracingmw.NewForClient(),
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
			w.Output.MSs = w.MT.FinishedSpans()
		})
	testcase.RunList(t,
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					results.MyOnOff = true
					return nil
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
			}).
			AddTask(21, func(w *Workspace) {
				mss := w.Output.MSs
				if !assert.Len(w.T(), mss, 1) {
					w.T().FailNow()
				}
				ms := mss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ms.OperationName)
				assert.Equal(w.T(), map[string]interface{}{
					"span.kind":         ext.SpanKindRPCClientEnum,
					"component":         "JROH",
					"jroh.is_requested": true,
					"http.status_code":  uint16(200),
					"jroh.error_code":   int32(0),
				}, ms.Tags())
				mlrs := ms.Logs()
				if assert.Len(w.T(), mlrs, 1) {
					mlr := mlrs[0]
					assert.Equal(w.T(), []mocktracer.MockKeyValue{
						{Key: "event", ValueKind: reflect.String, ValueString: "outgoing rpc"},
						{Key: "trace_id", ValueKind: reflect.String, ValueString: "tid1"},
						{Key: "url", ValueKind: reflect.String, ValueString: "http://127.0.0.1/rpc/Foo.Test.DoSomething2"},
					}, mlr.Fields)
				}
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					return errors.New("just wrong")
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
			}).
			AddTask(21, func(w *Workspace) {
				mss := w.Output.MSs
				if !assert.Len(w.T(), mss, 1) {
					w.T().FailNow()
				}
				ms := mss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ms.OperationName)
				assert.Equal(w.T(), map[string]interface{}{
					"span.kind":         ext.SpanKindRPCClientEnum,
					"component":         "JROH",
					"jroh.is_requested": true,
					"http.status_code":  uint16(200),
					"jroh.error_code":   int32(-32603),
				}, ms.Tags())
				mlrs := ms.Logs()
				if assert.Len(w.T(), mlrs, 1) {
					mlr := mlrs[0]
					assert.Equal(w.T(), []mocktracer.MockKeyValue{
						{Key: "event", ValueKind: reflect.String, ValueString: "outgoing rpc"},
						{Key: "trace_id", ValueKind: reflect.String, ValueString: "tid1"},
						{Key: "url", ValueKind: reflect.String, ValueString: "http://127.0.0.1/rpc/Foo.Test.DoSomething2"},
						{Key: "api_error", ValueKind: reflect.String, ValueString: "internal error"},
					}, mlr.Fields)
				}
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
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
				mss := w.Output.MSs
				if !assert.Len(w.T(), mss, 1) {
					w.T().FailNow()
				}
				ms := mss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ms.OperationName)
				assert.Equal(w.T(), map[string]interface{}{
					"span.kind":         ext.SpanKindRPCClientEnum,
					"component":         "JROH",
					"jroh.is_requested": false,
					"error":             true,
				}, ms.Tags())
				mlrs := ms.Logs()
				if assert.Len(w.T(), mlrs, 1) {
					mlr := mlrs[0]
					assert.Equal(w.T(), []mocktracer.MockKeyValue{
						{Key: "event", ValueKind: reflect.String, ValueString: "outgoing rpc"},
						{Key: "url", ValueKind: reflect.String, ValueString: "http://127.0.0.1/rpc/Foo.Test.DoSomething2"},
						{Key: "native_error", ValueKind: reflect.String, ValueString: "params encoding failed: json: error calling MarshalJSON for type *fooapi.MyStructString: bad word"},
					}, mlr.Fields)
				}
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				apicommon.DebugMode = true
				w.AddCleanup(func() {
					apicommon.DebugMode = false
				})
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					results.MyOnOff = true
					return nil
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
			}).
			AddTask(21, func(w *Workspace) {
				mss := w.Output.MSs
				if !assert.Len(w.T(), mss, 1) {
					w.T().FailNow()
				}
				ms := mss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ms.OperationName)
				assert.Equal(w.T(), map[string]interface{}{
					"span.kind":         ext.SpanKindRPCClientEnum,
					"component":         "JROH",
					"jroh.is_requested": true,
					"http.status_code":  uint16(200),
					"jroh.error_code":   int32(0),
				}, ms.Tags())
				mlrs := ms.Logs()
				if assert.Len(w.T(), mlrs, 1) {
					mlr := mlrs[0]
					assert.Equal(w.T(), []mocktracer.MockKeyValue{
						{Key: "event", ValueKind: reflect.String, ValueString: "outgoing rpc"},
						{Key: "trace_id", ValueKind: reflect.String, ValueString: "tid1"},
						{Key: "url", ValueKind: reflect.String, ValueString: "http://127.0.0.1/rpc/Foo.Test.DoSomething2"},
						{Key: "params", ValueKind: reflect.String, ValueString: "{\n  \"myOnOff\": false\n}\n"},
						{Key: "resp", ValueKind: reflect.String, ValueString: "{\n  \"results\": {\n    \"myOnOff\": true\n  }\n}\n"},
					}, mlr.Fields)
				}
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					r := ctx.Value(ReqKey{}).(*http.Request)
					sc, err := w.MT.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
					if assert.NoError(w.T(), err) {
						msc := sc.(mocktracer.MockSpanContext)
						assert.NotEqual(w.T(), -9999, msc.SpanID)
					}
					results.MyOnOff = true
					return nil
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
			}).
			AddTask(19, func(w *Workspace) {
				spanParent := w.MT.StartSpan("temp")
				spanParent.(*mocktracer.MockSpan).SpanContext.SpanID = -9999
				ctx := opentracing.ContextWithSpan(context.Background(), spanParent)
				w.Input.Ctx = ctx
			}).
			AddTask(21, func(w *Workspace) {
				mss := w.Output.MSs
				if !assert.Len(w.T(), mss, 1) {
					w.T().FailNow()
				}
				ms := mss[0]
				assert.Equal(w.T(), -9999, ms.ParentID)
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ms.OperationName)
				assert.Equal(w.T(), map[string]interface{}{
					"span.kind":         ext.SpanKindRPCClientEnum,
					"component":         "JROH",
					"jroh.is_requested": true,
					"http.status_code":  uint16(200),
					"jroh.error_code":   int32(0),
				}, ms.Tags())
				mlrs := ms.Logs()
				if assert.Len(w.T(), mlrs, 1) {
					mlr := mlrs[0]
					assert.Equal(w.T(), []mocktracer.MockKeyValue{
						{Key: "event", ValueKind: reflect.String, ValueString: "outgoing rpc"},
						{Key: "trace_id", ValueKind: reflect.String, ValueString: "tid1"},
						{Key: "url", ValueKind: reflect.String, ValueString: "http://127.0.0.1/rpc/Foo.Test.DoSomething2"},
					}, mlr.Fields)
				}
			}),
	)
}
