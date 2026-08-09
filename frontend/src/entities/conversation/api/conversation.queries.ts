import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { conversationApi } from './conversation.api'
import { useSessionStore } from '@/entities/session/model/session.store'
import { queryKeys } from '@/shared/api/query-keys'
import { getConversationActivityAt } from '@/entities/conversation/lib/display'
import type { Conversation } from '@/entities/conversation/model/types'

function sortConversations(conversations: Conversation[]): Conversation[] {
  return [...conversations].sort((left, right) => {
    const leftTime = activityTime(left)
    const rightTime = activityTime(right)
    return rightTime - leftTime
  })
}

function activityTime(conversation: Conversation): number {
  const iso = getConversationActivityAt(conversation)
  return iso ? new Date(iso).getTime() : 0
}

export function useConversations() {
  const session = useSessionStore()
  const enabled = computed(() => session.isAuthenticated && !session.needsUsername)

  const query = useQuery({
    queryKey: queryKeys.conversations,
    queryFn: conversationApi.list,
    enabled,
  })

  const conversations = computed(() => sortConversations(query.data.value?.conversations ?? []))

  return {
    query,
    conversations,
    isLoading: computed(() => query.isLoading.value),
    isError: computed(() => query.isError.value),
  }
}
