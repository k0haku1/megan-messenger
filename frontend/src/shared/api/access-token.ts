const ACCESS_TOKEN_KEY = 'mm_access_token'

type Listener = (token: string | null) => void

const listeners = new Set<Listener>()

export function getAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY)
}

export function setAccessToken(token: string | null): void {
  if (token) localStorage.setItem(ACCESS_TOKEN_KEY, token)
  else localStorage.removeItem(ACCESS_TOKEN_KEY)
  for (const listener of listeners) listener(token)
}

export function subscribeAccessToken(listener: Listener): () => void {
  listeners.add(listener)
  return () => {
    listeners.delete(listener)
  }
}
