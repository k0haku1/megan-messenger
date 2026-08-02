import type { Message } from '@/entities/message/model/types'

export function isOwnMessage(message: Message, currentUserId?: string): boolean {
  return Boolean(currentUserId) && message.sender.id === currentUserId
}

export function shouldShowSenderName(
  message: Message,
  index: number,
  messages: Message[],
  currentUserId?: string,
): boolean {
  if (isOwnMessage(message, currentUserId)) {
    return false
  }

  if (index === 0) {
    return true
  }

  return messages[index - 1]?.sender.id !== message.sender.id
}

export function formatMessageTime(isoDate: string): string {
  return new Date(isoDate).toLocaleTimeString('ru', {
    hour: '2-digit',
    minute: '2-digit',
  })
}
