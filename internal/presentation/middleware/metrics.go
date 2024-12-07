package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/FlutterDizaster/music-library/internal/domain/interfaces"
)

// Metrics is a statefull middleware that records metrics.
// Must be created with NewMetricsMiddleware function.
type Metrics struct {
	registry interfaces.HTTPMetricsRegistry
}

// NewMetricsMiddleware returns a new Metrics middleware.
func NewMetricsMiddleware(registry interfaces.HTTPMetricsRegistry) *Metrics {
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
