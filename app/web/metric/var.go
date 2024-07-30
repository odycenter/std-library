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
