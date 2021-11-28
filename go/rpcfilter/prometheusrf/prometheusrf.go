package prometheusrf

import (
	"context"
	"strconv"
	"time"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	errorsTotalCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jroh_client_errors_total",
		},
		[]string{
			"jroh_namespace",
			"jroh_service_name",
			"jroh_method_name",
		},
	)

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
	return func(ctx context.Context, rpc *apicommon.RPC) (returnedErr error) {
		outgoingRPC := rpc.OutgoingRPC()
		t0 := time.Now()
		// Before
		returnedErr = outgoingRPC.Do(ctx)
		// After
		var err error
		switch returnedErr.(type) {
		case *apicommon.UnexpectedStatusCodeError, *apicommon.Error:
		default:
			err = returnedErr
		}
		if err != nil {
			counter := errorsTotalCounterVec.WithLabelValues(
				outgoingRPC.Namespace(),
				outgoingRPC.ServiceName(),
				outgoingRPC.MethodName(),
			)
			counter.Add(1)
		}
		if !outgoingRPC.IsRequested() {
			return
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
			counter := rpcsTotalCounterVec.WithLabelValues(
				outgoingRPC.Namespace(),
				outgoingRPC.ServiceName(),
				outgoingRPC.MethodName(),
				strconv.FormatInt(int64(outgoingRPC.StatusCode()), 10),
				strconv.FormatInt(int64(outgoingRPC.Error().Code), 10),
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
		return
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
