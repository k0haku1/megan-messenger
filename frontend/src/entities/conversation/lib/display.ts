import type { Conversation, ConversationLastMessage } from '@/entities/conversation/model/types'

export function getConversationTitle(conversation: Conversation): string {
  if (conversation.type === 'group') {
    return conversation.title ?? 'Группа'
  }

  if (conversation.peer?.username) {
    return conversation.peer.username
  }

  return 'Личный диалог'
}

export function getConversationAvatarName(conversation: Conversation): string {
  if (conversation.type === 'group') {
    return conversation.title ?? 'G'
  }

  return conversation.peer?.username ?? 'U'
}

export function getConversationAvatarUrl(conversation: Conversation): string | undefined {
  return conversation.peer?.avatarUrl
}

export function getConversationPreview(conversation: Conversation): string {
  const last = conversation.lastMessage
  if (!last) return 'Нет сообщений'

  const body = formatLastMessageBody(last)
  if (conversation.type === 'group' && last.senderUsername) {
    return `${last.senderUsername}: ${body}`
  }
  return body
}

function formatLastMessageBody(last: ConversationLastMessage): string {
  const text = last.content?.trim()
  if (text) return text

  switch (last.attachmentKind) {
    case 'image':
      return 'Фото'
    case 'video':
      return 'Видео'
    case 'file':
      return 'Файл'
    default:
      return 'Сообщение'
  }
}

export function getConversationActivityAt(conversation: Conversation): string | undefined {
  return conversation.lastMessage?.createdAt ?? conversation.createdAt
}

export function formatConversationListTime(isoDate?: string): string {
  if (!isoDate) return ''

  const date = new Date(isoDate)
  if (Number.isNaN(date.getTime())) return ''

  const now = new Date()
  const sameDay =
    date.getFullYear() === now.getFullYear()
    && date.getMonth() === now.getMonth()
    && date.getDate() === now.getDate()

  if (sameDay) {
    return date.toLocaleTimeString('ru', { hour: '2-digit', minute: '2-digit' })
  }

  const yesterday = new Date(now)
  yesterday.setDate(now.getDate() - 1)
  const isYesterday =
    date.getFullYear() === yesterday.getFullYear()
    && date.getMonth() === yesterday.getMonth()
    && date.getDate() === yesterday.getDate()

  if (isYesterday) return 'вчера'

  if (date.getFullYear() === now.getFullYear()) {
    return date.toLocaleDateString('ru', { day: 'numeric', month: 'short' })
  }

  return date.toLocaleDateString('ru', { day: 'numeric', month: 'short', year: 'numeric' })
}
