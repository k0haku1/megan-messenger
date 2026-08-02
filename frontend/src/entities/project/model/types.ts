export type ProjectMemberRole = 'owner' | 'member'

export type DecisionStatus = 'proposed' | 'accepted' | 'deprecated'

export interface Project {
  id: string
  name: string
  description: string
  createdBy: string
  createdAt: string
  memberCount?: number
  docCount?: number
}

export interface ProjectMember {
  projectId: string
  userId: string
  role: ProjectMemberRole
  joinedAt: string
  username: string
  avatarUrl?: string
}

export interface ProjectDocSummary {
  id: string
  projectId: string
  slug: string
  title: string
  createdBy: string
  updatedAt: string
}

export interface ProjectDoc extends ProjectDocSummary {
  bodyMd: string
}

export interface ProjectDecision {
  id: string
  projectId: string
  summary: string
  context: string
  status: DecisionStatus
  messageId?: string
  createdBy: string
  createdAt: string
}

export interface ProjectConversationLink {
  projectId: string
  conversationId: string
  linkedAt: string
  type: 'dm' | 'group'
  title: string
  slug?: string
}

export type ProjectTab = 'docs' | 'decisions' | 'chats'
