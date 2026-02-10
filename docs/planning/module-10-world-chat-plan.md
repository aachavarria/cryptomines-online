# Module 10: World Chat - Implementation Plan

**Date:** 2026-02-07
**Module:** World Chat (Single World Channel)
**Estimated Hours:** 8-10 hours (2 days)
**Priority:** MEDIUM (Quality of Life)

---

## 1. Executive Summary

### Overview
World Chat provides a single global communication channel where all players can send and view messages in real-time. This is the simplest chat implementation - no private messages, no corps chat, no friends system. Players type a message, hit send, and it appears in the global feed for everyone to see. Recent 100 messages are displayed, and messages auto-refresh via polling.

### Scope
**IN SCOPE:**
- Single world channel (one global chat room)
- Send message (POST /api/chat/send)
- View recent messages (GET /api/chat/messages, limit 100)
- Real-time updates (5-second polling)
- Player name + message + timestamp display
- Optional: Basic profanity filter (simple word blacklist)
- Message length limit (500 characters)
- Rate limiting (1 message per 5 seconds per player)

**OUT OF SCOPE:**
- Corps chat (Corps system cut from scope)
- Private messages / DMs
- Friends system
- Mail system
- Loudspeaker item (chat is free, no premium messages)
- Message moderation tools (ban/mute players)
- Message history beyond 100 recent
- Message reactions/emojis
- Message editing/deletion
- User blocking
- Chat channels (only 1 world channel)

### Key Design Decisions
1. **Polling over WebSockets:** Use HTTP polling (every 5 seconds) instead of WebSockets/SSE for simplicity
2. **Recent 100 limit:** Only load last 100 messages (no pagination, no infinite scroll)
3. **Rate limiting:** 1 message per 5 seconds per player (prevent spam)
4. **No message history:** Messages older than recent 100 are not displayed (can be retained in DB for moderation)
5. **Simple profanity filter:** Optional word blacklist (replace with asterisks)
6. **Auto-scroll:** Chat auto-scrolls to bottom when new messages arrive

---

## 2. Current Database Analysis

### Existing Tables

**No existing chat tables** - This is a new system.

### Database Completeness: 0%
**Status:** ❌ NO CHAT TABLES EXIST - Full migration required

**Required Tables:**
1. `chat_messages` - Store all world chat messages

---

## 3. Galaxy Online 2 Research

### Sources
No specific Galaxy Online 2 chat documentation found in wiki archives. GO2 likely had a standard MMORPG chat system with world/corps/private channels. Since our scope is simplified (world channel only), we'll implement a standard chat pattern.

### Standard MMORPG Chat Mechanics

**Common Features:**
- World/Global channel: Public messages visible to all online players
- Message format: [Timestamp] PlayerName: Message
- Rate limiting: Prevent spam (e.g., 1 message per 5 seconds)
- Character limit: Typically 200-500 characters per message
- Profanity filter: Block/replace offensive words
- Auto-scroll: Chat window scrolls to bottom on new messages

### Cryptomines Online Implementation
- **Single world channel:** No tabs, no channel switching
- **Recent 100 messages:** Simple query: `ORDER BY created_at DESC LIMIT 100`
- **Polling:** Frontend polls GET /api/chat/messages every 5 seconds
- **Rate limiting:** Backend enforces 5-second cooldown between player messages
- **Profanity filter:** Optional - simple word blacklist (e.g., "bad_word" → "***")

---

## 4. System Architecture

### Component Overview
```
┌───────────────────────────────────────────────────────────────┐
│                       WORLD CHAT SYSTEM                        │
├───────────────────────────────────────────────────────────────┤
│                                                                │
│  ┌──────────────┐      ┌──────────────┐      ┌────────────┐ │
│  │   Frontend   │◄────►│   Backend    │◄────►│  Database  │ │
│  │              │ HTTP │              │ SQL  │            │ │
│  │  ChatPanel   │      │ Send/Get     │      │   chat_    │ │
│  │              │      │ handlers     │      │  messages  │ │
│  │  • Message   │      │              │      │            │ │
│  │    list      │      │ • Rate limit │      │ player_id  │ │
│  │  • Input     │      │ • Filter     │      │ message    │ │
│  │  • Auto-     │      │ • Trim       │      │ timestamp  │ │
│  │    scroll    │      │              │      │            │ │
│  │  • Polling   │      │              │      │            │ │
│  │    (5s)      │      │              │      │            │ │
│  └──────────────┘      └──────────────┘      └────────────┘ │
│                                                                │
└───────────────────────────────────────────────────────────────┘
```

### Data Flow: Send Message

```
1. Player types message in ChatPanel input
   ↓
2. Player clicks "Send" button
   ↓
3. Frontend validates:
   - Message not empty
   - Message length <= 500 characters
   ↓
4. POST /api/chat/send
   {
     "message": "Hello world!"
   }
   ↓
5. Backend validates:
   - Player authenticated
   - Rate limit check (last message >= 5 seconds ago)
   - Message length <= 500 characters
   - Apply profanity filter (optional)
   ↓
6. Insert into chat_messages:
   INSERT INTO chat_messages (player_id, username, message)
   VALUES (playerID, playerUsername, "Hello world!")
   ↓
7. Return success:
   {
     "success": true,
     "message_id": "uuid"
   }
   ↓
8. Frontend clears input, continues polling
```

### Data Flow: View Messages (Polling)

```
1. ChatPanel component mounts
   ↓
2. Start polling interval (every 5 seconds):
   setInterval(() => {
     fetchMessages()
   }, 5000)
   ↓
3. GET /api/chat/messages?since=<last_message_id>
   ↓
4. Backend queries recent 100 messages:
   SELECT id, player_id, username, message, created_at
   FROM chat_messages
   WHERE created_at > <since_timestamp>
   ORDER BY created_at DESC
   LIMIT 100
   ↓
5. Return messages:
   [
     {"id": "uuid1", "username": "Player1", "message": "Hello!", "created_at": "2026-02-07T15:30:00Z"},
     {"id": "uuid2", "username": "Player2", "message": "Hi!", "created_at": "2026-02-07T15:30:05Z"}
   ]
   ↓
6. Frontend appends new messages to chat list
   ↓
7. Auto-scroll to bottom if user was already at bottom
   ↓
8. Continue polling (loop back to step 3)
```

---

## 5. Database Schema Design

### 5.1 New Table: chat_messages

```sql
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    username TEXT NOT NULL,
    -- Denormalized username for faster reads (avoid JOIN on every message)
    message TEXT NOT NULL CHECK (LENGTH(message) > 0 AND LENGTH(message) <= 500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Index for fast recent message queries
    CONSTRAINT chk_message_length CHECK (LENGTH(message) <= 500)
);

CREATE INDEX idx_chat_messages_created_at ON chat_messages (created_at DESC);
CREATE INDEX idx_chat_messages_player ON chat_messages (player_id);
```

**Design Notes:**
- **Denormalized username:** Store username in chat_messages to avoid JOINing players table on every message query
- **500 character limit:** Enforced at DB level with CHECK constraint
- **Timestamp index:** Fast queries for recent 100 messages
- **No soft delete:** Messages are permanent (no edit/delete)
- **No read receipts:** No tracking of who read what

### 5.2 Optional: Rate Limiting Table

**Option 1: Use chat_messages table**
Check last message timestamp per player:
```sql
SELECT created_at
FROM chat_messages
WHERE player_id = $1
ORDER BY created_at DESC
LIMIT 1
```

**Option 2: Separate rate_limit table** (NOT RECOMMENDED - adds complexity)
```sql
CREATE TABLE chat_rate_limits (
    player_id UUID PRIMARY KEY REFERENCES players(id) ON DELETE CASCADE,
    last_message_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

**Recommendation:** Use Option 1 (query chat_messages) - simpler, no extra table.

---

## 6. Backend Implementation

### 6.1 New Handler: chat.go

**File:** `/backend/internal/handlers/chat.go`

```go
package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/cryptomines-online/backend/internal/database"
    "github.com/cryptomines-online/backend/internal/middleware"
)

// Profanity filter (simple word blacklist)
var profanityWords = []string{
    // Add offensive words here (example: "badword1", "badword2")
    // For production, use a comprehensive list
}

// SendMessage handles POST /api/chat/send
func SendMessage(w http.ResponseWriter, r *http.Request) {
    playerID := middleware.GetPlayerID(r)

    var req struct {
        Message string `json:"message"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
        return
    }

    // Validate message
    message := strings.TrimSpace(req.Message)
    if message == "" {
        http.Error(w, `{"error":"message cannot be empty"}`, http.StatusBadRequest)
        return
    }
    if len(message) > 500 {
        http.Error(w, `{"error":"message too long (max 500 characters)"}`, http.StatusBadRequest)
        return
    }

    // Rate limit check (5 seconds between messages)
    var lastMessageAt time.Time
    err := database.DB.QueryRow(`
        SELECT created_at
        FROM chat_messages
        WHERE player_id = $1
        ORDER BY created_at DESC
        LIMIT 1
    `, playerID).Scan(&lastMessageAt)
    if err != nil && err != sql.ErrNoRows {
        log.Printf("Failed to check rate limit: %v", err)
        http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
        return
    }

    if !lastMessageAt.IsZero() && time.Since(lastMessageAt) < 5*time.Second {
        http.Error(w, `{"error":"rate limit exceeded (wait 5 seconds)"}`, http.StatusTooManyRequests)
        return
    }

    // Apply profanity filter (optional)
    message = applyProfanityFilter(message)

    // Get player username
    var username string
    err = database.DB.QueryRow(`
        SELECT COALESCE(username, anonymous_id)
        FROM players
        WHERE id = $1
    `, playerID).Scan(&username)
    if err != nil {
        log.Printf("Failed to get player username: %v", err)
        http.Error(w, `{"error":"failed to get player"}`, http.StatusInternalServerError)
        return
    }

    // Insert message
    var messageID string
    err = database.DB.QueryRow(`
        INSERT INTO chat_messages (player_id, username, message)
        VALUES ($1, $2, $3)
        RETURNING id
    `, playerID, username, message).Scan(&messageID)
    if err != nil {
        log.Printf("Failed to insert message: %v", err)
        http.Error(w, `{"error":"failed to send message"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "success":    "true",
        "message_id": messageID,
    })
}

// GetMessages handles GET /api/chat/messages?since=<timestamp>
func GetMessages(w http.ResponseWriter, r *http.Request) {
    // Optional: since parameter for incremental updates
    sinceParam := r.URL.Query().Get("since")
    var sinceTime time.Time
    if sinceParam != "" {
        var err error
        sinceTime, err = time.Parse(time.RFC3339, sinceParam)
        if err != nil {
            http.Error(w, `{"error":"invalid since parameter"}`, http.StatusBadRequest)
            return
        }
    }

    // Query recent 100 messages (or messages since timestamp)
    query := `
        SELECT id, player_id, username, message, created_at
        FROM chat_messages
    `
    var args []interface{}
    if !sinceTime.IsZero() {
        query += ` WHERE created_at > $1`
        args = append(args, sinceTime)
    }
    query += ` ORDER BY created_at DESC LIMIT 100`

    rows, err := database.DB.Query(query, args...)
    if err != nil {
        log.Printf("Failed to query messages: %v", err)
        http.Error(w, `{"error":"failed to get messages"}`, http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    type message struct {
        ID        string    `json:"id"`
        PlayerID  string    `json:"player_id"`
        Username  string    `json:"username"`
        Message   string    `json:"message"`
        CreatedAt time.Time `json:"created_at"`
    }

    messages := []message{}
    for rows.Next() {
        var msg message
        if err := rows.Scan(&msg.ID, &msg.PlayerID, &msg.Username, &msg.Message, &msg.CreatedAt); err != nil {
            log.Printf("Failed to scan message: %v", err)
            continue
        }
        messages = append(messages, msg)
    }

    // Reverse order (oldest first for display)
    for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
        messages[i], messages[j] = messages[j], messages[i]
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(messages)
}

// applyProfanityFilter replaces profane words with asterisks
func applyProfanityFilter(message string) string {
    lower := strings.ToLower(message)
    for _, word := range profanityWords {
        // Simple case-insensitive replacement
        if strings.Contains(lower, strings.ToLower(word)) {
            replacement := strings.Repeat("*", len(word))
            message = strings.ReplaceAll(message, word, replacement)
            message = strings.ReplaceAll(message, strings.ToLower(word), replacement)
            message = strings.ReplaceAll(message, strings.ToUpper(word), replacement)
            message = strings.ReplaceAll(message, strings.Title(word), replacement)
        }
    }
    return message
}
```

### 6.2 Route Registration

**File:** `/backend/cmd/server/main.go` (update)

```go
// Chat endpoints
http.HandleFunc("POST /api/chat/send", middleware.AuthMiddleware(handlers.SendMessage))
http.HandleFunc("GET /api/chat/messages", middleware.AuthMiddleware(handlers.GetMessages))
```

### 6.3 Backend Tasks Breakdown

| Task | Description | Hours |
|------|-------------|-------|
| Create chat.go | SendMessage + GetMessages handlers | 2.0 |
| Rate limiting logic | Check last message timestamp, enforce 5-second cooldown | 0.5 |
| Profanity filter | Simple word blacklist (optional) | 0.5 |
| Database migration | Create chat_messages table + indexes | 0.5 |
| Route registration | Add 2 routes | 0.25 |
| Error handling | Comprehensive validation + error messages | 0.5 |
| Testing (manual) | Test send, rate limit, profanity filter, polling | 0.75 |
| **Subtotal** | | **5 hours** |

---

## 7. Frontend Implementation

### 7.1 New Component: ChatPanel.tsx

**File:** `/frontend/src/components/panels/ChatPanel.tsx`

```typescript
import { useState, useEffect, useRef } from 'react'
import { sendChatMessage, getChatMessages } from '../../api/chat'

interface ChatMessage {
  id: string
  player_id: string
  username: string
  message: string
  created_at: string
}

export default function ChatPanel() {
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const chatContainerRef = useRef<HTMLDivElement>(null)

  // Polling interval (5 seconds)
  useEffect(() => {
    fetchMessages() // Initial load

    const interval = setInterval(() => {
      fetchMessages()
    }, 5000) // Poll every 5 seconds

    return () => clearInterval(interval)
  }, [])

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    if (chatContainerRef.current) {
      const { scrollTop, scrollHeight, clientHeight } = chatContainerRef.current
      const isAtBottom = scrollHeight - scrollTop - clientHeight < 50
      if (isAtBottom) {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
      }
    }
  }, [messages])

  const fetchMessages = async () => {
    try {
      const lastMessage = messages[messages.length - 1]
      const since = lastMessage ? lastMessage.created_at : undefined
      const newMessages = await getChatMessages(since)

      if (newMessages.length > 0) {
        setMessages(prev => {
          const combined = [...prev, ...newMessages]
          // Keep only last 100 messages
          return combined.slice(-100)
        })
      }
    } catch (err) {
      console.error('Failed to fetch messages:', err)
    }
  }

  const handleSend = async () => {
    if (!input.trim()) return
    if (input.length > 500) {
      setError('Message too long (max 500 characters)')
      return
    }

    setLoading(true)
    setError(null)

    try {
      await sendChatMessage(input.trim())
      setInput('')
      // Fetch messages immediately after sending
      await fetchMessages()
    } catch (err: any) {
      if (err.message.includes('rate limit')) {
        setError('Please wait 5 seconds between messages')
      } else {
        setError(`Failed to send message: ${err.message}`)
      }
    } finally {
      setLoading(false)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  return (
    <div className="chat-panel">
      <h2>World Chat</h2>

      <div className="chat-messages" ref={chatContainerRef}>
        {messages.length === 0 && <p className="no-messages">No messages yet. Be the first to chat!</p>}
        {messages.map(msg => (
          <div key={msg.id} className="chat-message">
            <span className="chat-timestamp">{new Date(msg.created_at).toLocaleTimeString()}</span>
            <span className="chat-username">{msg.username}:</span>
            <span className="chat-text">{msg.message}</span>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      <div className="chat-input-container">
        {error && <div className="chat-error">{error}</div>}
        <input
          type="text"
          className="chat-input"
          placeholder="Type a message..."
          value={input}
          onChange={e => setInput(e.target.value)}
          onKeyPress={handleKeyPress}
          disabled={loading}
          maxLength={500}
        />
        <button
          className="chat-send-button"
          onClick={handleSend}
          disabled={loading || !input.trim()}
        >
          {loading ? 'Sending...' : 'Send'}
        </button>
        <div className="chat-char-count">{input.length}/500</div>
      </div>
    </div>
  )
}
```

### 7.2 API Functions

**File:** `/frontend/src/api/chat.ts`

```typescript
import { apiRequest } from './base'

export interface ChatMessage {
  id: string
  player_id: string
  username: string
  message: string
  created_at: string
}

export async function sendChatMessage(message: string): Promise<{ success: boolean; message_id: string }> {
  return apiRequest('/api/chat/send', {
    method: 'POST',
    body: JSON.stringify({ message }),
  })
}

export async function getChatMessages(since?: string): Promise<ChatMessage[]> {
  const url = since ? `/api/chat/messages?since=${encodeURIComponent(since)}` : '/api/chat/messages'
  return apiRequest<ChatMessage[]>(url)
}
```

### 7.3 ChatPanel Styling

**File:** `/frontend/src/styles/ChatPanel.css`

```css
.chat-panel {
  display: flex;
  flex-direction: column;
  height: 600px;
  border: 1px solid #ccc;
  border-radius: 8px;
  padding: 16px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 8px;
  margin-bottom: 16px;
  background-color: #f9f9f9;
}

.chat-message {
  margin-bottom: 8px;
  padding: 4px;
}

.chat-timestamp {
  color: #888;
  font-size: 0.85em;
  margin-right: 8px;
}

.chat-username {
  font-weight: bold;
  color: #2a5db0;
  margin-right: 4px;
}

.chat-text {
  color: #333;
}

.no-messages {
  text-align: center;
  color: #aaa;
  padding: 32px;
}

.chat-input-container {
  display: flex;
  gap: 8px;
  align-items: center;
}

.chat-input {
  flex: 1;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 14px;
}

.chat-send-button {
  padding: 8px 16px;
  background-color: #2a5db0;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.chat-send-button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.chat-char-count {
  font-size: 0.85em;
  color: #888;
}

.chat-error {
  color: red;
  font-size: 0.9em;
  margin-bottom: 8px;
}
```

### 7.4 Integration

**File:** `/frontend/src/App.tsx` (update)

```typescript
import ChatPanel from './components/panels/ChatPanel'

// Add chat button to main UI
<button onClick={() => setActivePanel('chat')}>Chat</button>

// Add to panel rendering
{activePanel === 'chat' && <ChatPanel />}
```

### 7.5 Frontend Tasks Breakdown

| Task | Description | Hours |
|------|-------------|-------|
| Create ChatPanel.tsx | Message list, input, send button, polling | 2.0 |
| Create api/chat.ts | API request functions (send, get) | 0.25 |
| Auto-scroll logic | Scroll to bottom on new messages (if user at bottom) | 0.5 |
| Polling implementation | setInterval 5-second polling | 0.25 |
| Styling (CSS) | Message list, input, timestamps, error display | 0.75 |
| Character counter | Show "X/500" character count | 0.25 |
| Error handling | Rate limit, network errors, validation | 0.5 |
| Testing | Manual UI testing for send, poll, scroll, rate limit | 0.5 |
| **Subtotal** | | **5 hours** |

---

## 8. QA Checklist

### 8.1 Backend Tests

**Manual Testing (Go):**

| # | Test Case | Expected Result | Status |
|---|-----------|----------------|--------|
| 1 | POST /api/chat/send with valid message | 200 OK, message inserted | ⬜ |
| 2 | POST /api/chat/send with empty message | 400 Bad Request: "message cannot be empty" | ⬜ |
| 3 | POST /api/chat/send with 501 characters | 400 Bad Request: "message too long" | ⬜ |
| 4 | POST /api/chat/send twice within 5 seconds | 429 Too Many Requests: "rate limit exceeded" | ⬜ |
| 5 | POST /api/chat/send after 5 seconds | 200 OK, second message allowed | ⬜ |
| 6 | POST /api/chat/send with profane word | Message filtered (*** replacement) | ⬜ |
| 7 | GET /api/chat/messages (empty chat) | 200 OK, empty array | ⬜ |
| 8 | GET /api/chat/messages (with messages) | 200 OK, recent 100 messages (oldest first) | ⬜ |
| 9 | GET /api/chat/messages?since=<timestamp> | Only messages after timestamp | ⬜ |
| 10 | Send 150 messages, GET /api/chat/messages | Only recent 100 returned | ⬜ |
| 11 | Username denormalized correctly | chat_messages.username matches players.username | ⬜ |
| 12 | Anonymous player sends message | username = anonymous_id | ⬜ |

### 8.2 Frontend Tests

**Manual Testing (UI):**

| # | Test Case | Expected Result | Status |
|---|-----------|----------------|--------|
| 13 | Open ChatPanel (empty chat) | "No messages yet" displayed | ⬜ |
| 14 | Open ChatPanel (with messages) | Recent 100 messages displayed | ⬜ |
| 15 | Type message, click Send | Message appears in chat, input clears | ⬜ |
| 16 | Type message, press Enter | Message sends (same as clicking Send) | ⬜ |
| 17 | Send message twice within 5 seconds | Error: "Please wait 5 seconds between messages" | ⬜ |
| 18 | Type 501 characters | Error: "Message too long" | ⬜ |
| 19 | Character counter updates | Shows "X/500" as user types | ⬜ |
| 20 | New message arrives (polling) | Chat updates automatically, auto-scrolls | ⬜ |
| 21 | Scroll up, new message arrives | Chat does NOT auto-scroll (user scrolled up) | ⬜ |
| 22 | At bottom, new message arrives | Chat auto-scrolls to show new message | ⬜ |
| 23 | Send button disabled while sending | Button shows "Sending..." and is disabled | ⬜ |
| 24 | Send button disabled for empty input | Cannot click Send with empty message | ⬜ |

### 8.3 Integration Tests

| # | Test Case | Expected Result | Status |
|---|-----------|----------------|--------|
| 25 | Two players chatting | Both see each other's messages | ⬜ |
| 26 | Player A sends, Player B sees within 5s | Message appears in Player B's chat via polling | ⬜ |
| 27 | 100+ messages, oldest scroll off | Only recent 100 visible | ⬜ |
| 28 | Profanity filter works end-to-end | Profane words replaced with *** | ⬜ |
| 29 | Rate limit enforced across page refreshes | Player cannot spam by refreshing | ⬜ |

---

## 9. Edge Cases & Error Handling

### 9.1 Edge Cases

| Case | Handling |
|------|----------|
| Empty message | 400 Bad Request: "message cannot be empty" |
| Message > 500 characters | 400 Bad Request: "message too long" |
| Rate limit (< 5 seconds) | 429 Too Many Requests: "rate limit exceeded" |
| Profane word in message | Replace with asterisks (*** ) |
| Anonymous player (no username) | Use anonymous_id as username |
| 100+ messages in chat | Only return recent 100 (ORDER BY created_at DESC LIMIT 100) |
| User scrolled up, new message | Do NOT auto-scroll (preserve scroll position) |
| User at bottom, new message | Auto-scroll to bottom |
| Network error during poll | Log error, continue polling (no user alert) |
| Network error during send | Show error message to user |

### 9.2 Validation Rules

**Backend:**
1. ✅ Message not empty (trimmed)
2. ✅ Message length <= 500 characters
3. ✅ Rate limit: 1 message per 5 seconds per player
4. ✅ Player must be authenticated (middleware)
5. ✅ Profanity filter applied (optional)

**Frontend:**
1. ✅ Message not empty (trimmed)
2. ✅ Message length <= 500 characters
3. ✅ Character counter shows "X/500"
4. ✅ Send button disabled while loading
5. ✅ Send button disabled for empty input
6. ✅ Error messages displayed for rate limit, network errors

---

## 10. Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| **Spam/flooding** | MEDIUM | MEDIUM | Rate limiting (5 seconds), profanity filter |
| **Offensive content** | MEDIUM | HIGH | Profanity filter (basic), future: report system |
| **Chat history loss** | LOW | LOW | Messages stored in DB (not in-memory) |
| **Polling overhead** | LOW | MEDIUM | 5-second interval, only load recent 100 |
| **Auto-scroll bug** | LOW | MEDIUM | Check scroll position before auto-scrolling |
| **Username desync** | LOW | LOW | Denormalized username in chat_messages |
| **Rate limit bypass (multiple tabs)** | LOW | LOW | Rate limit enforced server-side (per player_id) |

**Critical Paths:**
1. Rate limiting enforcement (prevent spam)
2. Polling efficiency (avoid loading full history on every poll)
3. Auto-scroll logic (only scroll if user at bottom)

---

## 11. Timeline & Estimates

### 11.1 Task Breakdown

| Phase | Task | Owner | Hours | Dependencies |
|-------|------|-------|-------|-------------|
| **Database** | | | | |
| 1 | Create migration (chat_messages table) | infra-dev | 0.5 | - |
| **Backend** | | | | |
| 2 | Create chat.go (SendMessage, GetMessages) | backend-dev | 2.0 | Task 1 |
| 3 | Rate limiting logic | backend-dev | 0.5 | Task 2 |
| 4 | Profanity filter (optional) | backend-dev | 0.5 | Task 2 |
| 5 | Route registration | backend-dev | 0.25 | Task 2 |
| 6 | Error handling | backend-dev | 0.5 | Task 2-4 |
| 7 | Manual testing (Postman) | backend-dev | 0.75 | Task 2-6 |
| **Frontend** | | | | |
| 8 | Create ChatPanel.tsx | frontend-dev | 2.0 | Backend Task 2 |
| 9 | Create api/chat.ts | frontend-dev | 0.25 | - |
| 10 | Auto-scroll logic | frontend-dev | 0.5 | Task 8 |
| 11 | Polling implementation | frontend-dev | 0.25 | Task 8 |
| 12 | Styling (CSS) | frontend-dev | 0.75 | Task 8 |
| 13 | Character counter | frontend-dev | 0.25 | Task 8 |
| 14 | Error handling | frontend-dev | 0.5 | Task 8 |
| 15 | Manual UI testing | frontend-dev | 0.5 | Task 8-14 |
| **QA** | | | | |
| 16 | Backend QA (12 tests) | qa-agent | 0.75 | Backend tasks |
| 17 | Frontend QA (12 tests) | qa-agent | 0.75 | Frontend tasks |
| 18 | Integration QA (5 tests) | qa-agent | 0.5 | All tasks |
| **TOTAL** | | | **12 hours** | |

### 11.2 Timeline (Parallel Execution)

**Day 1 (6 hours):**
- Infra: Task 1 (migration) - 0.5 hours
- Backend: Tasks 2-6 (chat.go, rate limit, filter, routes, errors) - 3.75 hours
- Frontend: Tasks 9, 11, 13 (API, polling, counter) - 1.0 hours

**Day 2 (6 hours):**
- Backend: Task 7 (testing) - 0.75 hours
- Frontend: Tasks 8, 10, 12, 14 (ChatPanel, auto-scroll, styling, errors) - 4.0 hours
- Frontend: Task 15 (UI testing) - 0.5 hours
- QA: Tasks 16-18 (all tests) - 2.0 hours (overlaps, can finish Day 2 or spill to Day 3)

**Total:** 12 hours (2 days with parallel work)

---

## 12. Success Criteria

### 12.1 Must-Have (MVP)
- ✅ Single world chat channel functional
- ✅ Players can send messages (POST /api/chat/send)
- ✅ Players can view recent 100 messages (GET /api/chat/messages)
- ✅ Polling updates chat every 5 seconds
- ✅ Rate limiting enforced (1 message per 5 seconds)
- ✅ Character limit enforced (500 characters)
- ✅ Auto-scroll to bottom on new messages (if user at bottom)
- ✅ Message format: [Timestamp] Username: Message

### 12.2 Nice-to-Have (Future)
- ⬜ Profanity filter (basic word blacklist)
- ⬜ Message timestamps formatted nicely (e.g., "5 minutes ago")
- ⬜ Player name clickable (view player profile)
- ⬜ Message reactions/emojis
- ⬜ User blocking (hide messages from specific players)
- ⬜ Chat channels (world, corps, trade, etc.)
- ⬜ WebSocket/SSE for true real-time updates (instead of polling)

### 12.3 Definition of Done
- All 29 QA tests pass
- Rate limiting prevents spam (verified)
- Auto-scroll works correctly (scrolls when at bottom, doesn't scroll when scrolled up)
- Messages display correctly (timestamp, username, message)
- Profanity filter working (if implemented)
- No polling performance issues (5-second interval acceptable)

---

## 13. Future Enhancements

### Phase 3+ Additions
1. **WebSocket/SSE:** Replace polling with WebSockets for true real-time updates
2. **Multiple Channels:** World, Trade, Corps (if Corps system added)
3. **Private Messages:** DM system between players
4. **Message Reactions:** Emoji reactions to messages
5. **User Blocking:** Hide messages from blocked players
6. **Message Moderation:** Admin tools to delete messages, ban users
7. **Message History:** Pagination to view older messages
8. **Loudspeaker Item:** Premium item for highlighted messages (if Mall system added)
9. **Link Previews:** Show previews for URLs in messages
10. **Message Formatting:** Bold, italic, colors

---

## 14. Appendix

### 14.1 Database Schema SQL

```sql
-- Migration file: /supabase/migrations/20260207010000_world_chat.sql

CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    username TEXT NOT NULL,
    message TEXT NOT NULL CHECK (LENGTH(message) > 0 AND LENGTH(message) <= 500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_message_length CHECK (LENGTH(message) <= 500)
);

CREATE INDEX idx_chat_messages_created_at ON chat_messages (created_at DESC);
CREATE INDEX idx_chat_messages_player ON chat_messages (player_id);

-- RLS policies
ALTER TABLE chat_messages ENABLE ROW LEVEL SECURITY;

-- All authenticated users can read all messages
CREATE POLICY chat_messages_select_all ON chat_messages
    FOR SELECT
    USING (true);

-- Players can only insert their own messages
CREATE POLICY chat_messages_insert_own ON chat_messages
    FOR INSERT
    WITH CHECK (auth.uid() = player_id);
```

### 14.2 Profanity Filter Word List (Example)

```go
var profanityWords = []string{
    // Add offensive words here
    // Example (placeholder - use a comprehensive list for production):
    "badword1",
    "badword2",
    "badword3",
    // For production, use a library like https://github.com/TwiN/go-away
}
```

### 14.3 Polling vs WebSocket Comparison

| Feature | Polling (Current) | WebSocket |
|---------|------------------|-----------|
| Complexity | Low | High |
| Real-time | ~5 second delay | Instant |
| Server load | Higher (constant polling) | Lower (push on change) |
| Scalability | Moderate | Better |
| Implementation time | 2 days | 4-5 days |
| Browser support | Universal | Modern browsers only |

**Recommendation:** Use polling for Phase 2 (simplicity), upgrade to WebSocket in Phase 3+ if needed.

---

## 15. Implementation Notes

### Backend Notes
- Use `ORDER BY created_at DESC LIMIT 100` to get recent messages efficiently
- Denormalize username in chat_messages to avoid JOIN on every query
- Rate limit check: query last message timestamp, enforce 5-second cooldown
- Profanity filter: simple word blacklist (replace with asterisks)
- No soft delete: messages are permanent (no edit/delete for MVP)

### Frontend Notes
- Polling interval: 5 seconds (balance between real-time and server load)
- Auto-scroll: only scroll if user was already at bottom (check scrollTop)
- Character counter: show "X/500" as user types
- Error handling: rate limit error shows friendly message ("Please wait 5 seconds")
- Message list: reverse order (oldest first for display)

### Testing Notes
- Test rate limiting: send 2 messages within 5 seconds (should fail)
- Test auto-scroll: scroll up, send message (should NOT auto-scroll)
- Test auto-scroll: at bottom, receive message (SHOULD auto-scroll)
- Test profanity filter: send message with banned word (should be filtered)
- Test polling: two players chatting (both see messages within 5 seconds)

---

**END OF PLAN**
