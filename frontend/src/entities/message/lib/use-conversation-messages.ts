import { computed, type Ref } from 'vue'
import { db } from '@/shared/model/db'
import { useLiveQuery } from '@/shared/lib/use-live-query'
import type { Message } from '@/entities/message/model/types'

export function useConversationMessages(conversationId: Ref<string | null>) {
  const messages = useLiveQuery(() => db.messages.orderBy('createdAt').toArray(), [] as Message[])
  const hiddenMessages = useLiveQuery(() => db.hiddenMessages.toArray(), [] as { id: string }[])

  return computed(() => {
    const id = conversationId.value
    if (!id) {
      return []
    }

    const hiddenIds = new Set(hiddenMessages.value.map((item) => item.id))
    return messages.value.filter(
      (message) => message.conversationId === id && !message.deletedAt && !hiddenIds.has(message.id),
    )
  })
}
