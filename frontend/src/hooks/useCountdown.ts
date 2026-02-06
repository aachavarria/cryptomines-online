import { useState, useEffect } from 'react'

export function useCountdown(finishAt: string | null): string {
  const [timeLeft, setTimeLeft] = useState(() => calcTimeLeft(finishAt))

  useEffect(() => {
    if (!finishAt) {
      setTimeLeft('')
      return
    }

    setTimeLeft(calcTimeLeft(finishAt))
    const interval = setInterval(() => {
      const left = calcTimeLeft(finishAt)
      setTimeLeft(left)
      if (!left) clearInterval(interval)
    }, 1000)

    return () => clearInterval(interval)
  }, [finishAt])

  return timeLeft
}

function calcTimeLeft(finishAt: string | null): string {
  if (!finishAt) return ''
  const diff = new Date(finishAt).getTime() - Date.now()
  if (diff <= 0) return ''

  const hours = Math.floor(diff / 3600000)
  const mins = Math.floor((diff % 3600000) / 60000)
  const secs = Math.floor((diff % 60000) / 1000)

  if (hours > 0) {
    return `${hours.toString().padStart(2, '0')}:${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

export function formatDuration(seconds: number): string {
  const hours = Math.floor(seconds / 3600)
  const mins = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60

  if (hours > 0) {
    return `${hours}h ${mins}m ${secs}s`
  }
  if (mins > 0) {
    return `${mins}m ${secs}s`
  }
  return `${secs}s`
}

export function formatNumber(n: number): string {
  if (n >= 1_000_000_000) return `${(n / 1_000_000_000).toFixed(2)}B`
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(2)}M`
  if (n >= 1000) return n.toLocaleString()
  return String(Math.floor(n))
}
