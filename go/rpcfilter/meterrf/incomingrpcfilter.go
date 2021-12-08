package meterrf

import (
	"context"
	"strconv"
	"time"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/prometheus/client_golang/prometheus"
)

func NewIncomingRPCFilter(registerer prometheus.Registerer) apicommon.IncomingRPCHandler {
	rpcsTotalCounterVec := prometheus.NewCounterVec(
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
	rpcDurationSecondsHistogramVec := prometheus.NewHistogramVec(
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
	errorsTotalCounterVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jroh_server_errors_total",
		},
		[]string{
			"jroh_namespace",
			"jroh_service_name",
			"jroh_method_name",
		},
	)
	paramsSizeBytesHistogramVec := prometheus.NewHistogramVec(
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
	resultsSizeBytesHistogramVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "jroh_server_results_size_bytes",
			Buckets: []float64{50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000, 50000, 100000},
		},
		[]string{
			"jroh_namespace",
			"jroh_service_name",
			"jroh_method_name",
		},
	)
	registerer.MustRegister(
		rpcsTotalCounterVec,
		rpcDurationSecondsHistogramVec,
		errorsTotalCounterVec,
		paramsSizeBytesHistogramVec,
		resultsSizeBytesHistogramVec,
	)
	return func(ctx context.Context, incomingRPC *apicommon.IncomingRPC) (returnedErr error) {
		t0 := time.Now()

		returnedErr = incomingRPC.Do(ctx)

		if returnedErr == nil {
			incomingRPC.EncodeResults()
		}
		rpcsTotalCounterVec.WithLabelValues(
			incomingRPC.Namespace,
			incomingRPC.ServiceName,
			incomingRPC.MethodName,
			strconv.Itoa(incomingRPC.StatusCode),
			strconv.Itoa(int(incomingRPC.ErrorCode)),
		).Add(1)
		rpcDurationSecondsHistogramVec.WithLabelValues(
			incomingRPC.Namespace,
			incomingRPC.ServiceName,
			incomingRPC.MethodName,
		).Observe(time.Since(t0).Seconds())
		if incomingRPC.StatusCode/100 == 5 {
			errorsTotalCounterVec.WithLabelValues(
				incomingRPC.Namespace,
				incomingRPC.ServiceName,
				incomingRPC.MethodName,
			).Add(1)
		}
		paramsSizeBytesHistogramVec.WithLabelValues(
			incomingRPC.Namespace,
			incomingRPC.ServiceName,
			incomingRPC.MethodName,
		).Observe(float64(len(incomingRPC.RawParams)))
		resultsSizeBytesHistogramVec.WithLabelValues(
			incomingRPC.Namespace,
			incomingRPC.ServiceName,
			incomingRPC.MethodName,
		).Observe(float64(len(incomingRPC.RawResults)))
		return
	}
}
