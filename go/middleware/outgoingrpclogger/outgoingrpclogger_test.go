package outgoingrpclogger_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/go-tk/jroh/go/apicommon/testdata/fooapi"
	. "github.com/go-tk/jroh/go/middleware/outgoingrpclogger"
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

		Buffer bytes.Buffer
	}
	tc := testcase.New().
		AddTask(10, func(w *Workspace) {
			sm := http.NewServeMux()
			logger := zerolog.New(&w.Buffer)
			so := apicommon.ServerOptions{
				TraceIDGenerator: w.Input.TraceIDGenerator,
			}
			fooapi.RegisterTestServer(&w.Input.TestServerFuncs, sm, so)
			co := apicommon.ClientOptions{
				Middlewares: map[apicommon.MethodIndex][]apicommon.ClientMiddleware{
					apicommon.AnyMethod: {
						New(logger, w.Input.OptionsSetters...),
					},
				},
				Transport: apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
					responseRecorder := httptest.NewRecorder()
					sm.ServeHTTP(responseRecorder, request)
					response := responseRecorder.Result()
					return response, nil
				}),
			}
			tc := fooapi.NewTestClient("http://127.0.0.1", co)
			tc.DoSomething2(context.Background(), &w.Input.Params)
		}).
		AddTask(20, func(w *Workspace) {
			var output Output
			output.Log = w.Buffer.String()
			assert.Equal(w.T(), w.ExpectedOutput, output)
		})
	testcase.RunList(t,
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					results.MyOnOff = true
					return nil
				}
				w.Input.TraceIDGenerator = func() string { return "abc123" }
				w.ExpectedOutput.Log = `{"level":"info","url":"http://127.0.0.1/rpc/Foo.Test.DoSomething2",` +
					`"params":"{\"myOnOff\":false}","statusCode":200,"resp":"{\"traceID\":\"abc123\",\"results\":{\"myOnOff\":true}}",` +
					`"message":"outgoing rpc"}` + "\n"
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.TestServerFuncs.DoSomething2Func = func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
					results.MyOnOff = true
					return nil
				}
				w.Input.TraceIDGenerator = func() string { return "abc123" }
				w.Input.OptionsSetters = []OptionsSetter{
					MaxRawParamsSize(10),
					MaxRawRespSize(11),
				}
				w.ExpectedOutput.Log = `{"level":"info","url":"http://127.0.0.1/rpc/Foo.Test.DoSomething2",` +
					`"truncatedParams":"{\"myOnOff\"","statusCode":200,"truncatedResp":"{\"traceID\":",` +
					`"message":"outgoing rpc"}` + "\n"
			}),
	)
}
