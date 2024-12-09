package detailsclient

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/FlutterDizaster/music-library/internal/apperrors"
	"github.com/FlutterDizaster/music-library/internal/domain/interfaces"
	"github.com/FlutterDizaster/music-library/internal/domain/models"
	"github.com/go-resty/resty/v2"
)

type Settings struct {
	Addr string

	RetryCount      int
	RetryBackoff    time.Duration
	RetryMaxBackoff time.Duration
}

type DetailsClient struct {
	client *resty.Client
}

var _ interfaces.DetailsRepository = (*DetailsClient)(nil)

func New(settings Settings) *DetailsClient {
	slog.Debug("Creating details client", slog.String("addr", settings.Addr))
	client := resty.New().
		SetBaseURL(settings.Addr).
		SetRetryCount(settings.RetryCount).
		SetRetryWaitTime(settings.RetryBackoff).
		SetRetryMaxWaitTime(settings.RetryMaxBackoff)

	slog.Debug("Details client created")
	return &DetailsClient{
		client: client,
	}
}

func (c *DetailsClient) GetSongDetails(
	ctx context.Context,
	title models.SongTitle,
) (models.SongDetail, error) {
	slog.Debug(
		"Getting song details",
		slog.String("song", title.Song),
		slog.String("group", title.Group),
	)
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"group": title.Group,
			"song":  title.Song,
		}).
		SetHeader("Accept", "application/json").
		Get("/info")
	if err != nil {
		slog.Debug(
			"Error while getting song details",
			slog.String("song", title.Song),
			slog.String("group", title.Group),
			slog.Any("error", err),
		)
		return models.SongDetail{}, err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		details := models.SongDetail{}

		if err = details.UnmarshalJSON(resp.Body()); err != nil {
			slog.Error("Error while unmarshalling song details",
				slog.String("song", title.Song),
				slog.String("group", title.Group),
				slog.Any("error", err),
			)
			return models.SongDetail{}, apperrors.ErrDetailsServerBadResponse
		}

		return details, nil
	case http.StatusBadRequest:
		slog.Debug(
			"Bad details request",
			slog.String("song", title.Song),
			slog.String("group", title.Group),
		)
		return models.SongDetail{}, apperrors.ErrBadDetailsRequest
	default:
		slog.Error(
			"Error while getting song details",
			slog.String("song", title.Song),
			slog.String("group", title.Group),
			slog.String("status", resp.Status()),
		)
		return models.SongDetail{}, errors.New(resp.Status())
	}
}
