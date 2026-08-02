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

  send: (conversationId: string, content: string) =>
    request<{ message: Message }>(`/conversations/${conversationId}/messages`, {
      method: 'POST',
      body: JSON.stringify({ content }),
    }),

  sendDM: (payload: { userId?: string; username?: string; content: string }) =>
    request<{ conversation: Conversation; message: Message }>('/conversations/dm/messages', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),
}
