package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateToken_ReturnsValidJWT(t *testing.T) {
	token, err := GenerateToken("test-player-123")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken returned empty string")
	}
	// JWT has 3 parts separated by dots
	parts := 0
	for _, c := range token {
		if c == '.' {
			parts++
		}
	}
	if parts != 2 {
		t.Errorf("JWT should have 3 parts (2 dots): got %d dots", parts)
	}
}

func TestGenerateToken_DifferentPlayersGetDifferentTokens(t *testing.T) {
	t1, _ := GenerateToken("player-1")
	t2, _ := GenerateToken("player-2")
	if t1 == t2 {
		t.Error("Different players should get different tokens")
	}
}

func TestAuth_MissingHeader(t *testing.T) {
	handler := Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when auth is missing")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Missing auth header: got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_InvalidFormat(t *testing.T) {
	handler := Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid format")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Authorization", "Basic sometoken")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Invalid auth format: got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	handler := Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid token")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Invalid token: got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_ValidToken(t *testing.T) {
	playerID := "test-player-abc"
	token, err := GenerateToken(playerID)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	var gotPlayerID string
	handler := Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPlayerID = GetPlayerID(r)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Valid token: got %d, want %d", w.Code, http.StatusOK)
	}
	if gotPlayerID != playerID {
		t.Errorf("Player ID: got %q, want %q", gotPlayerID, playerID)
	}
}

func TestGetPlayerID_NoContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	id := GetPlayerID(req)
	if id != "" {
		t.Errorf("GetPlayerID without context: got %q, want empty", id)
	}
}

func TestGetPlayerID_WithContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), PlayerIDKey, "my-player-id")
	req = req.WithContext(ctx)

	id := GetPlayerID(req)
	if id != "my-player-id" {
		t.Errorf("GetPlayerID with context: got %q, want %q", id, "my-player-id")
	}
}
