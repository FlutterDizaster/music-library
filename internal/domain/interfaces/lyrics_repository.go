package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type LyricsRepository interface {
	GetLyrics(ctx context.Context, id uuid.UUID) (string, error)
}
