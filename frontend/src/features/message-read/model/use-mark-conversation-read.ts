import { watch, type Ref } from 'vue'
import { conversationApi } from '@/entities/conversation/api/conversation.api'
import { clearConversationUnread } from '@/entities/conversation/lib/apply-conversation-activity'
import { useReadWatermarkStore } from '@/features/message-read/model/read-watermark.store'

const MARK_READ_DEBOUNCE_MS = 250

/**
 * Advances the viewer's read watermark when the chat is active and the user
 * is near the latest messages (Telegram-like open/scroll behavior).
 */
export function useMarkConversationRead(options: {
  conversationId: Ref<string | null>
  latestMessageId: Ref<string | null>
  enabled: Ref<boolean>
  stickToBottom: Ref<boolean>
}): void {
  const watermarks = useReadWatermarkStore()
  let timer: ReturnType<typeof setTimeout> | null = null
  let lastMarkedId: string | null = null
  let lastConversationId: string | null = null

  function clearTimer(): void {
    if (!timer) return
    clearTimeout(timer)
    timer = null
  }

  async function mark(conversationId: string, messageId: string): Promise<void> {
    if (lastMarkedId === messageId) return
    try {
      const result = await conversationApi.markRead(conversationId, messageId)
      lastMarkedId = messageId
      if (result.othersReadAt !== undefined) {
        watermarks.setOthersReadAt(conversationId, result.othersReadAt)
      }
      clearConversationUnread(conversationId)
    } catch {
      // Keep local UI; next open/scroll will retry.
    }
  }

  function schedule(): void {
    clearTimer()
    const conversationId = options.conversationId.value
    const messageId = options.latestMessageId.value
    if (!conversationId || !messageId || !options.enabled.value || !options.stickToBottom.value) {
      return
    }
    timer = setTimeout(() => {
      timer = null
      void mark(conversationId, messageId)
    }, MARK_READ_DEBOUNCE_MS)
  }

  watch(
    () =>
      [
        options.conversationId.value,
        options.latestMessageId.value,
        options.enabled.value,
        options.stickToBottom.value,
      ] as const,
    ([conversationId]) => {
      if (conversationId !== lastConversationId) {
        lastConversationId = conversationId
        lastMarkedId = null
      }
      schedule()
    },
    { immediate: true },
  )

  watch(
    () => options.conversationId.value,
    (id, _, onCleanup) => {
      onCleanup(clearTimer)
      if (!id) {
        lastMarkedId = null
        lastConversationId = null
      }
    },
  )
}
