package prometheusrf_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/go-tk/jroh/go/apicommon/testdata/fooapi"
	. "github.com/go-tk/jroh/go/rpcfilter/prometheusrf"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusHelper(t *testing.T) {
	r := prometheus.NewRegistry()
	MustRegisterCollectors(r)
	rr := apicommon.NewRouter()
	so := apicommon.ServerOptions{}
	fooapi.RegisterTestServer(&fooapi.TestServerFuncs{
		DoSomethingFunc: func(context.Context, *fooapi.DoSomethingParams) error {
			return nil
		},
		DoSomething2Func: func(context.Context, *fooapi.DoSomething2Params, *fooapi.DoSomething2Results) error {
			return nil
		},
	}, rr, so)
	var transportErr error
	co := apicommon.ClientOptions{
		RPCFilters: map[apicommon.MethodIndex][]apicommon.RPCHandler{
			apicommon.AnyMethod: {
				NewForClient(),
			},
		},
		Transport: apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
			if transportErr != nil {
				return nil, transportErr
			}
			responseRecorder := httptest.NewRecorder()
			rr.ServeHTTP(responseRecorder, request)
			response := responseRecorder.Result()
			return response, nil
		}),
	}
	tc := fooapi.NewTestClient("http://127.0.0.1", co)
	tmp := fooapi.MyStructString{}
	tmp.TheStringA = "taboo"
	err := tc.DoSomething(context.Background(), &fooapi.DoSomethingParams{
		MyStructString: &tmp,
	})
	if !assert.Error(t, err) {
		t.FailNow()
	}
	err = tc.DoSomething(context.Background(), &fooapi.DoSomethingParams{})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	_, err = tc.DoSomething2(context.Background(), &fooapi.DoSomething2Params{})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	err = tc.DoSomething3(context.Background())
	if !assert.Error(t, err) {
		t.FailNow()
	}
	transportErr = errors.New("something wrong")
	err = tc.DoSomething3(context.Background())
	if !assert.ErrorIs(t, err, transportErr) {
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

	assert.Contains(t, s, `jroh_client_rpc_duration_seconds_bucket{jroh_method_name="DoSomething",jroh_namespace="Foo",jroh_service_name="Test",le="0.1"} 1`)
	assert.Contains(t, s, `jroh_client_rpc_duration_seconds_bucket{jroh_method_name="DoSomething2",jroh_namespace="Foo",jroh_service_name="Test",le="0.1"} 1`)
	assert.Contains(t, s, `jroh_client_rpc_duration_seconds_bucket{jroh_method_name="DoSomething3",jroh_namespace="Foo",jroh_service_name="Test",le="0.1"} 2`)

	assert.Contains(t, s, `jroh_client_rpcs_total{jroh_error_code="0",jroh_method_name="DoSomething",jroh_namespace="Foo",jroh_service_name="Test",jroh_status_code="200"} 1`)
	assert.Contains(t, s, `jroh_client_rpcs_total{jroh_error_code="0",jroh_method_name="DoSomething2",jroh_namespace="Foo",jroh_service_name="Test",jroh_status_code="200"} 1`)
	assert.Contains(t, s, `jroh_client_rpcs_total{jroh_error_code="-32000",jroh_method_name="DoSomething3",jroh_namespace="Foo",jroh_service_name="Test",jroh_status_code="200"} 1`)
	assert.Contains(t, s, `jroh_client_rpcs_total{jroh_error_code="0",jroh_method_name="DoSomething3",jroh_namespace="Foo",jroh_service_name="Test",jroh_status_code="0"} 1`)

	assert.Contains(t, s, `jroh_client_params_size_bytes_bucket{jroh_method_name="DoSomething",jroh_namespace="Foo",jroh_service_name="Test",le="1000"} 1`)
	assert.Contains(t, s, `jroh_client_params_size_bytes_bucket{jroh_method_name="DoSomething2",jroh_namespace="Foo",jroh_service_name="Test",le="1000"} 1`)
	assert.NotContains(t, s, `jroh_client_params_size_bytes_bucket{jroh_method_name="DoSomething3",jroh_namespace="Foo",jroh_service_name="Test",le="1000"} 1`)

	assert.Contains(t, s, `jroh_client_resp_size_bytes_bucket{jroh_method_name="DoSomething",jroh_namespace="Foo",jroh_service_name="Test",le="1000"} 1`)
	assert.Contains(t, s, `jroh_client_resp_size_bytes_bucket{jroh_method_name="DoSomething2",jroh_namespace="Foo",jroh_service_name="Test",le="1000"} 1`)
	assert.Contains(t, s, `jroh_client_resp_size_bytes_bucket{jroh_method_name="DoSomething3",jroh_namespace="Foo",jroh_service_name="Test",le="1000"} 1`)
}
