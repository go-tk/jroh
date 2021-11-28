package prometheusmw_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/go-tk/jroh/go/apicommon/testdata/fooapi"
	. "github.com/go-tk/jroh/go/middleware/prometheusmw"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
)

func TestForServer(t *testing.T) {
	r := prometheus.NewRegistry()
	MustRegisterCollectors(r)
	rr := apicommon.NewRouter()
	so := apicommon.ServerOptions{
		Middlewares: apicommon.ServerMiddlewares{
			apicommon.AnyMethod: {
				NewForServer(),
			},
		},
	}
	fooapi.RegisterTestService(&fooapi.TestServiceFuncs{
		DoSomethingFunc: func(context.Context, *fooapi.DoSomethingParams) error {
			return nil
		},
		DoSomething2Func: func(_ context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
			if params.MyOnOff {
				return errors.New("somthing bad")
			}
			return nil
		},
	}, rr, so)
	co := apicommon.ClientOptions{
		Transport: apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
			responseRecorder := httptest.NewRecorder()
			rr.ServeHTTP(responseRecorder, request)
			response := responseRecorder.Result()
			return response, nil
		}),
	}
	tc := fooapi.NewTestClient("http://127.0.0.1", co)
	err := tc.DoSomething(context.Background(), &fooapi.DoSomethingParams{})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	_, err = tc.DoSomething2(context.Background(), &fooapi.DoSomething2Params{})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	_, err = tc.DoSomething2(context.Background(), &fooapi.DoSomething2Params{MyOnOff: true})
	if !assert.Error(t, err) {
		t.FailNow()
	}
	err = tc.DoSomething3(context.Background())
	if !assert.Error(t, err) {
		t.FailNow()
	}
	mfs, err := r.Gather()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	var buf bytes.Buffer
	enc := expfmt.NewEncoder(&buf, expfmt.FmtText)
	for _, mf := range mfs {
		enc.Encode(mf)
	}
	s := buf.String()

	assert.Contains(t, s, `jroh_server_errors_total{jroh_method_name="DoSomething2",jroh_namespace="Foo",jroh_service_name="Test"} 1`+"\n")

	assert.Contains(t, s, `jroh_server_rpc_duration_seconds_bucket{jroh_method_name="DoSomething",jroh_namespace="Foo",jroh_service_name="Test",le="0.1"} 1`+"\n")
	assert.Contains(t, s, `jroh_server_rpc_duration_seconds_bucket{jroh_method_name="DoSomething2",jroh_namespace="Foo",jroh_service_name="Test",le="0.1"} 2`+"\n")
	assert.Contains(t, s, `jroh_server_rpc_duration_seconds_bucket{jroh_method_name="DoSomething3",jroh_namespace="Foo",jroh_service_name="Test",le="0.1"} 1`+"\n")

	assert.Contains(t, s, `jroh_server_rpcs_total{jroh_error_code="0",jroh_method_name="DoSomething",jroh_namespace="Foo",jroh_service_name="Test",jroh_status_code="200"} 1`+"\n")
	assert.Contains(t, s, `jroh_server_rpcs_total{jroh_error_code="0",jroh_method_name="DoSomething2",jroh_namespace="Foo",jroh_service_name="Test",jroh_status_code="200"} 1`+"\n")
	assert.Contains(t, s, `jroh_server_rpcs_total{jroh_error_code="-32603",jroh_method_name="DoSomething2",jroh_namespace="Foo",jroh_service_name="Test",jroh_status_code="200"} 1`+"\n")
	assert.Contains(t, s, `jroh_server_rpcs_total{jroh_error_code="-32000",jroh_method_name="DoSomething3",jroh_namespace="Foo",jroh_service_name="Test",jroh_status_code="200"} 1`+"\n")

	assert.Contains(t, s, `jroh_server_params_size_bytes_bucket{jroh_method_name="DoSomething",jroh_namespace="Foo",jroh_service_name="Test",le="1000"} 1`+"\n")
	assert.Contains(t, s, `jroh_server_params_size_bytes_bucket{jroh_method_name="DoSomething2",jroh_namespace="Foo",jroh_service_name="Test",le="1000"} 2`+"\n")
	assert.NotContains(t, s, `jroh_server_params_size_bytes_bucket{jroh_method_name="DoSomething3",jroh_namespace="Foo",jroh_service_name="Test",le="1000"} 1`+"\n")

	assert.Contains(t, s, `jroh_server_resp_size_bytes_bucket{jroh_method_name="DoSomething",jroh_namespace="Foo",jroh_service_name="Test",le="1000"} 1`+"\n")
	assert.Contains(t, s, `jroh_server_resp_size_bytes_bucket{jroh_method_name="DoSomething2",jroh_namespace="Foo",jroh_service_name="Test",le="1000"} 2`+"\n")
	assert.Contains(t, s, `jroh_server_resp_size_bytes_bucket{jroh_method_name="DoSomething3",jroh_namespace="Foo",jroh_service_name="Test",le="1000"} 1`+"\n")
}
