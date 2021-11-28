package prometheusmw

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	errorsTotalCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jroh_server_errors_total",
		},
		[]string{
			"jroh_namespace",
			"jroh_service_name",
			"jroh_method_name",
		},
	)

	rpcDurationSecondsHistogramVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "jroh_server_rpc_duration_seconds",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{
			"jroh_namespace",
			"jroh_service_name",
			"jroh_method_name",
		},
	)

	rpcsTotalCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jroh_server_rpcs_total",
		},
		[]string{
			"jroh_namespace",
			"jroh_service_name",
			"jroh_method_name",
			"jroh_status_code",
			"jroh_error_code",
		},
	)

	paramsSizeBytesHistogramVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "jroh_server_params_size_bytes",
			Buckets: []float64{50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000, 50000, 100000},
		},
		[]string{
			"jroh_namespace",
			"jroh_service_name",
			"jroh_method_name",
		},
	)

	respSizeBytesHistogramVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "jroh_server_resp_size_bytes",
			Buckets: []float64{50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000, 50000, 100000},
		},
		[]string{
			"jroh_namespace",
			"jroh_service_name",
			"jroh_method_name",
		},
	)
)

func NewForServer() apicommon.ServerMiddleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			incomingRPC := apicommon.MustGetRPCFromContext(r.Context()).IncomingRPC()
			t0 := time.Now()
			// Before
			handler.ServeHTTP(w, r)
			// After
			t1 := time.Now()
			if incomingRPC.StatusCode()/100 == 5 || incomingRPC.Error().Code == apicommon.ErrorInternal {
				counter := errorsTotalCounterVec.WithLabelValues(
					incomingRPC.Namespace(),
					incomingRPC.ServiceName(),
					incomingRPC.MethodName(),
				)
				counter.Add(1)
			}
			{
				observer := rpcDurationSecondsHistogramVec.WithLabelValues(
					incomingRPC.Namespace(),
					incomingRPC.ServiceName(),
					incomingRPC.MethodName(),
				)
				observer.Observe(t1.Sub(t0).Seconds())
			}
			{
				counter := rpcsTotalCounterVec.WithLabelValues(
					incomingRPC.Namespace(),
					incomingRPC.ServiceName(),
					incomingRPC.MethodName(),
					strconv.FormatInt(int64(incomingRPC.StatusCode()), 10),
					strconv.FormatInt(int64(incomingRPC.Error().Code), 10),
				)
				counter.Add(1)
			}
			if n := len(incomingRPC.RawParams()); n >= 1 {
				observer := paramsSizeBytesHistogramVec.WithLabelValues(
					incomingRPC.Namespace(),
					incomingRPC.ServiceName(),
					incomingRPC.MethodName(),
				)
				observer.Observe(float64(n))
			}
			if n := len(incomingRPC.RawResp()); n >= 1 {
				observer := respSizeBytesHistogramVec.WithLabelValues(
					incomingRPC.Namespace(),
					incomingRPC.ServiceName(),
					incomingRPC.MethodName(),
				)
				observer.Observe(float64(n))
			}
		})
	}
}

func MustRegisterCollectors(registerer prometheus.Registerer) {
	registerer.MustRegister(
		errorsTotalCounterVec,
		rpcDurationSecondsHistogramVec,
		rpcsTotalCounterVec,
		paramsSizeBytesHistogramVec,
		respSizeBytesHistogramVec,
	)
}
