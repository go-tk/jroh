package signatureattachermw_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/go-tk/jroh/go/apicommon/testdata/fooapi"
	. "github.com/go-tk/jroh/go/middleware/signatureattachermw"
	"github.com/go-tk/jroh/go/middleware/signaturverifiermw"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
)

func TestForClient(t *testing.T) {
	type Input struct {
		SenderID        string
		Key             []byte
		OptionsBuilders []OptionsBuilder
	}
	type Workspace struct {
		testcase.WorkspaceBase

		Input Input
	}
	tc := testcase.New().
		AddTask(10, func(w *Workspace) {
			r := apicommon.NewRouter()
			so := apicommon.ServerOptions{
				Middlewares: apicommon.ServerMiddlewares{
					apicommon.AnyMethod: {
						signaturverifiermw.NewForServer(func(ctx context.Context, senderID string) (key []byte, ok bool, err error) {
							if senderID != w.Input.SenderID {
								return nil, false, nil
							}
							return w.Input.Key, true, nil
						}, signaturverifiermw.TimestampGetter(func() int64 { return 1234567890 })),
					},
				},
			}
			fooapi.RegisterTestService(&fooapi.TestServiceFuncs{
				DoSomething3Func: func(context.Context) error { return nil },
			}, r, so)
			obs := append(w.Input.OptionsBuilders, TimestampGetter(func() int64 { return 1234567890 }))
			co := apicommon.ClientOptions{
				Middlewares: apicommon.ClientMiddlewares{
					apicommon.AnyMethod: {
						NewForClient(w.Input.SenderID, w.Input.Key, obs...),
					},
				},
				Transport: apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
					a := request.Header.Get("Authorization")
					if !assert.NotEmpty(t, a) {
						t.FailNow()
					}
					w.T().Logf("Authorization: %s", a)
					responseRecorder := httptest.NewRecorder()
					r.ServeHTTP(responseRecorder, request.WithContext(context.Background()))
					response := responseRecorder.Result()
					return response, nil
				}),
			}
			tc := fooapi.NewTestClient("http://127.0.0.1", co)
			err := tc.DoSomething3(context.Background())
			assert.NoError(w.T(), err)
		})
	testcase.RunList(
		t,
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.SenderID = "hello"
				w.Input.Key = []byte("world")
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.SenderID = "foo_MD5"
				w.Input.Key = []byte("bar_MD5")
				w.Input.OptionsBuilders = []OptionsBuilder{AlgorithmMD5()}
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.SenderID = "foo_SHA1"
				w.Input.Key = []byte("bar_SHA1")
				w.Input.OptionsBuilders = []OptionsBuilder{AlgorithmSHA1()}
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.SenderID = "foo_SHA256"
				w.Input.Key = []byte("bar_SHA256")
				w.Input.OptionsBuilders = []OptionsBuilder{AlgorithmSHA256()}
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.SenderID = "foo_SHA512"
				w.Input.Key = []byte("bar_SHA512")
				w.Input.OptionsBuilders = []OptionsBuilder{AlgorithmSHA512()}
			}),
	)
}
