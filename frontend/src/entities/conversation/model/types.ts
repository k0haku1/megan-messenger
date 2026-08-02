export type ConversationType = 'dm' | 'group'

export interface ConversationPeer {
  id: string
  username: string
  avatarUrl?: string
}

export interface Conversation {
  id: string
  type: ConversationType
  title?: string
  slug?: string
  peer?: ConversationPeer
  createdAt?: string
}
