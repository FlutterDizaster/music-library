package handler

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/FlutterDizaster/music-library/internal/apperrors"
	"github.com/FlutterDizaster/music-library/internal/domain/interfaces"
	"github.com/FlutterDizaster/music-library/internal/domain/models"
	"github.com/google/uuid"
)

// musicHandler handles music data operations.
// Must be created with newmusicHandler function.
type musicHandler struct {
	router  *http.ServeMux
	service interfaces.MusicService
}

func newMusicHandler(service interfaces.MusicService) *musicHandler {
	h := &musicHandler{
		service: service,
	}

	h.registerRoutes()

	return h
}

// ServeHTTP implements http.Handler interface.
func (h *musicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *musicHandler) registerRoutes() {
	router := http.NewServeMux()

	router.HandleFunc("GET /library", h.getLibraryHandler)
	router.HandleFunc("GET /song/{id}/lyrics", h.getSongLyricsHandler)
	router.HandleFunc("POST /song", h.addSongHandler)
	router.HandleFunc("PATCH /song/{id}", h.updateSongHandler)
	router.HandleFunc("DELETE /song/{id}", h.deleteSongHandler)

	h.router = router
}

//	@Summary		Get songs library
//	@Description	Retrieves the music library with optional filters and pagination.
//	@Tags			Songs
//	@Accept			json
//	@Produce		json
//	@Param			limit		query		int		false	"Maximum number of songs to return"
//	@Param			offset		query		int		false	"Number of songs to skip"
//	@Param			title		query		string	false	"Song title"
//	@Param			group		query		string	false	"Song band"
//	@Param			releaseDate	query		string	false	"Song release date (format: DD.MM.YYYY, valid values: >DD.MM.YYYY, <DD.MM.YYYY, DD.MM.YYYY, DD.MM.YYYY-DD.MM.YYYY)"
//	@Param			text		query		string	false	"Song lyrics fragment"
//	@Param			link		query		string	false	"Song link"
//	@Success		200			{object}	models.Library
//	@Failure		400			{string}	string	"Invalid query parameters"
//	@Failure		500			{string}	string	"Internal server error"
//	@Router			/library [get]
//
// .
//
//nolint:lll // too long comment
func (h *musicHandler) getLibraryHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	library, err := h.service.GetLibrary(r.Context(), params)
	var apperror *apperrors.Error
	switch {
	case errors.As(err, &apperror):
		http.Error(w, apperror.Message, http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	libraryJSON, err := library.MarshalJSON()
	if err != nil {
		slog.Error("Failed to marshal library", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(libraryJSON)
	if err != nil {
		slog.Error("Failed to write response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//	@Summary		Get song lyrics
//	@Description	Retrieves the lyrics for a specific song by ID.
//	@Tags			Songs
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Song ID"
//	@Param			limit	query		string	false	"Maximum number of verses to return"
//	@Param			offset	query		string	false	"Number of verses to skip"
//	@Success		200		{object}	models.Lyrics
//	@Failure		400		{string}	string	"Invalid ID or query parameters"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/song/{id}/lyrics [get]
//
// .
func (h *musicHandler) getSongLyricsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	params := r.URL.Query()

	lyrics, err := h.service.GetSongLyrics(r.Context(), id, params)
	var apperror *apperrors.Error
	switch {
	case errors.As(err, &apperror):
		http.Error(w, apperror.Message, http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lyricsJSON, err := lyrics.MarshalJSON()
	if err != nil {
		slog.Error("Failed to marshal lyrics", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(lyricsJSON)
	if err != nil {
		slog.Error("Failed to write response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//	@Summary		Add a new song
//	@Description	Adds a new song to the music library.
//	@Tags			Songs
//	@Accept			json
//	@Produce		text/plain
//	@Param			song	body		models.SongTitle	true	"Song title details"
//	@Success		201		{string}	string				"Song ID"
//	@Failure		400		{string}	string				"Invalid request body"
//	@Failure		500		{string}	string				"Internal server error"
//	@Router			/song [post]
//
// .
func (h *musicHandler) addSongHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		http.Error(w, "Wrong content type", http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var title models.SongTitle
	err = title.UnmarshalJSON(buf.Bytes())
	if err != nil {
		http.Error(w, "Invalid body content", http.StatusBadRequest)
		return
	}

	id, err := h.service.AddSong(r.Context(), title)
	var apperror *apperrors.Error
	switch {
	case errors.As(err, &apperror):
		http.Error(w, apperror.Message, http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(id.String()))
	if err != nil {
		slog.Error("Failed to write response", slog.Any("error", err))
	}
}

//	@Summary		Update a song
//	@Description	Updates details of an existing song.
//	@Tags			Songs
//	@Accept			json
//	@Produce		text/plain
//	@Param			id		path		string		true	"Song ID"
//	@Param			song	body		models.Song	true	"Updated song details"
//	@Success		200		{string}	string		"Update successful (no content returned)"
//	@Failure		400		{string}	string		"Invalid ID or request body"
//	@Failure		500		{string}	string		"Internal server error"
//	@Router			/song/{id} [patch]
//
// .
func (h *musicHandler) updateSongHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		http.Error(w, "Wrong content type", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var song models.Song
	err = song.UnmarshalJSON(buf.Bytes())
	if err != nil {
		http.Error(w, "Invalid body content", http.StatusBadRequest)
		return
	}
	song.ID = id

	err = h.service.UpdateSong(r.Context(), song)
	var apperror *apperrors.Error
	switch {
	case errors.As(err, &apperror):
		http.Error(w, apperror.Message, http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//	@Summary		Delete a song
//	@Description	Deletes a song by ID.
//	@Tags			Songs
//	@Produce		text/plain
//	@Param			id	path		string	true	"Song ID"
//	@Success		204	{string}	string	"Delete successful (no content returned)"
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/song/{id} [delete]
//
// .
func (h *musicHandler) deleteSongHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteSong(r.Context(), id)
	var apperror *apperrors.Error
	switch {
	case errors.As(err, &apperror):
		http.Error(w, apperror.Message, http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
