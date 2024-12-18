package interfaces

import (
	"context"

	"github.com/FlutterDizaster/music-library/internal/domain/models"
)

type DetailsClient interface {
	GetSongDetails(ctx context.Context, title models.SongTitle) (models.SongDetail, error)
}
