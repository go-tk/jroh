package meterrf

import (
	"context"
	"strconv"
	"time"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	clientRPCsTotal = prometheus.NewCounterVec(
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

	clientRPCDurationSeconds = prometheus.NewHistogramVec(
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

	clientErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jroh_client_errors_total",
		},
		[]string{
			"jroh_namespace",
			"jroh_service_name",
			"jroh_method_name",
		},
	)

	clientParamsSizeBytes = prometheus.NewHistogramVec(
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

	clientResultsSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "jroh_client_results_size_bytes",
			Buckets: []float64{50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000, 50000, 100000},
		},
		[]string{
			"jroh_namespace",
			"jroh_service_name",
			"jroh_method_name",
		},
	)
)

func NewOutgoingRPCFilter(registerer prometheus.Registerer) apicommon.OutgoingRPCHandler {
	for _, collector := range []prometheus.Collector{
		clientRPCsTotal,
		clientRPCDurationSeconds,
		clientErrorsTotal,
		clientParamsSizeBytes,
		clientResultsSizeBytes,
	} {
		if err := registerer.Register(collector); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				panic(err)
			}
		}
	}
	return func(ctx context.Context, outgoingRPC *apicommon.OutgoingRPC) (returnedErr error) {
		t0 := time.Now()

		returnedErr = outgoingRPC.Do(ctx)

		if returnedErr == nil {
			outgoingRPC.ReadRawResults()
		}
		clientRPCsTotal.WithLabelValues(
			outgoingRPC.Namespace,
			outgoingRPC.ServiceName,
			outgoingRPC.MethodName,
			strconv.Itoa(outgoingRPC.StatusCode),
			strconv.Itoa(int(outgoingRPC.ErrorCode)),
		).Add(1)
		clientRPCDurationSeconds.WithLabelValues(
			outgoingRPC.Namespace,
			outgoingRPC.ServiceName,
			outgoingRPC.MethodName,
		).Observe(time.Since(t0).Seconds())
		if apicommon.ClientShouldReportError(returnedErr, outgoingRPC.StatusCode) {
			clientErrorsTotal.WithLabelValues(
				outgoingRPC.Namespace,
				outgoingRPC.ServiceName,
				outgoingRPC.MethodName,
			).Add(1)
		}
		clientParamsSizeBytes.WithLabelValues(
			outgoingRPC.Namespace,
			outgoingRPC.ServiceName,
			outgoingRPC.MethodName,
		).Observe(float64(len(outgoingRPC.RawParams)))
		clientResultsSizeBytes.WithLabelValues(
			outgoingRPC.Namespace,
			outgoingRPC.ServiceName,
			outgoingRPC.MethodName,
		).Observe(float64(len(outgoingRPC.RawResults)))
		return
	}
}
