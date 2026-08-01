export type ConversationType = 'dm' | 'group'

export interface Conversation {
  id: string
  type: ConversationType
  title?: string
  slug?: string
  createdAt?: string
}

export interface Message {
  id: string
  conversationId: string
  content: string
  createdAt: string
  sender: { id: string; username: string; avatarUrl?: string }
}
