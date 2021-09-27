package outgoingrpclogger_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/go-tk/jroh/go/apicommon/testdata/fooapi"
	. "github.com/go-tk/jroh/go/rpcfilter/outgoingrpclogger"
	"github.com/go-tk/testcase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestOutgoingRPCLogger(t *testing.T) {
	type Input struct {
		TestServerFuncs  fooapi.TestServerFuncs
		TraceIDGenerator apicommon.TraceIDGenerator
		OptionsSetters   []OptionsSetter
		Params           fooapi.DoSomething2Params
	}
	type Output struct {
		Log string
	}
	type Workspace struct {
		testcase.WorkspaceBase

		Input          Input
		ExpectedOutput Output

		Buf bytes.Buffer
		TC  fooapi.TestClient
	}
	tc := testcase.New().
		AddTask(10, func(w *Workspace) {
			sm := http.NewServeMux()
			logger := zerolog.New(&w.Buf)
			so := apicommon.ServerOptions{
				TraceIDGenerator: w.Input.TraceIDGenerator,
			}
			fooapi.RegisterTestServer(&w.Input.TestServerFuncs, sm, so)
			co := apicommon.ClientOptions{
				RPCFilters: map[apicommon.MethodIndex][]apicommon.RPCHandler{
					apicommon.AnyMethod: {
						NewForClient(logger, w.Input.OptionsSetters...),
					},
				},
				Transport: apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
					responseRecorder := httptest.NewRecorder()
					sm.ServeHTTP(responseRecorder, request)
					response := responseRecorder.Result()
					return response, nil
				}),
			}
			w.TC = fooapi.NewTestClient("http://127.0.0.1", co)
			w.TC.DoSomething2(context.Background(), &w.Input.Params)
		}).
		AddTask(20, func(w *Workspace) {
			var output Output
			output.Log = w.Buf.String()
			assert.Equal(w.T(), w.ExpectedOutput, output)
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
				w.ExpectedOutput.Log = `{"level":"info","traceID":"tid1","url":"http://127.0.0.1/rpc/Foo.Test.DoSomething2",` +
					`"params":"{\"myOnOff\":false}","statusCode":200,"resp":"{\"traceID\":\"tid1\",` +
					`\"results\":{\"myOnOff\":true}}","message":"outgoing rpc"}` + "\n"
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					results.MyOnOff = true
					return nil
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
				w.ExpectedOutput.Log = `{"level":"info","traceID":"tid1","url":"http://127.0.0.1/rpc/Foo.Test.DoSomething2",` +
					`"params":"{\"myOnOff\":false}","statusCode":200,"resp":"{\"traceID\":\"tid1\",` +
					`\"results\":{\"myOnOff\":true}}","message":"outgoing rpc"}` + "\n" +
					`{"level":"info","traceID":"tid2","url":"http://127.0.0.1/rpc/Foo.Test.DoSomething2",` +
					`"params":"{\"myOnOff\":false}","statusCode":200,"resp":"{\"traceID\":\"tid2\",` +
					`\"results\":{\"myOnOff\":true}}","message":"outgoing rpc"}` + "\n"
			}).
			AddTask(19, func(w *Workspace) {
				w.TC.DoSomething2(context.Background(), &w.Input.Params)
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					results.MyOnOff = true
					return nil
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
				w.Input.OptionsSetters = []OptionsSetter{
					MaxRawParamsSize(10),
					MaxRawRespSize(11),
				}
				w.ExpectedOutput.Log = `{"level":"info","traceID":"tid1","url":"http://127.0.0.1/rpc/Foo.Test.DoSomething2",` +
					`"paramsSize":17,"truncatedParams":"{\"myOnOff\"","statusCode":200,"respSize":45,` +
					`"truncatedResp":"{\"traceID\":","message":"outgoing rpc"}` + "\n"
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					logger := zerolog.Ctx(ctx)
					logger.UpdateContext(func(context zerolog.Context) zerolog.Context {
						return context.Str("foo", "bar")
					})
					logger.Info().Msg("test")
					results.MyOnOff = true
					return nil
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
				w.ExpectedOutput.Log = `{"level":"info","foo":"bar","message":"test"}` + "\n" +
					`{"level":"info","foo":"bar","traceID":"tid1","url":"http://127.0.0.1/rpc/Foo.Test.DoSomething2",` +
					`"params":"{\"myOnOff\":false}","statusCode":200,"resp":"{\"traceID\":\"tid1\",` +
					`\"results\":{\"myOnOff\":true}}","message":"outgoing rpc"}` + "\n"
			}),
	)
}
