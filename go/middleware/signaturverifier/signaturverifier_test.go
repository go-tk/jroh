package signaturverifier_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/go-tk/jroh/go/apicommon/testdata/fooapi"
	. "github.com/go-tk/jroh/go/middleware/signaturverifier"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
)

func TestSignatureVerifier(t *testing.T) {
	type Input struct {
		KeyFetcher      KeyFetcher
		OptionsBuilders []OptionsBuilder
		HeaderStr       string
	}
	type Output struct {
		RespBody string
		ErrStr   string
	}
	type Workspace struct {
		testcase.WorkspaceBase

		Input          Input
		ExpectedOutput Output
	}
	tc := testcase.New().
		AddTask(10, func(w *Workspace) {
			r := apicommon.NewRouter(nil)
			so := apicommon.ServerOptions{
				Middlewares: map[apicommon.MethodIndex][]apicommon.ServerMiddleware{
					apicommon.AnyMethod: {
						NewForServer(w.Input.KeyFetcher, w.Input.OptionsBuilders...),
					},
				},
				TraceIDGenerator: func() string { return "tid" },
			}
			fooapi.RegisterTestServer(&fooapi.TestServerFuncs{
				DoSomething3Func: func(context.Context) error { return nil },
			}, r, so)
			var output Output
			co := apicommon.ClientOptions{
				Transport: apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
					if headerStr := w.Input.HeaderStr; headerStr != "" {
						request.Header.Set("Authorization", headerStr)
					}
					responseRecorder := httptest.NewRecorder()
					r.ServeMux().ServeHTTP(responseRecorder, request.WithContext(context.Background()))
					output.RespBody = string(responseRecorder.Body.Bytes())
					response := responseRecorder.Result()
					return response, nil
				}),
			}
			tc := fooapi.NewTestClient("http://127.0.0.1", co)
			err := tc.DoSomething3(context.Background())
			if err != nil {
				output.ErrStr = err.Error()
			}
			assert.Equal(w.T(), w.ExpectedOutput, output)
		})
	testcase.RunList(
		t,
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.ExpectedOutput.ErrStr = `rpc failed; fullMethodName="Foo.Test.DoSomething3" traceID="tid": http request failed (3): unexpected status code - 401`
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.OptionsBuilders = []OptionsBuilder{MaxHeaderStrLength(3)}
				w.Input.HeaderStr = "1234"
				w.ExpectedOutput.RespBody = `signaturverifier: header str too long; headerStrLength=4 maxHeaderStrLength=3` + "\n"
				w.ExpectedOutput.ErrStr = `rpc failed; fullMethodName="Foo.Test.DoSomething3" traceID="tid": http request failed (3): unexpected status code - 422`
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.HeaderStr = "aslkdfjasldf"
				w.ExpectedOutput.RespBody = `signaturverifier: invalid header str; headerStr="aslkdfjasldf"` + "\n"
				w.ExpectedOutput.ErrStr = `rpc failed; fullMethodName="Foo.Test.DoSomething3" traceID="tid": http request failed (3): unexpected status code - 400`
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.HeaderStr = `Signature t=0,sid="",at="abc",s=""`
				w.ExpectedOutput.RespBody = `signaturverifier: unknown algorithm type; algorithmTypeStr="abc"` + "\n"
				w.ExpectedOutput.ErrStr = `rpc failed; fullMethodName="Foo.Test.DoSomething3" traceID="tid": http request failed (3): unexpected status code - 400`
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.OptionsBuilders = []OptionsBuilder{
					TimestampGetter(func() int64 { return 100 }),
					MaxTimestampSkew(5),
				}
				w.Input.HeaderStr = `Signature t=106,sid="",at="sha1",s=""`
				w.ExpectedOutput.RespBody = `signaturverifier: unexpected timestamp; timestamp=106 maxTimestampSkew=5` + "\n"
				w.ExpectedOutput.ErrStr = `rpc failed; fullMethodName="Foo.Test.DoSomething3" traceID="tid": http request failed (3): unexpected status code - 422`
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.KeyFetcher = func(ctx context.Context, senderID string) (key []byte, ok bool, err error) {
					return nil, false, errors.New("something wrong")
				}
				w.Input.OptionsBuilders = []OptionsBuilder{
					TimestampGetter(func() int64 { return 1234567890 }),
					MaxTimestampSkew(5),
				}
				w.Input.HeaderStr = `Signature t=1234567890,sid="",at="sha1",s=""`
				w.ExpectedOutput.ErrStr = `rpc failed; fullMethodName="Foo.Test.DoSomething3" traceID="tid": http request failed (3): unexpected status code - 500`
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				apicommon.DebugMode = true
				w.AddCleanup(func() { apicommon.DebugMode = false })
				w.Input.KeyFetcher = func(ctx context.Context, senderID string) (key []byte, ok bool, err error) {
					return nil, false, errors.New("something wrong")
				}
				w.Input.OptionsBuilders = []OptionsBuilder{
					TimestampGetter(func() int64 { return 1234567890 }),
					MaxTimestampSkew(5),
				}
				w.Input.HeaderStr = `Signature t=1234567890,sid="",at="sha1",s=""`
				w.ExpectedOutput.RespBody = `signaturverifier: key fetching failed: something wrong` + "\n"
				w.ExpectedOutput.ErrStr = `rpc failed; fullMethodName="Foo.Test.DoSomething3" traceID="tid": http request failed (3): unexpected status code - 500`
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.KeyFetcher = func(ctx context.Context, senderID string) (key []byte, ok bool, err error) {
					assert.Equal(w.T(), "user", senderID)
					return nil, false, nil
				}
				w.Input.OptionsBuilders = []OptionsBuilder{
					TimestampGetter(func() int64 { return 1234567890 }),
					MaxTimestampSkew(5),
				}
				w.Input.HeaderStr = `Signature t=1234567890,sid="user",at="sha1",s=""`
				w.ExpectedOutput.RespBody = `signaturverifier: key not found; senderID="user"` + "\n"
				w.ExpectedOutput.ErrStr = `rpc failed; fullMethodName="Foo.Test.DoSomething3" traceID="tid": http request failed (3): unexpected status code - 422`
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.KeyFetcher = func(ctx context.Context, senderID string) (key []byte, ok bool, err error) {
					assert.Equal(w.T(), "user", senderID)
					return []byte("pass"), true, nil
				}
				w.Input.OptionsBuilders = []OptionsBuilder{
					TimestampGetter(func() int64 { return 1234567890 }),
					MaxTimestampSkew(5),
				}
				w.Input.HeaderStr = `Signature t=1234567890,sid="user",at="sha1",s="123"`
				w.ExpectedOutput.RespBody = `signaturverifier: unexpected signature; signature="123"` + "\n"
				w.ExpectedOutput.ErrStr = `rpc failed; fullMethodName="Foo.Test.DoSomething3" traceID="tid": http request failed (3): unexpected status code - 422`
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.KeyFetcher = func(ctx context.Context, senderID string) (key []byte, ok bool, err error) {
					assert.Equal(w.T(), "foo_MD5", senderID)
					return []byte("bar_MD5"), true, nil
				}
				w.Input.OptionsBuilders = []OptionsBuilder{
					TimestampGetter(func() int64 { return 1234567890 }),
					MaxTimestampSkew(5),
				}
				w.Input.HeaderStr = `Signature t=1234567890,sid="foo_MD5",at="md5",s="QwnOmlfLbJh8d0jjvWLeIQ=="`
				w.ExpectedOutput.RespBody = `{}`
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.KeyFetcher = func(ctx context.Context, senderID string) (key []byte, ok bool, err error) {
					assert.Equal(w.T(), "foo_SHA1", senderID)
					return []byte("bar_SHA1"), true, nil
				}
				w.Input.OptionsBuilders = []OptionsBuilder{
					TimestampGetter(func() int64 { return 1234567890 }),
					MaxTimestampSkew(5),
				}
				w.Input.HeaderStr = `Signature t=1234567890,sid="foo_SHA1",at="sha1",s="bCMJKH66jFi6W6YNjDg/dM+ev+w="`
				w.ExpectedOutput.RespBody = `{}`
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.KeyFetcher = func(ctx context.Context, senderID string) (key []byte, ok bool, err error) {
					assert.Equal(w.T(), "foo_SHA256", senderID)
					return []byte("bar_SHA256"), true, nil
				}
				w.Input.OptionsBuilders = []OptionsBuilder{
					TimestampGetter(func() int64 { return 1234567890 }),
					MaxTimestampSkew(5),
				}
				w.Input.HeaderStr = `Signature t=1234567890,sid="foo_SHA256",at="sha256",s="CitHJeETGak9c7epf4E9crzF7E3r2CEHC7CvR9QR/T8="`
				w.ExpectedOutput.RespBody = `{}`
			}),
		tc.Copy().
			AddTask(9, func(w *Workspace) {
				w.Input.KeyFetcher = func(ctx context.Context, senderID string) (key []byte, ok bool, err error) {
					assert.Equal(w.T(), "foo_SHA512", senderID)
					return []byte("bar_SHA512"), true, nil
				}
				w.Input.OptionsBuilders = []OptionsBuilder{
					TimestampGetter(func() int64 { return 1234567890 }),
					MaxTimestampSkew(5),
				}
				w.Input.HeaderStr = `Signature t=1234567890,sid="foo_SHA512",at="sha512",s="8ybJd3HFkSYCMYlHMqFFbV5br931KBC1hHVGXcmiX6t7ymrJlZad/Vtue8zEovqzkPFX/QDI3wGqraYrTgAHDQ=="`
				w.ExpectedOutput.RespBody = `{}`
			}),
	)
}
