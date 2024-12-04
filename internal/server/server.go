package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	defaultShutdownMaxTime = 5 * time.Second
)

// Settings used to create Server.
// Settings must be provided to New function.
// All fields are required and cant be nil.
type Settings struct {
	Addr    string
	Handler http.Handler
}

// Server is a http server.
// Can be started with Start method.
// Gracefully stops when context is canceled.
// Must be initialized with New function.
type Server struct {
	server          *http.Server
	shutdownMaxTime time.Duration

	addr string
}

// New initializes and returns a new Server instance based on the provided settings.
// It configures the server address and handler, and sets the shutdown timeout to a default value.
func New(settings Settings) *Server {
	s := &Server{
		addr:            settings.Addr,
		shutdownMaxTime: defaultShutdownMaxTime,
	}

	s.server = &http.Server{
		Addr:              settings.Addr,
		Handler:           settings.Handler,
		ReadHeaderTimeout: time.Second,
	}

	return s
}

// Start starts the server. It returns an error if it fails to start the server.
// The error is not returned if the context is canceled, the server is shutdown
// and the shutdown is successful.
//
// The server is shutdown when the context is canceled. The shutdown is
// cancelled if the server failed to start.
func (s *Server) Start(ctx context.Context) error {
	eg, egCtx := errgroup.WithContext(ctx)

	errOnStart := false

	eg.Go(func() error {
		<-egCtx.Done()

		// Skip shutdown if error on start
		if errOnStart {
			return nil
		}

		slog.Info("Shutdown server", "addr", s.addr)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownMaxTime)
		defer cancel()

		return s.server.Shutdown(shutdownCtx)
	})

	eg.Go(func() error {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", slog.Any("error", err))
			errOnStart = true
			return err
		}
		return nil
	})

	return eg.Wait()
}
