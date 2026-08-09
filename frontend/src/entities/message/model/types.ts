export interface Message {
  id: string
  conversationId: string
  content: string
  replyToId?: string
  forwardedFromId?: string
  deletedAt?: string
  reactions: MessageReaction[]
  reactionUpdatedBy?: string
  createdAt: string
  sender: { id: string; username: string; avatarUrl?: string }
}

export interface MessageReaction {
  emoji: string
  count: number
  reactedByMe: boolean
  reactedBy: Message['sender'][]
}
