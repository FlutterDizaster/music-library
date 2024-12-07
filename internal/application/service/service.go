package service

import (
	"context"

	"github.com/FlutterDizaster/music-library/internal/domain/interfaces"
	"github.com/FlutterDizaster/music-library/internal/domain/models"
	"github.com/google/uuid"
)

type Settings struct {
	MusicRepo   interfaces.MusicRepository
	DetailsRepo interfaces.DetailsRepository
}

type Service struct {
	musicRepo  interfaces.MusicRepository
	detailRepo interfaces.DetailsRepository
}

func New(settings Settings) *Service {
	return &Service{
		musicRepo:  settings.MusicRepo,
		detailRepo: settings.DetailsRepo,
	}
}

func (c *Service) GetLibrary(
	ctx context.Context,
	params map[string][]string,
) (models.Library, error) {
	// buildedParams, err := c.paramsBuilder.Build(params)
	// if err != nil {
	// 	return models.Library{}, err
	// }

	// library, err := c.musicRepo.GetLibrary(ctx, buildedParams)
	// if err != nil {
	// 	return models.Library{}, err
	// }

	// return library, nil
	return models.Library{}, nil
}

// FIXME: rewrite.
func (c *Service) GetSongLyrics(
	ctx context.Context,
	id uuid.UUID,
	params map[string][]string,
) (models.Lyrics, error) {
	// TODO: Get lyrics
	// Filter pagination
	// return result
	return models.Lyrics{}, nil
}

func (c *Service) AddSong(ctx context.Context, title models.SongTitle) (uuid.UUID, error) {
	details, err := c.detailRepo.GetSongDetails(ctx, title)
	if err != nil {
		return uuid.Nil, err
	}

	song := models.Song{
		SongTitle:  title,
		SongDetail: details,
	}

	id, err := c.musicRepo.AddSong(ctx, song)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (c *Service) DeleteSong(ctx context.Context, id uuid.UUID) error {
	return c.musicRepo.DeleteSong(ctx, id)
}

func (c *Service) UpdateSong(ctx context.Context, song models.Song) error {
	return c.musicRepo.UpdateSong(ctx, song)
}
