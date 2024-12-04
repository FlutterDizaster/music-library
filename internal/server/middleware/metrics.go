package middleware

import (
	"net/http"
	"strconv"
	"time"
)

// MetricsRegistry is a registry for metrics.
type MetricsRegistry interface {
	DecrementActiveConnections()
	IncrementActiveConnections()
	IncrementRequestsTotal(method string, status string, endpoint string)
	ObserveRequestDuration(method string, status string, endpoint string, duration float64)
	ObserveRequestSize(method string, status string, endpoint string, size float64)
	ObserveResponseSize(method string, status string, endpoint string, size float64)
}

// Metrics is a statefull middleware that records metrics.
// Must be created with NewMetricsMiddleware function.
type Metrics struct {
	registry MetricsRegistry
}

// NewMetricsMiddleware returns a new Metrics middleware.
func NewMetricsMiddleware(registry MetricsRegistry) *Metrics {
	return &Metrics{
		registry: registry,
	}
}

// Handle is a middleware that records metrics.
func (m *Metrics) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.registry.IncrementActiveConnections()
		defer m.registry.DecrementActiveConnections()

		start := time.Now()

		mw := memoryWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			size:           0,
		}
		next.ServeHTTP(&mw, r)

		m.registry.ObserveRequestDuration(
			r.Method,
			strconv.Itoa(mw.statusCode),
			r.URL.Path,
			time.Since(start).Seconds(),
		)

		m.registry.ObserveRequestSize(
			r.Method,
			strconv.Itoa(mw.statusCode),
			r.URL.Path,
			float64(mw.size),
		)

		m.registry.IncrementRequestsTotal(
			r.Method,
			strconv.Itoa(mw.statusCode),
			r.URL.Path,
		)
	})
}
