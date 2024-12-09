package models

import "github.com/google/uuid"

//go:generate easyjson -all song.go
type Song struct {
	ID uuid.UUID `json:"id"` // ID is a unique identifier for the song
	SongTitle
	SongDetail
}
