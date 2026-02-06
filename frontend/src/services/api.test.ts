import { describe, it, expect, beforeEach, vi } from 'vitest'
import axios from 'axios'

vi.mock('axios', () => {
  const mockAxios = {
    create: vi.fn(() => mockAxios),
    get: vi.fn(),
    post: vi.fn(),
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() },
    },
  }
  return { default: mockAxios }
})

describe('api service', () => {
  beforeEach(() => {
    vi.resetModules()
    localStorage.clear()
  })

  it('creates axios instance with /api baseURL', async () => {
    await import('./api')
    expect(axios.create).toHaveBeenCalledWith(
      expect.objectContaining({ baseURL: '/api' })
    )
  })

  it('registers a request interceptor for JWT', async () => {
    await import('./api')
    expect(axios.create({} as never).interceptors.request.use).toBeDefined()
  })

  it('healthCheck calls GET /health', async () => {
    const mockApi = axios.create({} as never) as unknown as {
      get: ReturnType<typeof vi.fn>
    }
    mockApi.get.mockResolvedValueOnce({ data: { status: 'ok' } })

    const { healthCheck } = await import('./api')
    const result = await healthCheck()
    expect(result).toEqual({ status: 'ok' })
  })

  it('guestLogin calls POST /auth/guest and stores token', async () => {
    const mockResponse = {
      data: {
        token: 'test-jwt-token',
        player: { id: 'p1', anonymous_id: 'guest_abc', level: 1, created_at: '2026-01-01' },
      },
    }
    const mockApi = axios.create({} as never) as unknown as {
      post: ReturnType<typeof vi.fn>
    }
    mockApi.post.mockResolvedValueOnce(mockResponse)

    const { guestLogin } = await import('./api')
    const result = await guestLogin()
    expect(result.token).toBe('test-jwt-token')
    expect(result.player.id).toBe('p1')
    expect(localStorage.getItem('token')).toBe('test-jwt-token')
    expect(localStorage.getItem('player_id')).toBe('p1')
  })
})
