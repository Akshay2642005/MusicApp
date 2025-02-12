package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"musicapp-server/middlewares"
	"musicapp-server/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.UserCreate
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

		newUser := models.User{
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
		var user models.User
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
		claims := &models.Claims{
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(models.JwtKey)
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
		userEmail, ok := r.Context().Value(middlewares.Key("userEmail")).(string)
		if !ok {
			http.Error(w, "Unable to fetch user information", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Welcome to the Home page, %s!", userEmail)
		/* http.Redirect(w, r, "home/songs", http.StatusSeeOther) */
	}
}
