package observability

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP request counter
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// Request duration histogram
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Active WebSocket connections
	activeWebsocketConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_active_connections",
			Help: "Number of active WebSocket connections",
		},
	)
)

// InitMetrics registers Prometheus metrics.
func InitMetrics() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(activeWebsocketConnections)
}

// MetricsHandler returns an HTTP handler for Prometheus metrics.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// IncHTTPRequest increments the HTTP request counter.
func IncHTTPRequest(method, path, status string) {
	httpRequestsTotal.WithLabelValues(method, path, status).Inc()
}

// ObserveHTTPRequestDuration records request duration.
func ObserveHTTPRequestDuration(method, path string, duration float64) {
	httpRequestDuration.WithLabelValues(method, path).Observe(duration)
}

// SetActiveWebsocketConnections sets the active WebSocket connections gauge.
func SetActiveWebsocketConnections(count float64) {
	activeWebsocketConnections.Set(count)
}
