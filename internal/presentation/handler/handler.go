package handler

import (
	"log/slog"
	"net/http"

	"github.com/FlutterDizaster/music-library/internal/domain/interfaces"
	"github.com/FlutterDizaster/music-library/internal/infrastructure/metrics"
	"github.com/FlutterDizaster/music-library/internal/presentation/middleware"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//	@Title				Music Library API
//	@Version			1.0
//	@Description		This is the API for managing music library data, including songs, lyrics, and related operations.
//
//	@Contact.name		Dmitriy Loginov
//	@Contact.email		dmitriy@loginoff.space
//
//	@License.name		MIT
//	@License.url		https://opensource.org/licenses/MIT
//
//	@BasePath			/api/v1
//
//	@Tag.Name			Songs
//	@Tag.Description	Operations about songs

// Handler is a general HTTP handler.
// Must be created with New function.
type Handler struct {
	router   *http.ServeMux
	service  interfaces.MusicService
	registry metrics.HTTPMetricsRegistry
}

// New returns a new Handler.
// Accepts a data service and a registry for metrics.
func New(service interfaces.MusicService, registry metrics.HTTPMetricsRegistry) *Handler {
	slog.Debug("Creating handler")
	h := &Handler{
		service:  service,
		registry: registry,
	}

	h.registerRoutes()

	slog.Debug("Handler created")
	return h
}

// ServeHTTP implements http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Handling request in general handler", slog.String("path", r.URL.Path))
	h.router.ServeHTTP(w, r)
}

func (h *Handler) registerRoutes() {
	router := http.NewServeMux()

	slog.Debug("Creating handler", slog.String("stage", "registering middlewares"))
	// Setup middlewares
	metricsMiddleware := middleware.NewMetricsMiddleware(h.registry)
	middlewareChain := middleware.MakeChain(
		middleware.Logger,
		metricsMiddleware.Handle,
	)

	slog.Debug("Creating handler", slog.String("stage", "registering routes"))
	// Setup routes
	musicDataRouter := newMusicHandler(h.service)

	router.Handle("/api/v1/", middlewareChain(http.StripPrefix("/api/v1", musicDataRouter)))
	router.Handle("GET /metrics", promhttp.Handler())

	h.router = router
	slog.Debug("Creating handler", slog.String("stage", "routes registered"))
}
