package prometheusrf

import (
	"context"
	"strconv"
	"time"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	rpcDurationSecondsHistogramVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "jroh_client_rpc_duration_seconds",
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
			Name: "jroh_client_rpcs_total",
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
			Name:    "jroh_client_params_size_bytes",
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
			Name:    "jroh_client_resp_size_bytes",
			Buckets: []float64{50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000, 50000, 100000},
		},
		[]string{
			"jroh_namespace",
			"jroh_service_name",
			"jroh_method_name",
		},
	)
)

func NewForClient() apicommon.RPCHandler {
	return func(ctx context.Context, rpc *apicommon.RPC) error {
		t0 := time.Now()
		// Before
		err := rpc.Do(ctx)
		// After
		outgoingRPC := rpc.OutgoingRPC()
		if !outgoingRPC.RequestIsSent() {
			return err
		}
		t1 := time.Now()
		{
			observer := rpcDurationSecondsHistogramVec.WithLabelValues(
				outgoingRPC.Namespace(),
				outgoingRPC.ServiceName(),
				outgoingRPC.MethodName(),
			)
			observer.Observe(t1.Sub(t0).Seconds())
		}
		{
			var statusCodeStr string
			if statusCode := outgoingRPC.StatusCode(); statusCode == 0 {
				statusCodeStr = "-"
			} else {
				statusCodeStr = strconv.FormatInt(int64(statusCode), 10)
			}
			var errorCodeStr string
			if errorCode := outgoingRPC.Error().Code; errorCode == 0 {
				errorCodeStr = "-"
			} else {
				errorCodeStr = strconv.FormatInt(int64(errorCode), 10)
			}
			counter := rpcsTotalCounterVec.WithLabelValues(
				outgoingRPC.Namespace(),
				outgoingRPC.ServiceName(),
				outgoingRPC.MethodName(),
				statusCodeStr,
				errorCodeStr,
			)
			counter.Add(1)
		}
		if n := len(outgoingRPC.RawParams()); n >= 1 {
			observer := paramsSizeBytesHistogramVec.WithLabelValues(
				outgoingRPC.Namespace(),
				outgoingRPC.ServiceName(),
				outgoingRPC.MethodName(),
			)
			observer.Observe(float64(n))
		}
		if n := len(outgoingRPC.RawResp()); n >= 1 {
			observer := respSizeBytesHistogramVec.WithLabelValues(
				outgoingRPC.Namespace(),
				outgoingRPC.ServiceName(),
				outgoingRPC.MethodName(),
			)
			observer.Observe(float64(n))
		}
		return err
	}
}

func MustRegisterCollectors(registerer prometheus.Registerer) {
	registerer.MustRegister(
		rpcDurationSecondsHistogramVec,
		rpcsTotalCounterVec,
		paramsSizeBytesHistogramVec,
		respSizeBytesHistogramVec,
	)
}
