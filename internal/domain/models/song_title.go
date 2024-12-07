package models

//go:generate easyjson -all song_title.go

// SongTitle contains song band and title.
// @Description Song band and title.
type SongTitle struct {
	Group string `json:"group"` // Name of the band
	Song  string `json:"song"`  // Song title
}
