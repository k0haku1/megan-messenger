import { request } from '@/shared/api/http'
import type {
  Project,
  ProjectConversationLink,
  ProjectDecision,
  ProjectDoc,
  ProjectDocSummary,
  ProjectMember,
} from '@/entities/project/model/types'

export const projectApi = {
  list: () => request<{ projects: Project[] }>('/projects'),

  get: (projectId: string) => request<{ id: string; name: string; description: string; createdBy: string; createdAt: string }>(`/projects/${projectId}`),

  create: (payload: { name: string; description?: string }) =>
    request<{ id: string; name: string; description: string; createdBy: string; createdAt: string }>('/projects', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  listMembers: (projectId: string) =>
    request<{ members: ProjectMember[] }>(`/projects/${projectId}/members`),

  addMember: (projectId: string, username: string) =>
    request<void>(`/projects/${projectId}/members`, {
      method: 'POST',
      body: JSON.stringify({ username }),
    }),

  listDocs: (projectId: string) =>
    request<{ docs: ProjectDocSummary[] }>(`/projects/${projectId}/docs`),

  getDoc: (projectId: string, docId: string) =>
    request<{ doc: ProjectDoc }>(`/projects/${projectId}/docs/${docId}`),

  createDoc: (projectId: string, payload: { title: string; bodyMd?: string }) =>
    request<{ doc: ProjectDoc }>(`/projects/${projectId}/docs`, {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  updateDoc: (projectId: string, docId: string, payload: { title: string; bodyMd: string }) =>
    request<{ doc: ProjectDoc }>(`/projects/${projectId}/docs/${docId}`, {
      method: 'PATCH',
      body: JSON.stringify(payload),
    }),

  deleteDoc: (projectId: string, docId: string) =>
    request<void>(`/projects/${projectId}/docs/${docId}`, { method: 'DELETE' }),

  listDecisions: (projectId: string) =>
    request<{ decisions: ProjectDecision[] }>(`/projects/${projectId}/decisions`),

  createDecision: (
    projectId: string,
    payload: { summary: string; context?: string; messageId?: string },
  ) =>
    request<{ decision: ProjectDecision }>(`/projects/${projectId}/decisions`, {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  updateDecisionStatus: (projectId: string, decisionId: string, status: string) =>
    request<{ decision: ProjectDecision }>(`/projects/${projectId}/decisions/${decisionId}`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    }),

  listConversations: (projectId: string) =>
    request<{ conversations: ProjectConversationLink[] }>(`/projects/${projectId}/conversations`),

  linkConversation: (projectId: string, conversationId: string) =>
    request<void>(`/projects/${projectId}/conversations`, {
      method: 'POST',
      body: JSON.stringify({ conversationId }),
    }),

  unlinkConversation: (projectId: string, conversationId: string) =>
    request<void>(`/projects/${projectId}/conversations/${conversationId}`, {
      method: 'DELETE',
    }),
}
