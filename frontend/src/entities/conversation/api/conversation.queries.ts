import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { conversationApi } from './conversation.api'
import { useSessionStore } from '@/entities/session/model/session.store'
import { queryKeys } from '@/shared/api/query-keys'
import type { Conversation } from '@/entities/conversation/model/types'

function sortConversations(conversations: Conversation[]): Conversation[] {
  return [...conversations].sort((left, right) => {
    const leftTime = left.createdAt ? new Date(left.createdAt).getTime() : 0
    const rightTime = right.createdAt ? new Date(right.createdAt).getTime() : 0
    return rightTime - leftTime
  })
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
