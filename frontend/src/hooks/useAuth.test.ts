import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'

// Stub the Supabase client before importing the hook so the mock is in place
// for the module-level `supabase.auth` calls inside useEffect.
const getSession = vi.fn()
const signInAnonymously = vi.fn()
const onAuthStateChange = vi.fn().mockReturnValue({
  data: { subscription: { unsubscribe: vi.fn() } },
})

vi.mock('../lib/supabase', () => ({
  supabase: {
    auth: {
      getSession: () => getSession(),
      signInAnonymously: () => signInAnonymously(),
      onAuthStateChange: (cb: unknown) => onAuthStateChange(cb),
    },
  },
}))

import { useAuth } from './useAuth'

describe('useAuth', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    onAuthStateChange.mockReturnValue({
      data: { subscription: { unsubscribe: vi.fn() } },
    })
  })

  it('uses an existing session if Supabase already has one', async () => {
    getSession.mockResolvedValue({
      data: {
        session: {
          user: { id: 'u-existing', created_at: '2026-04-28T00:00:00Z' },
        },
      },
    })
    const { result } = renderHook(() => useAuth())
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(result.current.player?.id).toBe('u-existing')
    expect(signInAnonymously).not.toHaveBeenCalled()
  })

  it('falls back to anonymous sign-in when no session is present', async () => {
    getSession.mockResolvedValue({ data: { session: null } })
    signInAnonymously.mockResolvedValue({
      data: { user: { id: 'u-anon', created_at: '2026-04-28T00:00:00Z' } },
      error: null,
    })
    const { result } = renderHook(() => useAuth())
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(signInAnonymously).toHaveBeenCalledTimes(1)
    expect(result.current.player?.id).toBe('u-anon')
  })

  it('sets an error when anonymous sign-in fails', async () => {
    getSession.mockResolvedValue({ data: { session: null } })
    signInAnonymously.mockResolvedValue({
      data: { user: null },
      error: new Error('boom'),
    })
    const { result } = renderHook(() => useAuth())
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(result.current.error).toBe('Failed to connect to server')
    expect(result.current.player).toBeNull()
  })
})
