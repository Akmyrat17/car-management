package middleware

import (
	// "github.com/armon/go-metrics/prometheus"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of http requests",
		},
		[]string{"path", "method"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_requests_duration_seconds",
			Help: "DUration of http requests in seconds",
		},
		[]string{"path", "method"},
	)

	statusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_total",
			Help: "Total Number of Http resp[onse by status code]",
		},
		[]string{"path", "method", "status_code"},
	)
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func MetricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()
		requestCounter.WithLabelValues(r.URL.Path, r.Method)
		requestDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration)
		statusCounter.WithLabelValues(r.URL.Path, r.Method, http.StatusText(ww.statusCode)).Inc()
	})
}

func (re *responseWriter) WriteHeader(status_code int) {
	re.statusCode = status_code
	re.ResponseWriter.WriteHeader(status_code)
}

func init() {
	prometheus.MustRegister(requestCounter, requestDuration, statusCounter)
}
