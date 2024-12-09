package metrics

import "github.com/prometheus/client_golang/prometheus"

type MetricsRegistry struct {
	service string

	httpActiveConnections *prometheus.GaugeVec
	httpRequestsTotal     *prometheus.CounterVec
	httpRequestDuration   *prometheus.HistogramVec
	httpRequestSize       *prometheus.HistogramVec
	httpResponseSize      *prometheus.HistogramVec
}

// New returns a new MetricsRegistry instance based on the provided service name.
// It configures the following metrics:
// - http_active_connections: number of active HTTP connections
// - http_requests_total: number of HTTP requests
// - http_request_duration: duration of HTTP requests
// - http_request_size: size of HTTP requests
// - http_response_size: size of HTTP responses
//
// The metrics are registered with the global prometheus registry.
func New(service string) *MetricsRegistry {
	httpLabels := []string{
		"method",   // http method
		"status",   // http status code
		"endpoint", // http endpoint
		"service",  // application name
	}

	r := &MetricsRegistry{
		service: service,
		httpActiveConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_active_connections",
				Help: "Number of active HTTP connections",
			},
			[]string{"service"},
		),
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Number of HTTP requests",
			},
			httpLabels,
		),
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration",
				Help:    "Duration of HTTP requests",
				Buckets: prometheus.DefBuckets,
			},
			httpLabels,
		),
		httpRequestSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_size",
				Help:    "Size of HTTP requests",
				Buckets: prometheus.DefBuckets,
			},
			httpLabels,
		),
		httpResponseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size",
				Help:    "Size of HTTP responses",
				Buckets: prometheus.DefBuckets,
			},
			httpLabels,
		),
	}

	prometheus.MustRegister(r.httpActiveConnections)
	prometheus.MustRegister(r.httpRequestsTotal)
	prometheus.MustRegister(r.httpRequestDuration)
	prometheus.MustRegister(r.httpRequestSize)
	prometheus.MustRegister(r.httpResponseSize)

	return r
}

// IncrementActiveConnections increments the "http_active_connections" counter
// for the given service by one.
func (r *MetricsRegistry) IncrementActiveConnections() {
	r.httpActiveConnections.WithLabelValues(r.service).Inc()
}

// DecrementActiveConnections decrements the "http_active_connections" gauge
// for the given service by one.
func (r *MetricsRegistry) DecrementActiveConnections() {
	r.httpActiveConnections.WithLabelValues(r.service).Dec()
}

// IncrementRequestsTotal increments the "http_requests_total" counter
// for the specified HTTP method, status, and endpoint within the context
// of the registered service.
func (r *MetricsRegistry) IncrementRequestsTotal(method, status, endpoint string) {
	r.httpRequestsTotal.WithLabelValues(method, status, endpoint, r.service).Inc()
}

// ObserveRequestDuration observes the duration of an HTTP request.
// method, status, and endpoint are labels that identify the request.
// duration is the duration of the request in seconds.
func (r *MetricsRegistry) ObserveRequestDuration(
	method, status, endpoint string,
	duration float64,
) {
	r.httpRequestDuration.WithLabelValues(method, status, endpoint, r.service).Observe(duration)
}

// ObserveRequestSize records the size of an HTTP request.
// method, status, and endpoint are labels that identify the request.
// size is the size of the request in bytes.
func (r *MetricsRegistry) ObserveRequestSize(method, status, endpoint string, size float64) {
	r.httpRequestSize.WithLabelValues(method, status, endpoint, r.service).Observe(size)
}

// ObserveResponseSize records the size of an HTTP response.
// method, status, and endpoint are labels that identify the request.
// size is the size of the response in bytes.
func (r *MetricsRegistry) ObserveResponseSize(method, status, endpoint string, size float64) {
	r.httpResponseSize.WithLabelValues(method, status, endpoint, r.service).Observe(size)
}
