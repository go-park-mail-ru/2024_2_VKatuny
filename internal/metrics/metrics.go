package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Hits   *prometheus.CounterVec
	// Errors *prometheus.CounterVec
	Timings *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		Hits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "hits",
				Help: "Number of requests with method. path and status code",
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
	}
}

func InitMetrics(metrics *Metrics) {
	prometheus.MustRegister(metrics.Hits)
	prometheus.MustRegister(metrics.Timings)
}
