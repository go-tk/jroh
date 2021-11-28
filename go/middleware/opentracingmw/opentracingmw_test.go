package opentracingmw_test

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
	. "github.com/go-tk/jroh/go/middleware/opentracingmw"
	"github.com/go-tk/testcase"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestForServer(t *testing.T) {
	type Input struct {
		TestServiceFuncs fooapi.TestServiceFuncs
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
	mocktracer.New()
	tc := testcase.New().
		AddTask(10, func(w *Workspace) {
			r := apicommon.NewRouter()
			mt := mocktracer.New()
			w.MT = mt
			so := apicommon.ServerOptions{
				TraceIDGenerator: w.Input.TraceIDGenerator,
				Middlewares: apicommon.ServerMiddlewares{
					apicommon.AnyMethod: {
						NewForServer(mt),
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
			w.Output.MSs = w.MT.FinishedSpans()
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
				mss := w.Output.MSs
				if !assert.Len(w.T(), mss, 1) {
					w.T().FailNow()
				}
				ms := mss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ms.OperationName)
				assert.Equal(w.T(), map[string]interface{}{
					"span.kind":        ext.SpanKindRPCServerEnum,
					"component":        "JROH",
					"http.status_code": uint16(200),
					"jroh.error_code":  int32(0),
				}, ms.Tags())
				mlrs := ms.Logs()
				if assert.Len(w.T(), mlrs, 1) {
					mlr := mlrs[0]
					assert.Equal(w.T(), []mocktracer.MockKeyValue{
						{Key: "event", ValueKind: reflect.String, ValueString: "incoming rpc"},
						{Key: "trace_id", ValueKind: reflect.String, ValueString: "tid1"},
					}, mlr.Fields)
				}
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
				mss := w.Output.MSs
				if !assert.Len(w.T(), mss, 1) {
					w.T().FailNow()
				}
				ms := mss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ms.OperationName)
				assert.Equal(w.T(), map[string]interface{}{
					"span.kind":        ext.SpanKindRPCServerEnum,
					"component":        "JROH",
					"http.status_code": uint16(200),
					"jroh.error_code":  int32(-32603),
					"error":            true,
				}, ms.Tags())
				mlrs := ms.Logs()
				if assert.Len(w.T(), mlrs, 1) {
					mlr := mlrs[0]
					assert.Equal(w.T(), []mocktracer.MockKeyValue{
						{Key: "event", ValueKind: reflect.String, ValueString: "incoming rpc"},
						{Key: "trace_id", ValueKind: reflect.String, ValueString: "tid1"},
						{Key: "api_error", ValueKind: reflect.String, ValueString: "internal error"},
						{Key: "native_error", ValueKind: reflect.String, ValueString: "just wrong"},
					}, mlr.Fields)
				}
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
				mss := w.Output.MSs
				if !assert.Len(w.T(), mss, 1) {
					w.T().FailNow()
				}
				ms := mss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ms.OperationName)
				assert.Equal(w.T(), map[string]interface{}{
					"span.kind":        ext.SpanKindRPCServerEnum,
					"component":        "JROH",
					"http.status_code": uint16(500),
					"jroh.error_code":  int32(0),
					"error":            true,
				}, ms.Tags())
				mlrs := ms.Logs()
				if assert.Len(w.T(), mlrs, 1) {
					mlr := mlrs[0]
					assert.Equal(w.T(), []mocktracer.MockKeyValue{
						{Key: "event", ValueKind: reflect.String, ValueString: "incoming rpc"},
						{Key: "trace_id", ValueKind: reflect.String, ValueString: "tid1"},
						{Key: "native_error", ValueKind: reflect.String, ValueString: "resp encoding failed: json: error calling MarshalJSON for type *fooapi.MyStructString: bad word"},
					}, mlr.Fields)
				}
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
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
					"span.kind":        ext.SpanKindRPCServerEnum,
					"component":        "JROH",
					"http.status_code": uint16(200),
					"jroh.error_code":  int32(-32000),
				}, ms.Tags())
				mlrs := ms.Logs()
				if assert.Len(w.T(), mlrs, 1) {
					mlr := mlrs[0]
					assert.Equal(w.T(), []mocktracer.MockKeyValue{
						{Key: "event", ValueKind: reflect.String, ValueString: "incoming rpc"},
						{Key: "trace_id", ValueKind: reflect.String, ValueString: "tid1"},
						{Key: "api_error", ValueKind: reflect.String, ValueString: "not implemented"},
					}, mlr.Fields)
				}
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
				mss := w.Output.MSs
				if !assert.Len(w.T(), mss, 1) {
					w.T().FailNow()
				}
				ms := mss[0]
				assert.Equal(w.T(), "Foo.Test.DoSomething2", ms.OperationName)
				assert.Equal(w.T(), map[string]interface{}{
					"span.kind":        ext.SpanKindRPCServerEnum,
					"component":        "JROH",
					"http.status_code": uint16(200),
					"jroh.error_code":  int32(0),
				}, ms.Tags())
				mlrs := ms.Logs()
				if assert.Len(w.T(), mlrs, 1) {
					mlr := mlrs[0]
					assert.Equal(w.T(), []mocktracer.MockKeyValue{
						{Key: "event", ValueKind: reflect.String, ValueString: "incoming rpc"},
						{Key: "trace_id", ValueKind: reflect.String, ValueString: "tid1"},
						{Key: "params", ValueKind: reflect.String, ValueString: "{\n  \"myOnOff\": false\n}\n"},
						{Key: "resp", ValueKind: reflect.String, ValueString: "{\n  \"results\": {\n    \"myOnOff\": true\n  }\n}\n"},
					}, mlr.Fields)
				}
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServiceFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					ms := opentracing.SpanFromContext(ctx).(*mocktracer.MockSpan)
					assert.Equal(w.T(), "Foo.Test.DoSomething2", ms.OperationName)
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
					"span.kind":        ext.SpanKindRPCServerEnum,
					"component":        "JROH",
					"http.status_code": uint16(200),
					"jroh.error_code":  int32(0),
				}, ms.Tags())
				mlrs := ms.Logs()
				if assert.Len(w.T(), mlrs, 1) {
					mlr := mlrs[0]
					assert.Equal(w.T(), []mocktracer.MockKeyValue{
						{Key: "event", ValueKind: reflect.String, ValueString: "incoming rpc"},
						{Key: "trace_id", ValueKind: reflect.String, ValueString: "tid1"},
					}, mlr.Fields)
				}
			}),
	)
}
