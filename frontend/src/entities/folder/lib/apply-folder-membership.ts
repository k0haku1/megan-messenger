import type { Conversation } from '@/entities/conversation/model/types'
import { queryClient } from '@/shared/api/query-client'
import { queryKeys } from '@/shared/api/query-keys'

type ConversationsCache = {
  conversations: Conversation[]
}

export function patchConversationFolderMembership(
  conversationId: string,
  folderId: string,
  add: boolean,
): boolean {
  const existing = queryClient.getQueryData<ConversationsCache>(queryKeys.conversations)
  if (!existing) return false

  let changed = false
  const conversations = existing.conversations.map((conversation) => {
    if (conversation.id !== conversationId) return conversation

    const folderIds = new Set(conversation.folderIds ?? [])
    const had = folderIds.has(folderId)
    if (add && !had) {
      folderIds.add(folderId)
      changed = true
    } else if (!add && had) {
      folderIds.delete(folderId)
      changed = true
    }

    return changed ? { ...conversation, folderIds: [...folderIds] } : conversation
  })

  if (!changed) return true

  queryClient.setQueryData<ConversationsCache>(queryKeys.conversations, { conversations })
  return true
}
