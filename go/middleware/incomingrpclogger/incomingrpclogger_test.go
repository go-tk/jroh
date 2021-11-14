package incomingrpclogger_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/go-tk/jroh/go/apicommon/testdata/fooapi"
	. "github.com/go-tk/jroh/go/middleware/incomingrpclogger"
	"github.com/go-tk/testcase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestIncomingRPCLogger(t *testing.T) {
	type Input struct {
		TestServerFuncs  fooapi.TestServerFuncs
		OptionsBuilders  []OptionsBuilder
		TraceIDGenerator apicommon.TraceIDGenerator
		Params           fooapi.DoSomething2Params
	}
	type Output struct {
		Log string
	}
	type Workspace struct {
		testcase.WorkspaceBase

		Buf bytes.Buffer
		TC  fooapi.TestClient

		Input          Input
		ExpectedOutput Output
	}
	tc := testcase.New().
		AddTask(10, func(w *Workspace) {
			r := apicommon.NewRouter()
			logger := zerolog.New(&w.Buf)
			so := apicommon.ServerOptions{
				Middlewares: map[apicommon.MethodIndex][]apicommon.ServerMiddleware{
					apicommon.AnyMethod: {
						NewForServer(logger, w.Input.OptionsBuilders...),
					},
				},
				TraceIDGenerator: w.Input.TraceIDGenerator,
			}
			fooapi.RegisterTestServer(&w.Input.TestServerFuncs, r, so)
			co := apicommon.ClientOptions{
				Transport: apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
					request = request.Clone(request.Context())
					request.RemoteAddr = "127.0.0.1:12345"
					responseRecorder := httptest.NewRecorder()
					r.ServeHTTP(responseRecorder, request)
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
				w.ExpectedOutput.Log = `{"level":"info","traceID":"tid1","remoteAddr":"127.0.0.1:12345","fullMethodName":"Foo.Test.DoSomething2",` +
					`"params":"{\"myOnOff\":false}","resp":"{\"results\":{\"myOnOff\":true}}","statusCode":200,"message":"incoming rpc"}` + "\n"
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					results.MyOnOff = true
					return nil
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
				w.ExpectedOutput.Log = `{"level":"info","traceID":"tid1","remoteAddr":"127.0.0.1:12345","fullMethodName":"Foo.Test.DoSomething2",` +
					`"params":"{\"myOnOff\":false}","resp":"{\"results\":{\"myOnOff\":true}}","statusCode":200,"message":"incoming rpc"}` + "\n" +
					`{"level":"info","traceID":"tid2","remoteAddr":"127.0.0.1:12345","fullMethodName":"Foo.Test.DoSomething2",` +
					`"params":"{\"myOnOff\":false}","resp":"{\"results\":{\"myOnOff\":true}}","statusCode":200,"message":"incoming rpc"}` + "\n"
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
				w.Input.OptionsBuilders = []OptionsBuilder{
					MaxRawParamsSize(10),
					MaxRawRespSize(11),
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
				w.ExpectedOutput.Log = `{"level":"info","traceID":"tid1","remoteAddr":"127.0.0.1:12345","fullMethodName":"Foo.Test.DoSomething2",` +
					`"paramsSize":17,"truncatedParams":"{\"myOnOff\"","respSize":28,` +
					`"truncatedResp":"{\"results\":","statusCode":200,"message":"incoming rpc"}` + "\n"
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					tmp := fooapi.MyStructString{}
					tmp.TheStringA = "taboo"
					results.MyStructString = &tmp
					return nil
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
				w.ExpectedOutput.Log = `{"level":"error","traceID":"tid1","remoteAddr":"127.0.0.1:12345","fullMethodName":"Foo.Test.DoSomething2",` +
					`"params":"{\"myOnOff\":false}","err":"resp encoding failed: json: error calling MarshalJSON for type *fooapi.MyStructString: bad word",` +
					`"statusCode":500,"message":"incoming rpc"}` + "\n"
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					panic("hello")
				}
				lastTID := 0
				w.Input.TraceIDGenerator = func() string { lastTID++; return fmt.Sprintf("tid%d", lastTID) }
				w.ExpectedOutput.Log = `{"level":"error","traceID":"tid1","remoteAddr":"127.0.0.1:12345","fullMethodName":"Foo.Test.DoSomething2",` +
					`"params":"{\"myOnOff\":false}","errorCode":-32603,"err":"hello",` +
					`"stackTrace":"goroutine...","resp":"{\"error\":{\"code\":-32603,\"message\":\"internal error\"}}",` +
					`"statusCode":200,"message":"incoming rpc"}` + "\n"
			}).
			AddTask(19, func(w *Workspace) {
				s := w.Buf.String()
				ss := `,"stackTrace":"goroutine`
				i := strings.Index(s, ss)
				if !assert.GreaterOrEqual(w.T(), i, 0) {
					t.FailNow()
				}
				i += len(ss)
				w.Buf.Truncate(i)
				w.Buf.WriteString("...")
				j := strings.Index(s[i:], `"`)
				if !assert.GreaterOrEqual(w.T(), j, 0) {
					t.FailNow()
				}
				j += i
				w.Buf.WriteString(s[j:])
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
				w.ExpectedOutput.Log = `{"level":"info","traceID":"tid1","foo":"bar","message":"test"}` + "\n" +
					`{"level":"info","traceID":"tid1","foo":"bar","remoteAddr":"127.0.0.1:12345","fullMethodName":"Foo.Test.DoSomething2",` +
					`"params":"{\"myOnOff\":false}","resp":"{\"results\":{\"myOnOff\":true}}","statusCode":200,"message":"incoming rpc"}` + "\n"
			}),
	)
}
