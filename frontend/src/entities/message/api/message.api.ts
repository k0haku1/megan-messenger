import { request } from '@/shared/api/http'
import type { Conversation } from '@/entities/conversation/model/types'
import type { Message } from '@/entities/message/model/types'

type MessageListResponse = {
  messages: Message[]
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

  send: (conversationId: string, content: string, replyToId?: string) =>
    request<{ message: Message }>(`/conversations/${conversationId}/messages`, {
      method: 'POST',
      body: JSON.stringify({ content, replyToId }),
    }),

  sendDM: (payload: { userId?: string; username?: string; content: string }) =>
    request<{ conversation: Conversation; message: Message }>('/conversations/dm/messages', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  hideForMe: (conversationId: string, messageId: string) =>
    request<void>(`/conversations/${conversationId}/messages/${messageId}`, { method: 'DELETE' }),

  deleteForEveryone: (conversationId: string, messageId: string) =>
    request<{ message: Message }>(`/conversations/${conversationId}/messages/${messageId}/everyone`, { method: 'DELETE' }),

  forward: (conversationId: string, messageId: string, targetConversationId: string) =>
    request<{ message: Message }>(`/conversations/${conversationId}/messages/${messageId}/forward`, {
      method: 'POST', body: JSON.stringify({ conversationId: targetConversationId }),
    }),

  react: (conversationId: string, messageId: string, emoji: string, add: boolean) =>
    request<{ message: Message }>(`/conversations/${conversationId}/messages/${messageId}/reactions`, {
      method: add ? 'POST' : 'DELETE', body: JSON.stringify({ emoji }),
    }),
}
