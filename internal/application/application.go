package application

import (
	"context"
	"log/slog"

	"github.com/FlutterDizaster/music-library/internal/application/config"
	"github.com/FlutterDizaster/music-library/internal/application/service"
	"github.com/FlutterDizaster/music-library/internal/infrastructure/http/detailsclient"
	"github.com/FlutterDizaster/music-library/internal/infrastructure/metrics"
	"github.com/FlutterDizaster/music-library/internal/infrastructure/persistance/migrator"
	"github.com/FlutterDizaster/music-library/internal/infrastructure/persistance/postgres"
	"github.com/FlutterDizaster/music-library/internal/presentation/handler"
	"github.com/FlutterDizaster/music-library/internal/presentation/server"
)

type Service interface {
	Start(context.Context) error
}

func New(ctx context.Context, settings config.Config) (Service, error) {
	slog.Debug("Application initialization", slog.String("stage", "running migrations"))
	err := migrator.RunMigrations(ctx, settings.DatabaseDSN, settings.MigrationsPath)
	if err != nil {
		return nil, err
	}

	slog.Debug("Application initialization", slog.String("stage", "creating metrics registry"))
	metricsRegistry := metrics.New("music-library")

	slog.Debug("Application initialization", slog.String("stage", "connecting to database"))
	songRepo, err := postgres.New(ctx, postgres.Settings{
		DatabaseDSN:  settings.DatabaseDSN,
		RetryCount:   settings.DBRetryCount,
		RetryBackoff: settings.DBRetryBackoff,
	})
	if err != nil {
		return nil, err
	}

	slog.Debug("Application initialization", slog.String("stage", "creating details server client"))
	detailsRepo := detailsclient.New(detailsclient.Settings{
		Addr:            settings.DetailsServerAddr,
		RetryCount:      settings.DetailsServerRetryCount,
		RetryBackoff:    settings.DetailsServerRetryBackoff,
		RetryMaxBackoff: settings.DetailsServerMaxRetryBackoff,
	})

	slog.Debug("Application initialization", slog.String("stage", "creating service"))
	service := service.New(service.Settings{
		MusicRepo:   songRepo,
		DetailsRepo: detailsRepo,
	})

	slog.Debug("Application initialization", slog.String("stage", "creating http handler"))
	handler := handler.New(service, metricsRegistry)

	slog.Debug("Application initialization", slog.String("stage", "creating server"))
	server := server.New(server.Settings{
		Addr:    settings.HTTPAddr,
		Handler: handler,
	})

	slog.Debug("Application initialization", slog.String("stage", "finished"))

	return server, nil
}
