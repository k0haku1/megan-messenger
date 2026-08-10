import { request, uploadForm } from '@/shared/api/http'
import type { Conversation } from '@/entities/conversation/model/types'
import type { Message, MessageAttachment } from '@/entities/message/model/types'

type MessageListResponse = {
  messages: Message[]
  nextCursor: string
  othersReadAt?: string
}

type MediaListResponse = {
  items: MessageAttachment[]
  nextCursor: string
}

export const messageApi = {
  list: (conversationId: string, params?: { limit?: number; cursor?: string }) => {
    const search = new URLSearchParams()
    if (params?.limit) search.set('limit', String(params.limit))
    if (params?.cursor) search.set('cursor', params.cursor)
    const query = search.toString()

    return request<MessageListResponse>(
      `/conversations/${conversationId}/messages${query ? `?${query}` : ''}`,
    )
  },

  send: (
    conversationId: string,
    payload: { content?: string; replyToId?: string; attachmentIds?: string[] },
  ) =>
    request<{ message: Message }>(`/conversations/${conversationId}/messages`, {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  sendDM: (payload: { userId?: string; username?: string; content: string }) =>
    request<{ conversation: Conversation; message: Message }>('/conversations/dm/messages', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  uploadAttachment: (conversationId: string, file: File) => {
    const form = new FormData()
    form.append('file', file)
    return uploadForm<{ attachment: MessageAttachment }>(
      `/conversations/${conversationId}/attachments`,
      form,
    )
  },

  listMedia: (
    conversationId: string,
    params?: { limit?: number; cursor?: string; kind?: 'image' | 'video' | 'file' },
  ) => {
    const search = new URLSearchParams()
    if (params?.limit) search.set('limit', String(params.limit))
    if (params?.cursor) search.set('cursor', params.cursor)
    if (params?.kind) search.set('kind', params.kind)
    const query = search.toString()
    return request<MediaListResponse>(
      `/conversations/${conversationId}/media${query ? `?${query}` : ''}`,
    )
  },

  hideForMe: (conversationId: string, messageId: string) =>
    request<void>(`/conversations/${conversationId}/messages/${messageId}`, { method: 'DELETE' }),

  deleteForEveryone: (conversationId: string, messageId: string) =>
    request<{ message: Message }>(`/conversations/${conversationId}/messages/${messageId}/everyone`, {
      method: 'DELETE',
    }),

  forward: (conversationId: string, messageIds: string[], targetConversationId: string) =>
    request<{ messages: Message[] }>(`/conversations/${conversationId}/messages/forward`, {
      method: 'POST',
      body: JSON.stringify({ conversationId: targetConversationId, messageIds }),
    }),

  react: (conversationId: string, messageId: string, emoji: string, add: boolean) =>
    request<{ message: Message }>(`/conversations/${conversationId}/messages/${messageId}/reactions`, {
      method: add ? 'POST' : 'DELETE',
      body: JSON.stringify({ emoji }),
    }),
}
