import { describe, expect, it, beforeEach } from 'vitest'
import { http, HttpResponse } from 'msw'
import { server } from '@/shared/lib/test/msw'
import { api, getAccessToken, refreshSession, setAccessToken } from '@/shared/api'

describe('api client', () => {
  beforeEach(() => {
    setAccessToken(null)
    document.cookie = 'ssp_csrf=csrf-123; path=/'
  })

  it('single-flight: concurrent refreshes issue exactly one /auth/refresh', async () => {
    let calls = 0
    server.use(
      http.post('*/api/v1/auth/refresh', ({ request }) => {
        calls++
        expect(request.headers.get('X-CSRF-Token')).toBe('csrf-123')
        return HttpResponse.json({ access_token: 'tok-1', token_type: 'Bearer', expires_in: 900 })
      }),
    )
    const results = await Promise.all([refreshSession(), refreshSession(), refreshSession()])
    expect(results).toEqual([true, true, true])
    expect(calls).toBe(1)
    expect(getAccessToken()).toBe('tok-1')
  })

  it('retries once after 401 with a refreshed token', async () => {
    setAccessToken('expired')
    const seen: (string | null)[] = []
    server.use(
      http.get('*/api/v1/auth/me', ({ request }) => {
        seen.push(request.headers.get('Authorization'))
        return request.headers.get('Authorization') === 'Bearer fresh'
          ? HttpResponse.json({ id: '1' })
          : new HttpResponse(null, { status: 401 })
      }),
      http.post('*/api/v1/auth/refresh', () => HttpResponse.json({ access_token: 'fresh', token_type: 'Bearer', expires_in: 900 })),
    )
    const r = await api.GET('/auth/me')
    expect(r.response.status).toBe(200)
    expect(seen).toEqual(['Bearer expired', 'Bearer fresh'])
  })

  it('never writes tokens to Web Storage', async () => {
    server.use(http.post('*/api/v1/auth/refresh', () => HttpResponse.json({ access_token: 'secret-token', token_type: 'Bearer', expires_in: 900 })))
    await refreshSession()
    const dump = JSON.stringify({ ...localStorage }) + JSON.stringify({ ...sessionStorage }) + document.cookie
    expect(dump).not.toContain('secret-token')
  })

  it('without csrf cookie no refresh is attempted', async () => {
    document.cookie = 'ssp_csrf=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/'
    expect(await refreshSession()).toBe(false)
  })
})
