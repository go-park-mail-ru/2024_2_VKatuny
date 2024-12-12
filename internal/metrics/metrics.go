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
	}
}

func InitMetrics(metrics *Metrics) {
	prometheus.MustRegister(metrics.Hits)
}
