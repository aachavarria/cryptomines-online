import { useState, useEffect, useCallback, useRef } from 'react'
import { getChatMessages, sendChatMessage } from '../services/api.ts'
import type { ChatMessage } from '../types'
import { useGameContext } from '../contexts/GameContext.tsx'

export function useChat(channel: 'world' | 'alliance' = 'world') {
  const { dispatch } = useGameContext()
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [loading, setLoading] = useState(false)
  const [sending, setSending] = useState(false)
  const [rateLimitCooldown, setRateLimitCooldown] = useState(0)
  const cooldownTimerRef = useRef<NodeJS.Timeout | null>(null)

  const fetchMessages = useCallback(
    async (before?: string) => {
      setLoading(true)
      try {
        const data = await getChatMessages(channel, before)
        if (before) {
          // Loading older messages (pagination)
          setMessages((prev) => [...data, ...prev])
        } else {
          // Initial load or refresh
          setMessages(data)
        }
      } catch (error) {
        console.error('Failed to fetch chat messages:', error)
      } finally {
        setLoading(false)
      }
    },
    [channel]
  )

  const sendMessage = useCallback(
    async (message: string) => {
      if (rateLimitCooldown > 0) {
        dispatch({
          type: 'SET_ERROR',
          payload: `Please wait ${rateLimitCooldown} seconds before sending another message`,
        })
        return false
      }

      setSending(true)
      try {
        await sendChatMessage(message, channel)
        // Refresh messages to show new message
        await fetchMessages()

        // Start 3-second cooldown
        setRateLimitCooldown(3)
        if (cooldownTimerRef.current) {
          clearInterval(cooldownTimerRef.current)
        }
        cooldownTimerRef.current = setInterval(() => {
          setRateLimitCooldown((prev) => {
            const next = prev - 1
            if (next <= 0) {
              if (cooldownTimerRef.current) {
                clearInterval(cooldownTimerRef.current)
                cooldownTimerRef.current = null
              }
              return 0
            }
            return next
          })
        }, 1000)

        return true
      } catch (error: unknown) {
        const axiosErr = error as { response?: { data?: { error?: string } } }
        const message = axiosErr.response?.data?.error || 'Failed to send message'
        dispatch({ type: 'SET_ERROR', payload: message })
        return false
      } finally {
        setSending(false)
      }
    },
    [channel, rateLimitCooldown, dispatch, fetchMessages]
  )

  const loadOlderMessages = useCallback(async () => {
    if (messages.length === 0 || loading) return
    const oldestMessage = messages[0]
    await fetchMessages(oldestMessage.created_at)
  }, [messages, loading, fetchMessages])

  // Auto-refresh every 5 seconds
  useEffect(() => {
    fetchMessages()
    const interval = setInterval(() => {
      fetchMessages()
    }, 5000)
    return () => clearInterval(interval)
  }, [fetchMessages])

  // Cleanup cooldown timer
  useEffect(() => {
    return () => {
      if (cooldownTimerRef.current) {
        clearInterval(cooldownTimerRef.current)
      }
    }
  }, [])

  return {
    messages,
    loading,
    sending,
    rateLimitCooldown,
    sendMessage,
    loadOlderMessages,
    refresh: fetchMessages,
  }
}
