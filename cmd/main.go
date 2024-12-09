package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/FlutterDizaster/music-library/internal/application"
	"github.com/FlutterDizaster/music-library/internal/application/config"
)

func main() {
	os.Exit(mainWithCode())
}

func mainWithCode() int {
	// Logger initialization
	slog.SetDefault(
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)

	// Gracefull shutdown context
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	// Application settings
	settings, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", slog.Any("error", err))
		return 1
	}

	slog.Debug("Initializing application")
	// Application initialization
	app, err := application.New(ctx, *settings)
	if err != nil {
		slog.Error("Application initialization failed")
		return 1
	}

	slog.Debug("Starting application")
	// Application start
	if err = app.Start(ctx); err != nil {
		slog.Error("Application start failed", slog.Any("error", err))
		return 1
	}

	return 0
}
