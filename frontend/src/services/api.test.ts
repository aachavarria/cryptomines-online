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
    if (typeof localStorage !== 'undefined' && typeof localStorage.clear === 'function') {
      localStorage.clear()
    }
  })

  it('creates a single axios instance pointed at /api', async () => {
    await import('./api')
    expect(axios.create).toHaveBeenCalledWith(
      expect.objectContaining({ baseURL: '/api' }),
    )
  })

  it('registers a request interceptor so the JWT can be attached', async () => {
    await import('./api')
    const inst = (axios.create as unknown as { mock: { results: Array<{ value: unknown }> } }).mock.results[0]
      .value as { interceptors: { request: { use: ReturnType<typeof vi.fn> } } }
    expect(inst.interceptors.request.use).toHaveBeenCalled()
  })

  it('healthCheck round-trips the /health response', async () => {
    const mockApi = axios.create({} as never) as unknown as {
      get: ReturnType<typeof vi.fn>
    }
    mockApi.get.mockResolvedValueOnce({ data: { status: 'ok' } })

    const { healthCheck } = await import('./api')
    expect(await healthCheck()).toEqual({ status: 'ok' })
  })

  it('guestLogin is deprecated and throws — auth is handled by Supabase now', async () => {
    const { guestLogin } = await import('./api')
    await expect(guestLogin()).rejects.toThrow(/Supabase/i)
  })
})
