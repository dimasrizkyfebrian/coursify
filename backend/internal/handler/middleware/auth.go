package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Tipe custom untuk context key agar tidak bentrok
type contextKey string
const UserIDKey contextKey = "user_id"

// AuthMiddleware adalah middleware untuk memvalidasi JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 2. Pisahkan "Bearer " dari token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Authorization header must be in format 'Bearer {token}'", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// 3. Parse dan validasi token
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 4. Jika valid, ekstrak user_id dan simpan di context
		userID, ok := (*claims)["user_id"].(string)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Menambahkan user_id ke dalam context request
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// 5. Lanjutkan ke handler berikutnya dengan context yang sudah diperbarui
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}