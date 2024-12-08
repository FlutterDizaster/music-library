package interfaces

import (
	"context"

	"github.com/FlutterDizaster/music-library/internal/domain/models"
	"github.com/google/uuid"
)

type MusicRepository interface {
	GetLibrary(ctx context.Context, filters models.Filters) (models.Library, error)
	GetSongLyrics(ctx context.Context, id uuid.UUID) (string, error)
	AddSong(ctx context.Context, song models.Song) (uuid.UUID, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
	UpdateSong(ctx context.Context, song models.Song) error
}
