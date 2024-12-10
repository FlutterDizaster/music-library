package service

import (
	"context"
	"strings"

	"github.com/FlutterDizaster/music-library/internal/domain/interfaces"
	"github.com/FlutterDizaster/music-library/internal/domain/models"
	"github.com/google/uuid"
)

type Settings struct {
	SongsRepo   interfaces.SongsRepository
	LyricsRepo  interfaces.LyricsRepository
	LibraryRepo interfaces.LibraryRepository
	DetailsRepo interfaces.DetailsClient
}

type Service struct {
	songsRepo   interfaces.SongsRepository
	lyricsRepo  interfaces.LyricsRepository
	libraryRepo interfaces.LibraryRepository
	detailRepo  interfaces.DetailsClient
}

var _ interfaces.MusicService = (*Service)(nil)

// New creates a new Service instance with the provided settings.
// It initializes the music repository and details repository
// using the given interfaces from the settings parameter.
func New(settings Settings) *Service {
	return &Service{
		songsRepo:   settings.SongsRepo,
		lyricsRepo:  settings.LyricsRepo,
		libraryRepo: settings.LibraryRepo,
		detailRepo:  settings.DetailsRepo,
	}
}

// GetLibrary retrieves a filtered list of songs from the database.
//
// It logs debug messages before and after retrieving the library, and error messages if an error occurs.
//
// Parameters:
//   - ctx: The context for managing cancellation.
//   - filters: The filters to apply.
//
// Returns:
//   - The filtered library with pagination data if successful.
//   - An error if the operation fails or if no songs are found.
func (c *Service) GetLibrary(
	ctx context.Context,
	filters models.Filters,
) (models.Library, error) {
	return c.libraryRepo.GetLibrary(ctx, filters)
}

// GetSongLyrics retrieves the lyrics of a song from the database.
//
// It logs debug messages before and after retrieving the lyrics, and error messages if an error occurs.
//
// Parameters:
//   - ctx: The context for managing cancellation.
//   - id: The ID of the song.
//   - pagination: The pagination parameters.
//
// Returns:
//   - The filtered lyrics with pagination data if successful.
//   - An error if the operation fails or if the song is not found.
func (c *Service) GetSongLyrics(
	ctx context.Context,
	id uuid.UUID,
	pagination models.Pagination,
) (models.Lyrics, error) {
	rawLyrics, err := c.lyricsRepo.GetLyrics(ctx, id)
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

// AddSong adds a new song to the music library.
//
// It logs debug messages before and after retrieving the song details, and error messages if an error occurs.
//
// Parameters:
//   - ctx: The context for managing cancellation.
//   - title: The title and band of the song to add.
//
// Returns:
//   - The ID of the new song if successful.
//   - An error if the operation fails or if the song is not found.
func (c *Service) AddSong(ctx context.Context, title models.SongTitle) (uuid.UUID, error) {
	details, err := c.detailRepo.GetSongDetails(ctx, title)
	if err != nil {
		return uuid.Nil, err
	}

	song := models.Song{
		SongTitle:  title,
		SongDetail: details,
	}

	id, err := c.songsRepo.AddSong(ctx, song)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// DeleteSong deletes a song with the given ID from the music library.
//
// It logs debug messages before and after deleting, and error messages if an error occurs during the process.
//
// Parameters:
//   - ctx: The context for managing cancellation.
//   - id: The UUID of the song to be deleted.
//
// Returns:
//   - An error if the operation fails, otherwise nil.
func (c *Service) DeleteSong(ctx context.Context, id uuid.UUID) error {
	return c.songsRepo.DeleteSong(ctx, id)
}

// UpdateSong updates the details of an existing song in the library.
//
// It logs debug messages before and after updating, and error messages if an error occurs during the process.
//
// Parameters:
//   - ctx: The context for managing cancellation.
//   - song: The models.Song containing updated song details. Only non-empty fields
//     will be updated in the database.
//
// Returns:
//   - An error if the operation fails, or nil if the update is successful.
func (c *Service) UpdateSong(ctx context.Context, song models.Song) error {
	return c.songsRepo.UpdateSong(ctx, song)
}
