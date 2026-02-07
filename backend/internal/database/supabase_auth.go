package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// SupabaseAuthResponse represents the response from Supabase Auth signup.
type SupabaseAuthResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	User         struct {
		ID          string `json:"id"`
		IsAnonymous bool   `json:"is_anonymous"`
		CreatedAt   string `json:"created_at"`
	} `json:"user"`
}

// SignUpAnonymous creates an anonymous user via Supabase Auth REST API.
// POST {SUPABASE_URL}/auth/v1/signup with apikey header.
func SignUpAnonymous() (*SupabaseAuthResponse, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	if supabaseURL == "" {
		supabaseURL = "http://127.0.0.1:54321"
	}
	anonKey := os.Getenv("SUPABASE_ANON_KEY")
	if anonKey == "" {
		anonKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0"
	}

	url := supabaseURL + "/auth/v1/signup"

	// Anonymous sign-in: empty email/password triggers anonymous user creation
	// when enable_anonymous_sign_ins is true in config
	body, _ := json.Marshal(map[string]interface{}{})

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", anonKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call supabase auth: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("supabase auth error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var authResp SupabaseAuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, fmt.Errorf("failed to parse auth response: %w", err)
	}

	return &authResp, nil
}
