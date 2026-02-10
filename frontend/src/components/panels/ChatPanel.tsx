import { useState, useRef, useEffect } from 'react'
import { useChat } from '../../hooks/useChat.ts'
import LoadingButton from '../common/LoadingButton.tsx'
import '../../styles/chat.css'
import '../../styles/common.css'

export default function ChatPanel() {
  const [channel, setChannel] = useState<'world' | 'alliance'>('world')
  const [inputMessage, setInputMessage] = useState('')
  const { messages, loading, sending, rateLimitCooldown, sendMessage, loadOlderMessages } =
    useChat(channel)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const messagesContainerRef = useRef<HTMLDivElement>(null)
  const [isAtBottom, setIsAtBottom] = useState(true)

  // Auto-scroll to bottom when new messages arrive (if user is at bottom)
  useEffect(() => {
    if (isAtBottom && messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' })
    }
  }, [messages, isAtBottom])

  // Track if user is at bottom of message list
  const handleScroll = () => {
    if (!messagesContainerRef.current) return
    const { scrollTop, scrollHeight, clientHeight } = messagesContainerRef.current
    const threshold = 50 // pixels from bottom
    setIsAtBottom(scrollHeight - scrollTop - clientHeight < threshold)
  }

  const handleSend = async () => {
    const trimmedMessage = inputMessage.trim()
    if (!trimmedMessage || sending || rateLimitCooldown > 0) return

    const success = await sendMessage(trimmedMessage)
    if (success) {
      setInputMessage('')
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffMins = Math.floor(diffMs / 60000)

    if (diffMins < 1) return 'just now'
    if (diffMins < 60) return `${diffMins}m ago`
    if (diffMins < 1440) return `${Math.floor(diffMins / 60)}h ago`
    return date.toLocaleDateString()
  }

  return (
    <div className="panel chat-panel">
      <div className="panel-header">
        <h2>World Chat</h2>
        <div className="chat-channel-tabs">
          <button
            className={`chat-tab ${channel === 'world' ? 'active' : ''}`}
            onClick={() => setChannel('world')}
          >
            World
          </button>
          <button
            className={`chat-tab ${channel === 'alliance' ? 'active' : ''}`}
            onClick={() => setChannel('alliance')}
          >
            Alliance
          </button>
        </div>
      </div>

      <div className="panel-content chat-content">
        {/* Load Older Button */}
        <div className="chat-load-more">
          <LoadingButton
            className="btn btn-small btn-secondary"
            onClick={loadOlderMessages}
            loading={loading}
            disabled={messages.length === 0}
          >
            Load Older Messages
          </LoadingButton>
        </div>

        {/* Messages List */}
        <div
          className="chat-messages"
          ref={messagesContainerRef}
          onScroll={handleScroll}
        >
          {messages.length === 0 && !loading && (
            <div className="chat-empty">
              <p>No messages yet. Be the first to chat!</p>
            </div>
          )}

          {messages.map((msg) => (
            <div key={msg.id} className="chat-message">
              <div className="chat-message-header">
                <span className="chat-player-name">{msg.player_name || 'Anonymous'}</span>
                <span className="chat-timestamp">{formatTime(msg.created_at)}</span>
              </div>
              <div className="chat-message-body">{msg.message}</div>
            </div>
          ))}

          <div ref={messagesEndRef} />
        </div>

        {/* Input Area */}
        <div className="chat-input-area">
          <div className="chat-input-wrapper">
            <textarea
              className="chat-input"
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Type your message... (Max 500 characters)"
              maxLength={500}
              rows={2}
              disabled={sending || rateLimitCooldown > 0}
            />
            <div className="chat-input-footer">
              <span className="chat-char-count">
                {inputMessage.length}/500
              </span>
              {rateLimitCooldown > 0 && (
                <span className="chat-cooldown">
                  Cooldown: {rateLimitCooldown}s
                </span>
              )}
            </div>
          </div>
          <LoadingButton
            className="btn btn-primary chat-send-btn"
            onClick={handleSend}
            loading={sending}
            disabled={
              rateLimitCooldown > 0 ||
              !inputMessage.trim() ||
              inputMessage.length > 500
            }
          >
            Send
          </LoadingButton>
        </div>
      </div>
    </div>
  )
}
