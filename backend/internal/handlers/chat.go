package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/utils"
)

type sendMessageRequest struct {
	Message string `json:"message"`
	Channel string `json:"channel"`
}

type chatMessage struct {
	ID         string    `json:"id"`
	PlayerID   string    `json:"player_id"`
	PlayerName string    `json:"player_name"`
	Message    string    `json:"message"`
	Channel    string    `json:"channel"`
	CreatedAt  time.Time `json:"created_at"`
}

// SendChatMessage handles POST /api/chat/send
func SendChatMessage(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate message
	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		http.Error(w, `{"error":"message cannot be empty"}`, http.StatusBadRequest)
		return
	}

	if len(req.Message) > 500 {
		http.Error(w, `{"error":"message too long (max 500 characters)"}`, http.StatusBadRequest)
		return
	}

	// Validate channel
	if req.Channel != "world" && req.Channel != "alliance" {
		http.Error(w, `{"error":"invalid channel"}`, http.StatusBadRequest)
		return
	}

	// Check rate limit (3 seconds between messages)
	var lastMessageAt sql.NullTime
	err := database.DB.QueryRow(`
		SELECT last_message_at FROM chat_rate_limits WHERE player_id = $1
	`, playerID).Scan(&lastMessageAt)

	if err == nil && lastMessageAt.Valid {
		if time.Since(lastMessageAt.Time) < 3*time.Second {
			remaining := 3*time.Second - time.Since(lastMessageAt.Time)
			http.Error(w, `{"error":"rate limited","remaining_seconds":`+string(rune(int(remaining.Seconds())))+`}`, http.StatusTooManyRequests)
			return
		}
	}

	// Profanity filter
	if utils.ContainsProfanity(req.Message) {
		http.Error(w, `{"error":"message contains inappropriate content"}`, http.StatusBadRequest)
		return
	}

	// Insert message
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var messageID string
	err = tx.QueryRow(`
		INSERT INTO chat_messages (player_id, message, channel)
		VALUES ($1, $2, $3)
		RETURNING id
	`, playerID, req.Message, req.Channel).Scan(&messageID)

	if err != nil {
		log.Printf("Failed to insert chat message: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Update rate limit
	_, err = tx.Exec(`
		INSERT INTO chat_rate_limits (player_id, last_message_at, message_count)
		VALUES ($1, now(), 1)
		ON CONFLICT (player_id) DO UPDATE
		SET last_message_at = now(),
		    message_count = chat_rate_limits.message_count + 1,
		    updated_at = now()
	`, playerID)

	if err != nil {
		log.Printf("Failed to update rate limit: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"message_id": messageID,
	})
}

// GetChatMessages handles GET /api/chat/messages?channel=world&limit=50&before=timestamp
func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	channel := r.URL.Query().Get("channel")
	if channel == "" {
		channel = "world"
	}

	limit := 50
	before := r.URL.Query().Get("before")

	// Build query
	query := `
		SELECT cm.id, cm.player_id, p.username, cm.message, cm.channel, cm.created_at
		FROM chat_messages cm
		JOIN players p ON p.id = cm.player_id
		WHERE cm.channel = $1
	`

	args := []interface{}{channel}

	if before != "" {
		query += ` AND cm.created_at < $2`
		args = append(args, before)
	}

	query += fmt.Sprintf(` ORDER BY cm.created_at DESC LIMIT $%d`, len(args)+1)
	args = append(args, limit)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Failed to get chat messages: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	messages := []chatMessage{}
	for rows.Next() {
		var msg chatMessage
		err := rows.Scan(&msg.ID, &msg.PlayerID, &msg.PlayerName, &msg.Message, &msg.Channel, &msg.CreatedAt)
		if err != nil {
			log.Printf("Failed to scan chat message: %v", err)
			continue
		}
		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i := 0; i < len(messages)/2; i++ {
		messages[i], messages[len(messages)-1-i] = messages[len(messages)-1-i], messages[i]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
