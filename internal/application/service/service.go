package service

import (
	"context"
	"strings"

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

var _ interfaces.MusicService = (*Service)(nil)

func New(settings Settings) *Service {
	return &Service{
		musicRepo:  settings.MusicRepo,
		detailRepo: settings.DetailsRepo,
	}
}

func (c *Service) GetLibrary(
	ctx context.Context,
	filters models.Filters,
) (models.Library, error) {
	return c.musicRepo.GetLibrary(ctx, filters)
}

func (c *Service) GetSongLyrics(
	ctx context.Context,
	id uuid.UUID,
	pagination models.Pagination,
) (models.Lyrics, error) {
	rawLyrics, err := c.musicRepo.GetSongLyrics(ctx, id)
	if err != nil {
		return models.Lyrics{}, err
	}

	verbsArr := strings.Split(rawLyrics, "\n\n")
	pagination.Total = len(verbsArr)

	if pagination.Limit == 0 {
		pagination.Limit = len(verbsArr)
	}

	if len(verbsArr) <= pagination.Offset {
		pagination.Limit = 0
		pagination.Offset = len(verbsArr)
		return models.Lyrics{
			Lyrics:     "",
			Pagination: pagination,
		}, nil
	}

	if len(verbsArr) < pagination.Limit+pagination.Offset {
		pagination.Limit = len(verbsArr) - pagination.Offset
	}

	lyrics := strings.Join(verbsArr[pagination.Offset:pagination.Offset+pagination.Limit], "\n\n")

	return models.Lyrics{
		Lyrics:     lyrics,
		Pagination: pagination,
	}, nil
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
