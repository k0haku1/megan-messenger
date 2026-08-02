import { type Ref, watch } from 'vue'
import { messageApi } from '@/entities/message/api/message.api'
import { replaceMessagesForConversation } from '@/entities/message/api/message.repository'

export function useMessagesSync(conversationId: Ref<string | null>): void {
  watch(
    conversationId,
    (id, _, onCleanup) => {
      if (!id) return

      let cancelled = false
      onCleanup(() => {
        cancelled = true
      })

      void messageApi
        .list(id)
        .then(({ messages }) => {
          if (!cancelled) {
            return replaceMessagesForConversation(id, messages)
          }
        })
        .catch(() => {
          // Dexie keeps the last synced snapshot offline.
        })
    },
    { immediate: true },
  )
}
