package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCreate struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type Song struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Genre  string `json:"genre"`
}

func SignupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UserCreate
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		row := db.QueryRow("SELECT 1 FROM users WHERE email = $1", user.Email)
		var exists int
		err = row.Scan(&exists)
		if err == nil && exists == 1 {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userID := uuid.New().String()

		_, err = db.Exec("INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)", userID, user.Name, user.Email, string(hashedPassword))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newUser := User{
			ID:       userID,
			Name:     user.Name,
			Email:    user.Email,
			Password: string(hashedPassword),
		}
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var storedPasswordHash string

		row := db.QueryRow("SELECT password FROM users WHERE email = $1", user.Email)
		err = row.Scan(&storedPasswordHash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(user.Password))
		if err != nil {
			http.Error(w, "Invalid Password", http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(time.Hour * 24)
		claims := &Claims{
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(JwtKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{
			"token": tokenString,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func HomeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query to fetch songs
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
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(songs); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
