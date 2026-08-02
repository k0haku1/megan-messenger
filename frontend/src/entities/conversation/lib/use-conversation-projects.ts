import { useQuery } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { conversationApi } from '@/entities/conversation/api/conversation.api'
import { queryKeys } from '@/shared/api/query-keys'

export function useConversationProjects(conversationId: Ref<string | null>, enabled: Ref<boolean>) {
  return useQuery({
    queryKey: computed(() =>
      conversationId.value
        ? queryKeys.conversationProjects(conversationId.value)
        : ['conversations', 'projects', 'none'],
    ),
    queryFn: async () => {
      if (!conversationId.value) return []
      return (await conversationApi.listLinkedProjects(conversationId.value)).projects
    },
    enabled: computed(() => enabled.value && !!conversationId.value),
  })
}
