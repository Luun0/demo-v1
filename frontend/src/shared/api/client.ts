import createClient from 'openapi-fetch'
import type { paths } from './schema.gen'
import { ApiError, type Problem } from './errors'
import { getAccessToken, setAccessToken } from './token'

export const API_BASE = '/api/v1'

export function readCookie(name: string): string | null {
  const m = document.cookie.split('; ').find((c) => c.startsWith(`${name}=`))
  return m ? decodeURIComponent(m.slice(name.length + 1)) : null
}

let refreshing: Promise<boolean> | null = null

/** Single-flight refresh: concurrent 401s share one /auth/refresh call. */
export function refreshSession(): Promise<boolean> {
  refreshing ??= (async () => {
    try {
      const csrf = readCookie('ssp_csrf')
      if (!csrf) return false
      const res = await fetch(`${API_BASE}/auth/refresh`, {
        method: 'POST',
        credentials: 'same-origin',
        headers: { 'X-CSRF-Token': csrf },
      })
      if (!res.ok) {
        setAccessToken(null)
        return false
      }
      const body = (await res.json()) as { access_token: string }
      setAccessToken(body.access_token)
      return true
    } catch {
      return false
    } finally {
      queueMicrotask(() => {
        refreshing = null
      })
    }
  })()
  return refreshing
}

function withAuth(req: Request): Request {
  const token = getAccessToken()
  const r = new Request(req)
  if (token) r.headers.set('Authorization', `Bearer ${token}`)
  return r
}

/** fetch wrapper: bearer token + one retry after a successful refresh. */
export async function authFetch(input: Request): Promise<Response> {
  const retry = input.clone()
  const res = await fetch(withAuth(input))
  if (res.status !== 401 || /\/auth\/(login|register|refresh|logout)$/.test(new URL(input.url).pathname)) return res
  if (!(await refreshSession())) return res
  return fetch(withAuth(retry))
}

export const api = createClient<paths>({ baseUrl: `${window.location.origin}${API_BASE}`, credentials: 'same-origin', fetch: authFetch })

/** Unwraps an openapi-fetch result or throws ApiError. */
export function unwrap<T>(r: { data?: T; error?: unknown; response: Response }): T {
  if (!r.response.ok) throw new ApiError(r.response.status, r.error as Problem | undefined)
  return r.data as T
}
