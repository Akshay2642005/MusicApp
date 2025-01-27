package models

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Key string

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
