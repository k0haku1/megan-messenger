import { getAccessToken, setAccessToken } from './access-token'

export class ApiError extends Error {
  constructor(
    readonly status: number,
    message: string,
    readonly code?: string,
    readonly fields?: Record<string, string>,
  ) {
    super(message)
  }
}

type ErrorPayload = {
  error?: {
    message?: string
    code?: string
    fields?: Record<string, string>
  }
}

type RefreshResponse = {
  accessToken?: string
}

let refreshPromise: Promise<boolean> | null = null

function shouldSkipRefresh(path: string): boolean {
  return (
    path === '/auth/refresh' ||
    path === '/auth/logout' ||
    path.startsWith('/auth/phone/') ||
    path === '/auth/password/verify'
  )
}

async function refreshAccessToken(): Promise<boolean> {
  if (refreshPromise) return refreshPromise

  refreshPromise = (async () => {
    try {
      const response = await fetch('/api/auth/refresh', {
        method: 'POST',
        credentials: 'include',
      })
      if (!response.ok) return false

      const data = (await response.json().catch(() => null)) as RefreshResponse | null
      if (!data?.accessToken) return false

      setAccessToken(data.accessToken)
      return true
    } catch {
      return false
    } finally {
      refreshPromise = null
    }
  })()

  return refreshPromise
}

async function parseResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const payload = (await response.json().catch((): ErrorPayload => ({}))) as ErrorPayload
    throw new ApiError(
      response.status,
      payload.error?.message ?? response.statusText,
      payload.error?.code,
      payload.error?.fields,
    )
  }

  if (response.status === 204) {
    return undefined as T
  }

  const text = await response.text()
  if (!text) {
    return undefined as T
  }

  return JSON.parse(text) as T
}

function authHeaders(extra?: HeadersInit): Headers {
  const headers = new Headers(extra)
  const token = getAccessToken()
  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`)
  }
  return headers
}

async function fetchOnce(path: string, init?: RequestInit): Promise<Response> {
  const headers = authHeaders(init?.headers)
  if (!headers.has('Content-Type') && init?.body && !(init.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json')
  }

  return fetch(`/api${path}`, {
    ...init,
    headers,
    credentials: 'include',
  })
}

async function fetchWithRefresh(path: string, init?: RequestInit): Promise<Response> {
  const response = await fetchOnce(path, init)
  if (response.status !== 401 || shouldSkipRefresh(path)) {
    return response
  }

  const refreshed = await refreshAccessToken()
  if (!refreshed) {
    setAccessToken(null)
    return response
  }

  return fetchOnce(path, init)
}

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetchWithRefresh(path, init)
  return parseResponse<T>(response)
}

export async function uploadForm<T>(path: string, form: FormData): Promise<T> {
  const response = await fetchWithRefresh(path, {
    method: 'POST',
    body: form,
  })
  return parseResponse<T>(response)
}
