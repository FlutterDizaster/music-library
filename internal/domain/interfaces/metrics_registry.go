package interfaces

// MetricsRegistry is a registry for metrics.
type HTTPMetricsRegistry interface {
	DecrementActiveConnections()
	IncrementActiveConnections()
	IncrementRequestsTotal(method string, status string, endpoint string)
	ObserveRequestDuration(method string, status string, endpoint string, duration float64)
	ObserveRequestSize(method string, status string, endpoint string, size float64)
	ObserveResponseSize(method string, status string, endpoint string, size float64)
}
