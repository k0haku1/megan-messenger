export class ApiError extends Error {
  constructor(readonly status: number, message: string, readonly code?: string) { super(message) }
}

type ErrorPayload = { error?: { message?: string; code?: string } }

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`/api${path}`, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...init?.headers },
    credentials: 'include',
  })
  if (!response.ok) {
    const payload = await response.json().catch((): ErrorPayload => ({})) as ErrorPayload
    throw new ApiError(response.status, payload.error?.message ?? response.statusText, payload.error?.code)
  }
  return response.json() as Promise<T>
}
