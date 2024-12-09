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

// New returns a new DetailsClient instance based on the provided settings.
// It configures the Resty client with the provided settings, sets the retry count and
// backoff time, and sets the base URL.
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

// GetSongDetails retrieves the details of a song with the given title from the details server.
//
// It logs debug messages before and after retrieving the details, and error messages if an error occurs.
//
// Parameters:
//   - ctx: The context for managing cancellation.
//   - title: The models.SongTitle of the song whose details are to be retrieved.
//
// Returns:
//   - The models.SongDetail of the song if successful.
//   - An error if the operation fails or if no details are found.
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
