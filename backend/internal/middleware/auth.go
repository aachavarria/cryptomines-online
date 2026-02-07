package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const PlayerIDKey contextKey = "player_id"

// Auth is a middleware that validates JWT tokens from the Authorization header.
// Uses Supabase JWT secret for token validation.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		// Get Supabase JWT secret
		supabaseSecret := os.Getenv("SUPABASE_JWT_SECRET")
		if supabaseSecret == "" {
			// Default Supabase local dev JWT secret
			supabaseSecret = "super-secret-jwt-token-with-at-least-32-characters-long"
		}
		
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(supabaseSecret), nil
		})
		
		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		playerID, ok := claims["sub"].(string)
		if !ok || playerID == "" {
			http.Error(w, `{"error":"invalid token subject"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), PlayerIDKey, playerID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetPlayerID extracts the player ID from the request context.
func GetPlayerID(r *http.Request) string {
	if id, ok := r.Context().Value(PlayerIDKey).(string); ok {
		return id
	}
	return ""
}
