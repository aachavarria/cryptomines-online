import { describe, it, expect, beforeEach, vi } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { useAuth } from './useAuth'

const mockGuestLogin = vi.fn()

vi.mock('../services/api', () => ({
  guestLogin: () => mockGuestLogin(),
}))

describe('useAuth', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('starts in loading state', () => {
    mockGuestLogin.mockReturnValue(new Promise(() => {})) // never resolves
    const { result } = renderHook(() => useAuth())
    expect(result.current.loading).toBe(true)
    expect(result.current.player).toBeNull()
    expect(result.current.error).toBeNull()
  })

  it('auto-calls guestLogin when no token exists', async () => {
    mockGuestLogin.mockResolvedValue({
      token: 'new-token',
      player: { id: 'p1', anonymous_id: 'guest_x', level: 1, created_at: '2026-01-01' },
    })

    const { result } = renderHook(() => useAuth())

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })

    expect(mockGuestLogin).toHaveBeenCalledTimes(1)
    expect(result.current.player).not.toBeNull()
    expect(result.current.player!.id).toBe('p1')
  })

  it('uses stored token if already logged in', async () => {
    localStorage.setItem('token', 'existing-token')
    localStorage.setItem('player_id', 'stored-player-id')

    const { result } = renderHook(() => useAuth())

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })

    expect(mockGuestLogin).not.toHaveBeenCalled()
    expect(result.current.player).not.toBeNull()
    expect(result.current.player!.id).toBe('stored-player-id')
  })

  it('sets error when guestLogin fails', async () => {
    mockGuestLogin.mockRejectedValue(new Error('Network error'))

    const { result } = renderHook(() => useAuth())

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })

    expect(result.current.error).toBe('Failed to connect to server')
    expect(result.current.player).toBeNull()
  })
})
