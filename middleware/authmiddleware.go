package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/AungKyawPhyo1142/be-students-management-system/helpers"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			helpers.RespondWithErr(w, http.StatusUnauthorized, "Authorization header is missing!")
			return
		}

		// remove 'Bearer ' from token string
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected sigining method: %v", t.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			helpers.RespondWithErr(w, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err.Error()))
			return
		}

		// set the user's username to context so it can be used in the handlers

		// Define a custom type to avoid context key collisions
		type contextKey string
		const usernameKey contextKey = "username"
		ctx := context.WithValue(r.Context(), usernameKey, claims.Username)

		// pass it to next handler
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
