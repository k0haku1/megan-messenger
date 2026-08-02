export interface Message {
  id: string
  conversationId: string
  content: string
  createdAt: string
  sender: { id: string; username: string; avatarUrl?: string }
}
