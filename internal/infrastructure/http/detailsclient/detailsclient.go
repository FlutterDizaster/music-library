package detailsclient

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/FlutterDizaster/music-library/internal/apperrors"
	"github.com/FlutterDizaster/music-library/internal/domain/interfaces"
	"github.com/FlutterDizaster/music-library/internal/domain/models"
	"github.com/go-resty/resty/v2"
)

const (
	retryCount       = 3
	retryWaitTime    = 1 * time.Second
	retryMaxWaitTime = 10 * time.Second
)

type DetailsClient struct {
	client *resty.Client
}

var _ interfaces.DetailsRepository = (*DetailsClient)(nil)

func New(addr string) *DetailsClient {
	client := resty.New().
		SetBaseURL(addr).
		SetRetryCount(retryCount).
		SetRetryWaitTime(retryWaitTime).
		SetRetryMaxWaitTime(retryMaxWaitTime)

	return &DetailsClient{
		client: client,
	}
}

func (c *DetailsClient) GetSongDetails(
	ctx context.Context,
	title models.SongTitle,
) (models.SongDetail, error) {
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"group": title.Group,
			"song":  title.Song,
		}).
		SetHeader("Accept", "application/json").
		Get("/info")
	if err != nil {
		return models.SongDetail{}, err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		details := models.SongDetail{}

		if err = details.UnmarshalJSON(resp.Body()); err != nil {
			return models.SongDetail{}, apperrors.ErrDetailsServerBadResponse
		}

		return details, nil
	case http.StatusBadRequest:
		return models.SongDetail{}, apperrors.ErrBadDetailsRequest
	default:
		return models.SongDetail{}, errors.New(resp.Status())
	}
}
