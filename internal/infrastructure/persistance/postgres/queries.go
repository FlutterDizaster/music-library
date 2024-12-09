package postgres

const (
	libraryQuery = `SELECT COUNT(*) OVER() AS total_count, id, song, group, text, link, release_date FROM songs;`

	lyricsQuery = `SELECT text FROM songs WHERE id = $1 AND deleted = false;`

	addSongQuery = `INSERT INTO songs (song, group, text, link, release_date) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	updateFieldQuery = `UPDATE songs SET %s = $1 WHERE id = $2;`

	deleteSongQuery = `UPDATE songs SET deleted = true WHERE id = $1;`
)
