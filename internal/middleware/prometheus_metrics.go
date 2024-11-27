package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func PrometheusMiddleware(next http.Handler) http.Handler {
	var (
		httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "main_http_duration_seconds",
			Help: "Duration of HTTP requests.",
		}, []string{"path"})
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()

	})
}
