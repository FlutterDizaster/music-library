package controller

import (
	"context"

	"github.com/FlutterDizaster/music-library/internal/abstraction"
	"github.com/FlutterDizaster/music-library/internal/models"
	"github.com/google/uuid"
)

type MusicLibraryRepository interface {
	GetLibrary(ctx context.Context, params abstraction.Params) (models.Library, error)
	GetSongLyrics(
		ctx context.Context,
		id uuid.UUID,
		params abstraction.Params,
	) (models.Lyrics, error)
	AddSong(ctx context.Context, song models.Song) (uuid.UUID, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
	UpdateSong(ctx context.Context, song models.Song) error
}

type DetailsRepository interface {
	GetSongDetails(ctx context.Context, title models.SongTitle) (models.SongDetail, error)
}

type Settings struct {
	MusicRepo     MusicLibraryRepository
	DetailsRepo   DetailsRepository
	ParamsBuilder abstraction.ParamsBuilder
}

type Controller struct {
	musicRepo     MusicLibraryRepository
	detailRepo    DetailsRepository
	paramsBuilder abstraction.ParamsBuilder
}

func New(settings Settings) *Controller {
	return &Controller{
		musicRepo:     settings.MusicRepo,
		detailRepo:    settings.DetailsRepo,
		paramsBuilder: settings.ParamsBuilder,
	}
}

func (c *Controller) GetLibrary(
	ctx context.Context,
	params map[string][]string,
) (models.Library, error) {
	buildedParams, err := c.paramsBuilder.Build(params)
	if err != nil {
		return models.Library{}, err
	}

	library, err := c.musicRepo.GetLibrary(ctx, buildedParams)
	if err != nil {
		return models.Library{}, err
	}

	return library, nil
}

func (c *Controller) GetSongLyrics(
	ctx context.Context,
	id uuid.UUID,
	params map[string][]string,
) (models.Lyrics, error) {
	buildedParams, err := c.paramsBuilder.Build(params)
	if err != nil {
		return models.Lyrics{}, err
	}

	lyrics, err := c.musicRepo.GetSongLyrics(ctx, id, buildedParams)
	if err != nil {
		return models.Lyrics{}, err
	}

	return lyrics, nil
}

func (c *Controller) AddSong(ctx context.Context, title models.SongTitle) (uuid.UUID, error) {
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

func (c *Controller) DeleteSong(ctx context.Context, id uuid.UUID) error {
	return c.musicRepo.DeleteSong(ctx, id)
}

func (c *Controller) UpdateSong(ctx context.Context, song models.Song) error {
	return c.musicRepo.UpdateSong(ctx, song)
}
