import { type Ref, computed, unref } from 'vue'
import { useConversations } from '@/entities/conversation/api/conversation.queries'

export function useConversation(conversationId: Ref<string | null>) {
  const { conversations } = useConversations()

  return computed(() => {
    const id = unref(conversationId)
    if (!id) {
      return undefined
    }

    return conversations.value.find((conversation) => conversation.id === id)
  })
}
