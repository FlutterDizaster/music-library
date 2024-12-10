package interfaces

import (
	"context"

	"github.com/FlutterDizaster/music-library/internal/domain/models"
)

type LibraryRepository interface {
	GetLibrary(ctx context.Context, filters models.Filters) (models.Library, error)
}
