import { ApiError } from '@/shared/api/http'

type ErrorOptions = {
  fallback?: string
  forbidden?: string
}

export function getApiErrorMessage(error: unknown, options: ErrorOptions = {}): string {
  const fallback = options.fallback ?? 'Не удалось выполнить запрос'

  if (!(error instanceof ApiError)) {
    return fallback
  }

  if (error.status === 403 && options.forbidden) {
    return options.forbidden
  }

  if (error.fields) {
    const fieldMessage = Object.values(error.fields).find(Boolean)
    if (fieldMessage) {
      return fieldMessage
    }
  }

  return error.message || fallback
}
