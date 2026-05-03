import { useState, useRef, useEffect } from 'react'
import { createPortal } from 'react-dom'
import { MessageCircle, Send, X } from 'lucide-react'
import { useChat } from '../../hooks/useChat.ts'
import LoadingButton from '../common/LoadingButton.tsx'

interface ChatPanelProps {
  onClose: () => void
}

export default function ChatPanel({ onClose }: ChatPanelProps) {
  const [channel, setChannel] = useState<'world' | 'alliance'>('world')
  const [inputMessage, setInputMessage] = useState('')
  const { messages, loading, sending, rateLimitCooldown, sendMessage, loadOlderMessages } =
    useChat(channel)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const messagesContainerRef = useRef<HTMLDivElement>(null)
  const [isAtBottom, setIsAtBottom] = useState(true)

  // ESC to close
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

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

  return createPortal(
    <div
      className="ds-modal-backdrop"
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose()
      }}
    >
      <div
        className="ds-modal"
        style={{
          width: 'min(800px, 95vw)',
          height: 'min(600px, 85vh)',
          maxHeight: '85vh',
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <div className="ds-modal-header">
          <div className="ds-row">
            <MessageCircle size={20} aria-hidden style={{ color: 'var(--ds-teal)' }} />
            <h2 className="ds-modal-title">{channel === 'world' ? 'World Chat' : 'Alliance Chat'}</h2>
          </div>
          <button className="ds-btn-icon" onClick={onClose} aria-label="Close chat">
            <X size={16} />
          </button>
        </div>

        <div className="ds-tabs" style={{ padding: '0 var(--sp-6)' }}>
          <button
            className="ds-tab"
            aria-selected={channel === 'world'}
            onClick={() => setChannel('world')}
          >
            World
          </button>
          <button
            className="ds-tab"
            aria-selected={channel === 'alliance'}
            onClick={() => setChannel('alliance')}
          >
            Alliance
          </button>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', flex: 1, minHeight: 0 }}>
          {/* Load Older Button */}
          <div
            style={{
              padding: 'var(--sp-3) var(--sp-4)',
              borderBottom: '1px solid var(--ds-border)',
              display: 'flex',
              justifyContent: 'center',
            }}
          >
            <LoadingButton
              className="ds-btn-ghost ds-btn--sm"
              onClick={loadOlderMessages}
              loading={loading}
              disabled={messages.length === 0}
            >
              Load Older Messages
            </LoadingButton>
          </div>

          {/* Messages List */}
          <div
            ref={messagesContainerRef}
            onScroll={handleScroll}
            style={{
              flex: 1,
              overflowY: 'auto',
              padding: 'var(--sp-4)',
              display: 'flex',
              flexDirection: 'column',
              gap: 'var(--sp-3)',
            }}
          >
            {messages.length === 0 && !loading && (
              <div
                style={{
                  flex: 1,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: 'var(--ds-text-muted)',
                  fontStyle: 'italic',
                }}
              >
                <p>No messages yet. Be the first to chat!</p>
              </div>
            )}

            {messages.map((msg) => (
              <div key={msg.id} className="ds-card">
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    marginBottom: 'var(--sp-1)',
                  }}
                >
                  <span
                    style={{
                      fontWeight: 'var(--fw-semibold)',
                      color: 'var(--ds-teal-dark)',
                      fontSize: 'var(--fs-sm)',
                    }}
                  >
                    {msg.player_name || 'Anonymous'}
                  </span>
                  <span className="ds-text-muted" style={{ fontSize: 'var(--fs-caption)' }}>
                    {formatTime(msg.created_at)}
                  </span>
                </div>
                <div
                  style={{
                    color: 'var(--ds-text)',
                    lineHeight: 'var(--lh-normal)',
                    wordWrap: 'break-word',
                    whiteSpace: 'pre-wrap',
                    fontSize: 'var(--fs-body)',
                  }}
                >
                  {msg.message}
                </div>
              </div>
            ))}

            <div ref={messagesEndRef} />
          </div>

          {/* Input Area */}
          <div
            style={{
              display: 'flex',
              gap: 'var(--sp-3)',
              padding: 'var(--sp-4)',
              borderTop: '1px solid var(--ds-border)',
              background: 'var(--ds-surface-2)',
            }}
          >
            <div
              style={{
                flex: 1,
                display: 'flex',
                flexDirection: 'column',
                gap: 'var(--sp-1)',
              }}
            >
              <textarea
                className="ds-input"
                value={inputMessage}
                onChange={(e) => setInputMessage(e.target.value)}
                onKeyDown={handleKeyPress}
                placeholder="Type your message... (Max 500 characters)"
                maxLength={500}
                rows={2}
                disabled={sending || rateLimitCooldown > 0}
                style={{ minHeight: 60, resize: 'none', fontFamily: 'var(--ff-sans)' }}
              />
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  fontSize: 'var(--fs-caption)',
                  color: 'var(--ds-text-muted)',
                }}
              >
                <span className="ds-text-soft">{inputMessage.length}/500</span>
                {rateLimitCooldown > 0 && (
                  <span style={{ color: 'var(--ds-danger)', fontWeight: 'var(--fw-semibold)' }}>
                    Cooldown: {rateLimitCooldown}s
                  </span>
                )}
              </div>
            </div>
            <div style={{ alignSelf: 'flex-end', minWidth: 100, height: 'fit-content' }}>
              <LoadingButton
                className="ds-btn-secondary"
                onClick={handleSend}
                loading={sending}
                disabled={
                  rateLimitCooldown > 0 ||
                  !inputMessage.trim() ||
                  inputMessage.length > 500
                }
              >
                <Send size={14} aria-hidden />
                Send
              </LoadingButton>
            </div>
          </div>
        </div>
      </div>
    </div>,
    document.body,
  )
}
