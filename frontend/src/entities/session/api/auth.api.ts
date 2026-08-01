import { request } from '@/shared/api/http'
import type { AuthResult, CompleteUsernameResponse, SessionUser } from '../model/types'

export const authApi = {
  startPhone: (phone: string) =>
    request<void>('/auth/phone/start', {
      method: 'POST',
      body: JSON.stringify({ phone }),
    }),

  verifyPhone: (phone: string, code: string) =>
    request<AuthResult>('/auth/phone/verify', {
      method: 'POST',
      body: JSON.stringify({ phone, code }),
    }),

  verifyPassword: (challengeToken: string, password: string) =>
    request<AuthResult>('/auth/password/verify', {
      method: 'POST',
      body: JSON.stringify({ challengeToken, password }),
    }),

  completeUsername: (username: string) =>
    request<CompleteUsernameResponse>('/auth/onboarding/username', {
      method: 'POST',
      body: JSON.stringify({ username }),
    }),

  getMe: () => request<SessionUser>('/users/me'),

  setPassword: (password: string, currentPassword?: string) =>
    request<void>('/auth/password', {
      method: 'POST',
      body: JSON.stringify(
        currentPassword
          ? { password, currentPassword }
          : { password },
      ),
    }),

  removePassword: (currentPassword: string) =>
    request<void>('/auth/password', {
      method: 'DELETE',
      body: JSON.stringify({ currentPassword }),
    }),

  logout: () => request<void>('/auth/logout', { method: 'POST' }),
}
