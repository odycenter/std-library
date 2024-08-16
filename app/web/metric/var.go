package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	GRPCActiveConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "grpc_active_connections",
		Help: "The current number of active gRPC connections",
	})
	GRPCMaxConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "grpc_max_connections",
		Help: "The maximum number of allowed gRPC connections",
	})
	GRPCConnectionAttempts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_connection_attempts_total",
		Help: "The total number of connection attempts",
	})
	GRPCConnectionRejections = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_connection_rejections_total",
		Help: "The total number of rejected connections due to limit",
	})
)

var (
	ExecutorRunning = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "executor_running",
			Help: "Number of running tasks in executor",
		},
		[]string{"name"},
	)
	ExecutorFree = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "executor_free",
			Help: "Number of free workers in executor",
		},
		[]string{"name"},
	)
	ExecutorWaiting = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "executor_waiting",
			Help: "Number of waiting tasks in executor",
		},
		[]string{"name"},
	)
)
