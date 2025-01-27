package routes

import (
	"database/sql"
	"encoding/json"
	"musicapp-server/models"
	"net/http"
)

type Song = models.Song

func SongHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `SELECT id, title, artist, album, genre FROM songs`

		// Execute the query
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Failed to fetch songs from database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Slice to store the songs
		var songs []Song

		// Iterate through the rows and populate the songs slice
		for rows.Next() {
			var song Song
			if err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.Album, &song.Genre); err != nil {
				http.Error(w, "Error scanning song data", http.StatusInternalServerError)
				return
			}
			songs = append(songs, song)
		}

		// Check for errors encountered during iteration
		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating through songs", http.StatusInternalServerError)
			return
		}

		// Respond with the list of songs in JSON format
		if err := json.NewEncoder(w).Encode(songs); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
