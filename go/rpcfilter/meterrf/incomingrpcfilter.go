package meterrf

import (
	"context"
	"strconv"
	"time"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	serverRPCsTotal = prometheus.NewCounterVec(
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

	serverRPCDurationSeconds = prometheus.NewHistogramVec(
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

	serverErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jroh_server_errors_total",
		},
		[]string{
			"jroh_namespace",
			"jroh_service_name",
			"jroh_method_name",
		},
	)

	serverParamsSizeBytes = prometheus.NewHistogramVec(
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

	serverResultsSizeBytes = prometheus.NewHistogramVec(
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
)

func NewIncomingRPCFilter(registerer prometheus.Registerer) apicommon.IncomingRPCHandler {
	for _, collector := range []prometheus.Collector{
		serverRPCsTotal,
		serverRPCDurationSeconds,
		serverErrorsTotal,
		serverParamsSizeBytes,
		serverResultsSizeBytes,
	} {
		if err := registerer.Register(collector); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				panic(err)
			}
		}
	}
	return func(ctx context.Context, incomingRPC *apicommon.IncomingRPC) (returnedErr error) {
		t0 := time.Now()

		returnedErr = incomingRPC.Do(ctx)

		if returnedErr == nil {
			incomingRPC.EncodeResults()
		}
		serverRPCsTotal.WithLabelValues(
			incomingRPC.Namespace,
			incomingRPC.ServiceName,
			incomingRPC.MethodName,
			strconv.Itoa(incomingRPC.StatusCode),
			strconv.Itoa(int(incomingRPC.ErrorCode)),
		).Add(1)
		serverRPCDurationSeconds.WithLabelValues(
			incomingRPC.Namespace,
			incomingRPC.ServiceName,
			incomingRPC.MethodName,
		).Observe(time.Since(t0).Seconds())
		if apicommon.ServerShouldReportError(returnedErr, incomingRPC.StatusCode) {
			serverErrorsTotal.WithLabelValues(
				incomingRPC.Namespace,
				incomingRPC.ServiceName,
				incomingRPC.MethodName,
			).Add(1)
		}
		serverParamsSizeBytes.WithLabelValues(
			incomingRPC.Namespace,
			incomingRPC.ServiceName,
			incomingRPC.MethodName,
		).Observe(float64(len(incomingRPC.RawParams)))
		serverResultsSizeBytes.WithLabelValues(
			incomingRPC.Namespace,
			incomingRPC.ServiceName,
			incomingRPC.MethodName,
		).Observe(float64(len(incomingRPC.RawResults)))
		return
	}
}
