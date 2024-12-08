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

func (r *MetricsRegistry) IncrementActiveConnections() {
	r.httpActiveConnections.WithLabelValues(r.service).Inc()
}

func (r *MetricsRegistry) DecrementActiveConnections() {
	r.httpActiveConnections.WithLabelValues(r.service).Dec()
}

func (r *MetricsRegistry) IncrementRequestsTotal(method, status, endpoint string) {
	r.httpRequestsTotal.WithLabelValues(method, status, endpoint, r.service).Inc()
}

func (r *MetricsRegistry) ObserveRequestDuration(
	method, status, endpoint string,
	duration float64,
) {
	r.httpRequestDuration.WithLabelValues(method, status, endpoint, r.service).Observe(duration)
}

func (r *MetricsRegistry) ObserveRequestSize(method, status, endpoint string, size float64) {
	r.httpRequestSize.WithLabelValues(method, status, endpoint, r.service).Observe(size)
}

func (r *MetricsRegistry) ObserveResponseSize(method, status, endpoint string, size float64) {
	r.httpResponseSize.WithLabelValues(method, status, endpoint, r.service).Observe(size)
}
