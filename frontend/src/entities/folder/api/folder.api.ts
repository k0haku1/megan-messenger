import { request } from '@/shared/api/http'
import type { ChatFolder } from '../model/types'

export const folderApi = {
  list: () => request<{ folders: ChatFolder[] }>('/folders'),

  create: (payload: { name: string; icon?: string }) =>
    request<{ folder: ChatFolder }>('/folders', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  update: (folderId: string, payload: { name: string; icon?: string }) =>
    request<{ folder: ChatFolder }>(`/folders/${encodeURIComponent(folderId)}`, {
      method: 'PATCH',
      body: JSON.stringify(payload),
    }),

  delete: (folderId: string) =>
    request<void>(`/folders/${encodeURIComponent(folderId)}`, { method: 'DELETE' }),

  reorder: (folderIds: string[]) =>
    request<void>('/folders/reorder', {
      method: 'PUT',
      body: JSON.stringify({ folderIds }),
    }),

  addItems: (folderId: string, conversationIds: string[]) =>
    request<void>(`/folders/${encodeURIComponent(folderId)}/items`, {
      method: 'POST',
      body: JSON.stringify({ conversationIds }),
    }),

  removeItem: (folderId: string, conversationId: string) =>
    request<void>(
      `/folders/${encodeURIComponent(folderId)}/items/${encodeURIComponent(conversationId)}`,
      { method: 'DELETE' },
    ),
}
