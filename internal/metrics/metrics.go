package metrics

import "github.com/prometheus/client_golang/prometheus"

type HTTPMetricsRegistry struct {
	service string

	httpActiveConnections *prometheus.GaugeVec
	httpRequestsTotal     *prometheus.CounterVec
	httpRequestDuration   *prometheus.HistogramVec
	httpRequestSize       *prometheus.HistogramVec
	httpResponseSize      *prometheus.HistogramVec
}

func New(service string) *HTTPMetricsRegistry {
	httpLabels := []string{
		"method",   // http method
		"status",   // http status code
		"endpoint", // http endpoint
		"service",  // application name
	}

	r := &HTTPMetricsRegistry{
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

func (r *HTTPMetricsRegistry) IncrementActiveConnections() {
	r.httpActiveConnections.WithLabelValues(r.service).Inc()
}

func (r *HTTPMetricsRegistry) DecrementActiveConnections() {
	r.httpActiveConnections.WithLabelValues(r.service).Dec()
}

func (r *HTTPMetricsRegistry) IncrementRequestsTotal(method, status, endpoint string) {
	r.httpRequestsTotal.WithLabelValues(method, status, endpoint, r.service).Inc()
}

func (r *HTTPMetricsRegistry) ObserveRequestDuration(
	method, status, endpoint string,
	duration float64,
) {
	r.httpRequestDuration.WithLabelValues(method, status, endpoint, r.service).Observe(duration)
}

func (r *HTTPMetricsRegistry) ObserveRequestSize(method, status, endpoint string, size float64) {
	r.httpRequestSize.WithLabelValues(method, status, endpoint, r.service).Observe(size)
}

func (r *HTTPMetricsRegistry) ObserveResponseSize(method, status, endpoint string, size float64) {
	r.httpResponseSize.WithLabelValues(method, status, endpoint, r.service).Observe(size)
}
