package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/FlutterDizaster/music-library/internal/application"
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
	ctx, cancle := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancle()

	// Application settings
	settings := application.Settings{}

	// TODO: Load config to settings struct

	// Application initialization
	app, err := application.New(ctx, settings)
	if err != nil {
		slog.Error("Application initialization failed", slog.Any("error", err))
		return 1
	}

	// Application start
	if err = app.Start(ctx); err != nil {
		slog.Error("Application start failed", slog.Any("error", err))
		return 1
	}

	return 0
}
