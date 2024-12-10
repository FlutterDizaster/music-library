package interfaces

import (
	"context"

	"github.com/FlutterDizaster/music-library/internal/domain/models"
	"github.com/google/uuid"
)

type SongsRepository interface {
	AddSong(ctx context.Context, song models.Song) (uuid.UUID, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
	UpdateSong(ctx context.Context, song models.Song) error
}
