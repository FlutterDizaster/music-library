package models

//go:generate easyjson -all lyrics.go

// Lyrics contains filtered lyrics and pagination for song.
// @Description Response with lyrics and pagination.
type Lyrics struct {
	Lyrics     string     `json:"lyrics"`     // Lyrics of the song
	Pagination Pagination `json:"pagination"` // Pagination
}
