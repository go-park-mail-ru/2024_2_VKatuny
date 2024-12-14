package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Hits            *prometheus.CounterVec
	Timings         *prometheus.HistogramVec
	AuthHits        *prometheus.CounterVec
	AuthTimings     *prometheus.HistogramVec
	CompressHits    *prometheus.CounterVec
	CompressTimings *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		Hits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "hits",
				Help: "Number of requests with method, path and status code",
			},
			[]string{"method", "path", "status"},
		),
		Timings: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "timings",
				Help:    "Request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		AuthHits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_auth_hits",
				Help: "Number of requests with gRPC method and error",
			},
			[]string{"grpc_method", "error"},
		),
		AuthTimings: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_auth_timings",
				Help:    "Request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"grpc_method"},
		),
		CompressHits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_compress_hits",
				Help: "Number of requests with gRPC method and error",
			},
			[]string{"grpc_method", "error"},
		),
		CompressTimings: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_compress_timings",
				Help:    "Request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"grpc_method"},
		),
	}
}

func InitMetrics(metrics *Metrics) {
	prometheus.MustRegister(metrics.Hits)
	prometheus.MustRegister(metrics.Timings)
}

func InitAuthMetrics(metrics *Metrics) {
	prometheus.MustRegister(metrics.AuthHits)
	prometheus.MustRegister(metrics.AuthTimings)
}

func InitCompressMetrics(metrics *Metrics) {
	prometheus.MustRegister(metrics.CompressHits)
	prometheus.MustRegister(metrics.CompressTimings)
}
