package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/FlutterDizaster/music-library/internal/application/apperrors"
	"github.com/FlutterDizaster/music-library/internal/domain/interfaces"
	"github.com/FlutterDizaster/music-library/internal/domain/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type updatingInfo struct {
	id    uuid.UUID
	field string
	value any
}

type Settings struct {
	// DatabaseDSN is the connection string for the PostgreSQL database.
	DatabaseDSN string

	// RetryCount is the number of times to retry the database connection.
	RetryCount int
	// RetryBackoff is the duration to wait before retrying the database connection.
	RetryBackoff time.Duration
}

// Repository implements the MusicRepository interface.
// It represents a PostgreSQL database connection pool and provides methods for interacting with the database.
// Must be created with New function.
type Repository struct {
	pool    *pgxpool.Pool
	config  *pgxpool.Config
	connStr string

	retryCount   int
	retryBackoff time.Duration
}

var _ interfaces.LibraryRepository = (*Repository)(nil)
var _ interfaces.LyricsRepository = (*Repository)(nil)
var _ interfaces.SongsRepository = (*Repository)(nil)

// New returns a new Repository instance based on the provided settings.
// It configures the PostgreSQL connection based on the connection string,
// sets the retry count and backoff time, and establishes a connection to the database.
// If the connection string is invalid or the database connection fails, it returns an error.
func New(ctx context.Context, settings Settings) (*Repository, error) {
	repo := &Repository{
		connStr: settings.DatabaseDSN,

		retryCount:   settings.RetryCount,
		retryBackoff: settings.RetryBackoff,
	}

	config, err := pgxpool.ParseConfig(repo.connStr)
	if err != nil {
		slog.Error("Error while parsing connection string", slog.Any("err", err))
		return nil, err
	}
	repo.config = config

	err = repo.connect(ctx)
	if err != nil {
		slog.Error("Error while connecting to database", slog.Any("err", err))
		return nil, err
	}

	return repo, nil
}

func (r *Repository) connect(ctx context.Context) error {
	slog.Debug("Connecting to PostgreSQL")
	var (
		pool *pgxpool.Pool
		err  error
	)

	for i := range r.retryCount {
		pool, err = pgxpool.NewWithConfig(ctx, r.config)
		if err == nil {
			if err = pool.Ping(ctx); err == nil {
				break
			}
		}

		slog.Info(
			"Failed to connect to PostgreSQL",
			slog.Int("attempt", i+1),
			slog.Int("max attempts", r.retryCount),
			slog.Duration("retry_backoff", r.retryBackoff),
			slog.Any("error", err),
		)

		waitCtx, cancle := context.WithTimeout(context.Background(), r.retryBackoff)
		defer cancle()

		select {
		case <-waitCtx.Done():
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if err != nil {
		slog.Error("Error while connecting to PostgreSQL")
		return err
	}

	if r.pool != nil {
		r.pool.Close()
	}

	r.pool = pool
	slog.Debug("Connected to PostgreSQL")

	return nil
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
func (r *Repository) GetLibrary(
	ctx context.Context,
	filters models.Filters,
) (models.Library, error) {
	filtersQuery, values := filters.ToQueryParams()
	slog.Debug(
		"Retrieving library",
		slog.String("filters", filtersQuery),
		slog.Any("values", values),
	)

	query := fmt.Sprintf("%s%s;", strings.TrimSuffix(libraryQuery, ";"), filtersQuery)

	rows, err := r.pool.Query(ctx, query, values...)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		slog.Debug(
			"No song found",
			slog.String("filters", filtersQuery),
			slog.Any("values", values),
		)
		return models.Library{}, apperrors.ErrNoData

	case err != nil:
		slog.Error(
			"Error while getting library",
			slog.String("filters", filtersQuery),
			slog.Any("values", values),
			slog.Any("error", err),
		)
		return models.Library{}, err
	}

	var (
		songs      = make([]models.Song, 0)
		totalCount int
	)

	for rows.Next() {
		var (
			song        models.Song
			releaseDate time.Time
		)

		err = rows.Scan(
			&totalCount,
			&song.ID,
			&song.Song,
			&song.Group,
			&song.Text,
			&song.Link,
			&releaseDate,
		)
		if err != nil {
			return models.Library{}, err
		}

		song.ReleaseDate = releaseDate.Format(models.DateLayout)

		songs = append(songs, song)
	}

	rows.Close()

	if err = rows.Err(); err != nil {
		return models.Library{}, err
	}

	library := models.Library{
		Songs: songs,
		Pagination: models.Pagination{
			Limit:  len(songs),
			Offset: filters.Offset(),
			Total:  totalCount,
		},
	}

	slog.Debug("Library retrieved", slog.Int("songs_count", len(library.Songs)))

	return library, nil
}

// GetSongLyrics retrieves the lyrics of a song with the given ID from the database.
//
// It logs debug messages before and after retrieving the lyrics, and error messages if an error occurs.
//
// Parameters:
//   - ctx: The context for managing cancellation.
//   - id: The UUID of the song whose lyrics are to be retrieved.
//
// Returns:
//   - The lyrics as a string if successful.
//   - An error if the operation fails or if no lyrics are found.
func (r *Repository) GetLyrics(ctx context.Context, id uuid.UUID) (string, error) {
	slog.Debug("Retrieving song lyrics", slog.String("song_id", id.String()))

	var lyrics string

	err := r.pool.QueryRow(ctx, lyricsQuery, id).Scan(&lyrics)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		slog.Debug("Song lyrics not found", slog.String("song_id", id.String()))
		return "", apperrors.ErrNoData

	case err != nil:
		slog.Error(
			"Error while getting song lyrics",
			slog.String("song_id", id.String()),
			slog.Any("error", err),
		)
		return "", apperrors.AppError{}
	}

	slog.Debug("Song lyrics retrieved", slog.String("song_id", id.String()))

	return lyrics, nil
}

// AddSong adds a song to the library.
//
// The function logs debug messages before and after adding, and error messages if an error occurs during the process.
//
// Parameters:
//   - ctx: The context for managing cancellation.
//   - song: The models.Song to be added.
//
// Returns:
//   - The ID of the added song if successful, otherwise nil.
//   - An error if the operation fails, otherwise nil.
func (r *Repository) AddSong(ctx context.Context, song models.Song) (uuid.UUID, error) {
	slog.Debug(
		"Adding song",
		slog.Group("song", slog.String("title", song.Song), slog.String("group", song.Group)),
	)

	releaseDate, err := time.Parse(models.DateLayout, song.ReleaseDate)
	if err != nil {
		slog.Debug(
			"Song not added due to incorrect release date",
			slog.String("title", song.Song),
			slog.String("group", song.Group),
		)
		return uuid.Nil, apperrors.ErrInvalidDateLayout
	}

	err = r.pool.QueryRow(ctx, addSongQuery, song.Song, song.Group, song.Text, song.Link, releaseDate).
		Scan(&song.ID)
	if err != nil {
		slog.Error(
			"Error while adding song",
			slog.Group(
				"song",
				slog.String("title", song.Song),
				slog.String("group", song.Group),
			),
			slog.Any("error", err),
		)
		return uuid.Nil, err
	}

	slog.Debug(
		"Song added",
		slog.Group(
			"song",
			slog.String("title", song.Song),
			slog.String("group", song.Group),
			slog.String("id", song.ID.String()),
		),
	)
	return song.ID, nil
}

// DeleteSong deletes a song with the given ID by marking it as deleted in the database.
//
// It logs debug messages before and after the deletion, and error messages if an error occurs during the process.
//
// Parameters:
//   - ctx: The context for managing cancellation.
//   - id: The UUID of the song to be deleted.
//
// Returns:
//   - An error if the operation fails, otherwise nil.
func (r *Repository) DeleteSong(ctx context.Context, id uuid.UUID) error {
	slog.Debug("Deleting song", slog.String("song_id", id.String()))

	_, err := r.pool.Exec(ctx, deleteSongQuery, id)
	if err != nil {
		slog.Error(
			"Error while deleting song",
			slog.String("song_id", id.String()),
			slog.Any("error", err),
		)
		return err
	}

	slog.Debug("Song deleted", slog.String("song_id", id.String()))
	return nil
}

// Important: The approach chosen was not the most optimal one in order
// to demonstrate the ability to work with transactions.
//
// UpdateSong updates the details of an existing song in the database.
// It updates only the fields provided in the song parameter that are not empty.
// The function starts a transaction for updating multiple fields and ensures
// atomicity by rolling back if any update fails.
//
// Parameters:
//   - ctx: The context for managing cancellation.
//   - song: The models.Song containing updated song details. Only non-empty fields
//     will be updated in the database.
//
// Returns:
//   - An error if the operation fails, or nil if the update is successful.
func (r *Repository) UpdateSong(ctx context.Context, song models.Song) error {
	slog.Debug("Updating song", slog.String("song_id", song.ID.String()))

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		slog.Error(
			"Song not updated",
			slog.String("reason", "Error while starting transaction"),
			slog.String("song_id", song.ID.String()),
			slog.Any("error", err),
		)
		return err
	}

	defer func(e error) {
		if e != nil {
			slog.Debug(
				"Song not updated",
				slog.String("reason", "Rolling back transaction"),
				slog.String("song_id", song.ID.String()),
			)
			//nolint:errcheck // ignore
			tx.Rollback(ctx)
		}
	}(err)

	if song.Song != "" {
		err = r.updateField(ctx, tx, updatingInfo{
			id:    song.ID,
			field: "title",
			value: song.Song,
		})
	}

	if err == nil && song.Group != "" {
		err = r.updateField(ctx, tx, updatingInfo{
			id:    song.ID,
			field: "band",
			value: song.Group,
		})
	}

	if err == nil && song.Text != "" {
		err = r.updateField(ctx, tx, updatingInfo{
			id:    song.ID,
			field: "text",
			value: song.Text,
		})
	}

	if err == nil && song.Link != "" {
		err = r.updateField(ctx, tx, updatingInfo{
			id:    song.ID,
			field: "link",
			value: song.Link,
		})
	}

	if err == nil && song.ReleaseDate != "" {
		var releaseDate time.Time
		releaseDate, err = time.Parse(models.DateLayout, song.ReleaseDate)
		if err != nil {
			return apperrors.ErrInvalidDateLayout
		}
		err = r.updateField(ctx, tx, updatingInfo{
			id:    song.ID,
			field: "release_date",
			value: releaseDate,
		})
	}

	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error(
			"Song not updated",
			slog.String("reason", "Error while committing transaction"),
			slog.String("song_id", song.ID.String()),
			slog.Any("error", err),
		)
		return err
	}

	slog.Debug("Song updated", slog.String("song_id", song.ID.String()))
	return nil
}

func (r *Repository) updateField(ctx context.Context, tx pgx.Tx, info updatingInfo) error {
	slog.Debug(
		"Updating song field",
		slog.String("song_id", info.id.String()),
		slog.String("field", info.field),
		slog.Any("value", info.value),
	)
	_, err := tx.Exec(ctx, fmt.Sprintf(updateFieldQuery, info.field), info.value, info.id)
	if err != nil {
		slog.Error(
			"Song not updated",
			slog.String("reason", "Error while updating field"),
			slog.String("song_id", info.id.String()),
			slog.String("field", info.field),
			slog.Any("value", info.value),
			slog.Any("error", err),
		)
		return err
	}
	return nil
}
