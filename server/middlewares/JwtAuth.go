package middlewares

import (
	"context"
	"musicapp-server/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Key = models.Key

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]

		// Verify the token
		claims := &models.Claims{}
		secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		} else if claims.ExpiresAt < time.Now().Unix() {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler if token is valid
		// Add the email to the context
		ctx := context.WithValue(r.Context(), Key("userEmail"), claims.Email)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
