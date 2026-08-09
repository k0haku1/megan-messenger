export type ConversationType = 'dm' | 'group'

export interface ConversationPeer {
  id: string
  username: string
  avatarUrl?: string
}

export interface ConversationLastMessage {
  id: string
  content: string
  senderUsername?: string
  attachmentKind?: 'image' | 'video' | 'file' | string
  createdAt: string
}

export interface Conversation {
  id: string
  type: ConversationType
  title?: string
  slug?: string
  peer?: ConversationPeer
  lastMessage?: ConversationLastMessage
  createdAt?: string
}
