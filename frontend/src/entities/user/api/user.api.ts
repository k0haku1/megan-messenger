import { request } from '@/shared/api/http'
import type { PublicUser, UserPrivacySettings, UserSearchResult } from '../model/types'

export const userApi = {
  search: (query: string, limit = 20) =>
    request<{ users: UserSearchResult[] }>(
      `/users/search?q=${encodeURIComponent(query)}&limit=${limit}`,
    ),

  getByUsername: (username: string) =>
    request<PublicUser>(`/users/by-username/${encodeURIComponent(username)}`),

  updateUsername: (username: string) =>
    request<{ username: string }>('/users/me/username', {
      method: 'PATCH',
      body: JSON.stringify({ username }),
    }),

  updatePrivacy: (settings: UserPrivacySettings) =>
    request<UserPrivacySettings>('/users/me/privacy', {
      method: 'PATCH',
      body: JSON.stringify(settings),
    }),
}
