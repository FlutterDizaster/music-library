package interfaces

import (
	"context"

	"github.com/FlutterDizaster/music-library/internal/domain/models"
	"github.com/google/uuid"
)

// MusicService is an interface for music data operations.
type MusicService interface {
	GetLibrary(ctx context.Context, filters models.Filters) (models.Library, error)
	GetSongLyrics(
		ctx context.Context,
		id uuid.UUID,
		pagination models.Pagination,
	) (models.Lyrics, error)
	AddSong(ctx context.Context, title models.SongTitle) (uuid.UUID, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
	UpdateSong(ctx context.Context, song models.Song) error
}
