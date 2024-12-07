package models

//go:generate easyjson -all song_detail.go

// SongDetail contains song details.
// @Description Additional song data.
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"` // Realease date of song format: DD-MM-YYYY
	Text        string `json:"text"`        // Lyrics
	Link        string `json:"link"`        // Link to song
}
