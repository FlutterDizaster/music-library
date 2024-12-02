package application

import (
	"context"
)

type Service interface {
	Start(context.Context) error
}

type Settings struct {
	// App settings
}

func New(ctx context.Context, settings Settings) (Service, error) {
	// Init application dependencies

	// Assemble application service

	// Return application service
	return nil, nil
}
