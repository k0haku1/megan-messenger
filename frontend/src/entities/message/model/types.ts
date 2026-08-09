export interface MessageAttachment {
  id: string
  conversationId: string
  uploaderId?: string
  mime: string
  kind: 'image' | 'video' | 'file'
  sizeBytes: number
  width?: number
  height?: number
  durationMs?: number
  originalName: string
  url?: string
  thumbUrl?: string
  createdAt: string
}

export interface Message {
  id: string
  conversationId: string
  content: string
  replyToId?: string
  forwardedFromId?: string
  deletedAt?: string
  reactions: MessageReaction[]
  reactionUpdatedBy?: string
  attachments?: MessageAttachment[]
  createdAt: string
  sender: { id: string; username: string; avatarUrl?: string }
}

export interface MessageReaction {
  emoji: string
  count: number
  reactedByMe: boolean
  reactedBy: Message['sender'][]
}
