import type { Conversation } from '@/entities/conversation/model/types'

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
