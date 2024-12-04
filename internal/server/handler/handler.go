package handler

import (
	"net/http"

	"github.com/FlutterDizaster/music-library/internal/server/middleware"

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
	router     *http.ServeMux
	controller MusicDataController
	registry   middleware.MetricsRegistry
}

// New returns a new Handler.
// Accepts a data controller and a registry for metrics.
func New(controller MusicDataController, registry middleware.MetricsRegistry) *Handler {
	h := &Handler{
		controller: controller,
		registry:   registry,
	}

	h.registerRoutes()

	return h
}

// ServeHTTP implements http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) registerRoutes() {
	router := http.NewServeMux()

	// Setup middlewares
	metricsMiddleware := middleware.NewMetricsMiddleware(h.registry)
	middlewareChain := middleware.MakeChain(
		middleware.Logger,
		metricsMiddleware.Handle,
	)

	// Setup routes
	musicDataRouter := newMusicDataHandler(h.controller)

	router.Handle("/api/v1", middlewareChain(http.StripPrefix("/api/v1", musicDataRouter)))
	router.Handle("GET /metrics", promhttp.Handler())

	h.router = router
}
