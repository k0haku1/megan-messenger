import { request } from '@/shared/api/http'
import type { Conversation } from '@/entities/conversation/model/types'
import type { Project } from '@/entities/project/model/types'

export interface CreateGroupResponse {
  id: string
  slug: string
}

export interface JoinGroupResponse {
  id: string
}

export const conversationApi = {
  list: () => request<{ conversations: Conversation[] }>('/conversations'),

  createGroup: (payload: { title: string; memberIds: string[] }) =>
    request<CreateGroupResponse>('/conversations/group', {
      method: 'POST',
      body: JSON.stringify({
        title: payload.title,
        memberIds: payload.memberIds,
      }),
    }),

  joinBySlug: (slug: string) =>
    request<JoinGroupResponse>(`/conversations/join/${encodeURIComponent(slug)}`, {
      method: 'POST',
    }),

  leave: (conversationId: string) =>
    request<void>(`/conversations/${encodeURIComponent(conversationId)}/leave`, {
      method: 'POST',
    }),

  listLinkedProjects: (conversationId: string) =>
    request<{ projects: Pick<Project, 'id' | 'name'>[] }>(
      `/conversations/${encodeURIComponent(conversationId)}/projects`,
    ),

  markRead: (conversationId: string, messageId: string) =>
    request<{ lastReadAt: string; othersReadAt?: string }>(
      `/conversations/${encodeURIComponent(conversationId)}/read`,
      {
        method: 'POST',
        body: JSON.stringify({ messageId }),
      },
    ),

  readState: (conversationId: string) =>
    request<{ lastReadAt?: string; othersReadAt?: string }>(
      `/conversations/${encodeURIComponent(conversationId)}/read-state`,
    ),
}
