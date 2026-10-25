// Access token lives ONLY in memory (never localStorage/sessionStorage, ADR-0004).
let accessToken: string | null = null
const listeners = new Set<() => void>()

export function getAccessToken(): string | null {
  return accessToken
}

export function setAccessToken(token: string | null): void {
  accessToken = token
  listeners.forEach((l) => l())
}

export function subscribeToken(listener: () => void): () => void {
  listeners.add(listener)
  return () => listeners.delete(listener)
}
