package interfaces

import (
	"context"

	"github.com/FlutterDizaster/music-library/internal/domain/models"
	"github.com/google/uuid"
)

// MusicService is an interface for music data operations.
type MusicService interface {
	GetLibrary(ctx context.Context, params map[string][]string) (models.Library, error)
	GetSongLyrics(
		ctx context.Context,
		id uuid.UUID,
		params map[string][]string,
	) (models.Lyrics, error)
	AddSong(ctx context.Context, title models.SongTitle) (uuid.UUID, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
	UpdateSong(ctx context.Context, song models.Song) error
}
