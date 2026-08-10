import type { Conversation } from '@/entities/conversation/model/types'
import { getConversationActivityAt } from '@/entities/conversation/lib/display'
import { queryClient } from '@/shared/api/query-client'
import { queryKeys } from '@/shared/api/query-keys'

export type ConversationActivity = {
  conversationId: string
  messageId: string
  content: string
  senderId: string
  senderUsername: string
  attachmentKind?: string
  createdAt: string
}

type ConversationsCache = {
  conversations: Conversation[]
}

function activityTime(conversation: Conversation): number {
  const iso = getConversationActivityAt(conversation)
  return iso ? new Date(iso).getTime() : 0
}

function sortConversations(conversations: Conversation[]): Conversation[] {
  return [...conversations].sort((left, right) => activityTime(right) - activityTime(left))
}

/**
 * Patches the conversations list cache from an inbox activity event.
 * Avoids a full list refetch on every new message.
 */
export function applyConversationActivity(
  activity: ConversationActivity,
  options: {
    currentUserId?: string
    /** When the viewer is in this chat and caught up, unread stays 0. */
    viewingConversationId?: string | null
    stickToBottom?: boolean
  } = {},
): void {
  const existing = queryClient.getQueryData<ConversationsCache>(queryKeys.conversations)
  if (!existing) {
    void queryClient.invalidateQueries({ queryKey: queryKeys.conversations })
    return
  }

  const index = existing.conversations.findIndex((item) => item.id === activity.conversationId)
  if (index < 0) {
    void queryClient.invalidateQueries({ queryKey: queryKeys.conversations })
    return
  }

  const previous = existing.conversations[index]!
  const isOwn = Boolean(options.currentUserId && activity.senderId === options.currentUserId)
  const viewingCaughtUp =
    options.viewingConversationId === activity.conversationId && Boolean(options.stickToBottom)

  let unreadCount = previous.unreadCount ?? 0
  if (isOwn || viewingCaughtUp) {
    // Own messages / reading at bottom do not increase unread.
    if (viewingCaughtUp) unreadCount = 0
  } else {
    unreadCount += 1
  }

  const next: Conversation = {
    ...previous,
    unreadCount,
    lastMessage: {
      id: activity.messageId,
      content: activity.content,
      senderId: activity.senderId,
      senderUsername: activity.senderUsername,
      attachmentKind: activity.attachmentKind,
      createdAt: activity.createdAt,
    },
  }

  const conversations = [...existing.conversations]
  conversations[index] = next

  queryClient.setQueryData<ConversationsCache>(queryKeys.conversations, {
    conversations: sortConversations(conversations),
  })
}

/** Clears unread for a conversation after mark-read without refetching the list. */
export function clearConversationUnread(conversationId: string): void {
  const existing = queryClient.getQueryData<ConversationsCache>(queryKeys.conversations)
  if (!existing) return

  const conversations = existing.conversations.map((item) =>
    item.id === conversationId ? { ...item, unreadCount: 0 } : item,
  )
  queryClient.setQueryData<ConversationsCache>(queryKeys.conversations, { conversations })
}
