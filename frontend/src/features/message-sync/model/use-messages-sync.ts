import { nextTick, ref, watch, type Ref } from 'vue'
import { messageApi } from '@/entities/message/api/message.api'
import {
  replaceMessagesForConversation,
  upsertMessages,
} from '@/entities/message/api/message.repository'
import type { Message } from '@/entities/message/model/types'
import { useReadWatermarkStore } from '@/features/message-read/model/read-watermark.store'

const INITIAL_LIMIT = 50
const OLDER_LIMIT = 50
const MAX_OPEN_PAGES = 20

function newestFirstId(messages: Message[]): string | null {
  // API returns DESC (newest first).
  return messages[0]?.id ?? null
}

function findFirstUnreadId(
  messagesAsc: Array<{ id: string; createdAt: string; senderId: string }>,
  currentUserId: string | null | undefined,
  lastReadAt: string | null | undefined,
): string | null {
  const watermark = lastReadAt ? new Date(lastReadAt).getTime() : Number.NEGATIVE_INFINITY
  if (Number.isNaN(watermark)) return null

  for (const message of messagesAsc) {
    if (currentUserId && message.senderId === currentUserId) continue
    const created = new Date(message.createdAt).getTime()
    if (Number.isNaN(created)) continue
    if (created > watermark) return message.id
  }
  return null
}

async function waitFrames(count: number, cancelled: () => boolean) {
  for (let i = 0; i < count; i += 1) {
    if (cancelled()) return
    await new Promise<void>((resolve) => {
      requestAnimationFrame(() => resolve())
    })
  }
}

/**
 * Syncs the latest message page into Dexie and loads older pages on demand.
 * Dexie stays as the reactive source for WS/send; open-scroll waits on sync flags.
 */
export function useMessagesSync(
  conversationId: Ref<string | null>,
  options?: {
    messagesAsc?: Ref<Array<{ id: string; createdAt: string; senderId: string }>>
    currentUserId?: Ref<string | null | undefined>
    lastReadAt?: Ref<string | null | undefined>
    readStateReady?: Ref<boolean>
    openUnreadCount?: Ref<number>
  },
) {
  const watermarks = useReadWatermarkStore()
  const messagesSynced = ref(false)
  const openHistoryReady = ref(false)
  const syncedLatestMessageId = ref<string | null>(null)
  const syncedMessageCount = ref(0)
  const nextCursor = ref<string | null>(null)
  const hasMoreOlder = ref(false)
  const isLoadingOlder = ref(false)

  async function waitForLocalCount(id: string, minCount: number, cancelled: () => boolean) {
    const local = options?.messagesAsc
    if (!local) {
      await nextTick()
      return
    }
    for (let attempt = 0; attempt < 40; attempt += 1) {
      if (cancelled() || conversationId.value !== id) return
      if (local.value.length >= minCount) return
      await waitFrames(1, cancelled)
    }
  }

  async function waitForReadState(id: string, cancelled: () => boolean) {
    if (!options?.readStateReady) return
    for (let attempt = 0; attempt < 60; attempt += 1) {
      if (cancelled() || conversationId.value !== id) return
      if (options.readStateReady.value) return
      await waitFrames(1, cancelled)
    }
  }

  async function loadOlder(container?: HTMLElement | null): Promise<boolean> {
    const id = conversationId.value
    const cursor = nextCursor.value
    if (!id || !cursor || isLoadingOlder.value || !hasMoreOlder.value) return false

    isLoadingOlder.value = true
    const prevHeight = container?.scrollHeight ?? 0
    const prevTop = container?.scrollTop ?? 0

    try {
      const { messages, nextCursor: olderCursor } = await messageApi.list(id, {
        limit: OLDER_LIMIT,
        cursor,
      })
      if (conversationId.value !== id) return false
      await upsertMessages(messages)
      nextCursor.value = olderCursor || null
      hasMoreOlder.value = Boolean(olderCursor)
      syncedMessageCount.value += messages.length

      if (container) {
        await nextTick()
        await waitFrames(1, () => conversationId.value !== id)
        container.scrollTop = container.scrollHeight - prevHeight + prevTop
      }
      return messages.length > 0
    } catch {
      return false
    } finally {
      isLoadingOlder.value = false
    }
  }

  /**
   * Newest page alone may contain only own messages while unread peers are older.
   * Keep paging upward until the first unread sits below some older context (or history ends).
   */
  async function ensureOpenHistory(id: string, cancelled: () => boolean): Promise<void> {
    const unreadHint = options?.openUnreadCount?.value ?? 0
    const lastReadAt = options?.lastReadAt?.value
    const currentUserId = options?.currentUserId?.value

    for (let page = 0; page < MAX_OPEN_PAGES; page += 1) {
      if (cancelled() || conversationId.value !== id) return

      const local = options?.messagesAsc?.value ?? []
      const firstUnreadId = findFirstUnreadId(local, currentUserId, lastReadAt)
      const oldest = local[0]

      if (firstUnreadId && oldest && firstUnreadId !== oldest.id) {
        // Messages before the first unread are loaded — safe to anchor open scroll.
        return
      }
      if (firstUnreadId && !hasMoreOlder.value) return

      if (!firstUnreadId) {
        if (!hasMoreOlder.value || !nextCursor.value) return
        if (lastReadAt && oldest) {
          const oldestAt = new Date(oldest.createdAt).getTime()
          const watermark = new Date(lastReadAt).getTime()
          if (!Number.isNaN(oldestAt) && !Number.isNaN(watermark) && oldestAt <= watermark) {
            return
          }
        }
        if (!lastReadAt && unreadHint <= 0) return
      }

      if (!hasMoreOlder.value || !nextCursor.value) return

      const loaded = await loadOlder()
      if (!loaded) return
      await waitForLocalCount(id, syncedMessageCount.value, cancelled)
    }
  }

  watch(
    conversationId,
    (id, _, onCleanup) => {
      messagesSynced.value = false
      openHistoryReady.value = false
      syncedLatestMessageId.value = null
      syncedMessageCount.value = 0
      nextCursor.value = null
      hasMoreOlder.value = false
      isLoadingOlder.value = false

      if (!id) {
        messagesSynced.value = true
        openHistoryReady.value = true
        return
      }

      let cancelled = false
      onCleanup(() => {
        cancelled = true
      })

      void (async () => {
        try {
          const { messages, nextCursor: cursor, othersReadAt } = await messageApi.list(id, {
            limit: INITIAL_LIMIT,
          })
          if (cancelled) return

          watermarks.setOthersReadAt(id, othersReadAt)
          await replaceMessagesForConversation(id, messages)
          if (cancelled || conversationId.value !== id) return

          syncedLatestMessageId.value = newestFirstId(messages)
          syncedMessageCount.value = messages.length
          nextCursor.value = cursor || null
          hasMoreOlder.value = Boolean(cursor)
          messagesSynced.value = true

          await waitForLocalCount(id, messages.length, () => cancelled)
          if (cancelled || conversationId.value !== id) return

          await waitForReadState(id, () => cancelled)
          if (cancelled || conversationId.value !== id) return

          await ensureOpenHistory(id, () => cancelled)
          if (cancelled || conversationId.value !== id) return

          openHistoryReady.value = true
        } catch {
          if (!cancelled && conversationId.value === id) {
            messagesSynced.value = true
            openHistoryReady.value = true
          }
        }
      })()
    },
    { immediate: true },
  )

  return {
    messagesSynced,
    openHistoryReady,
    syncedLatestMessageId,
    syncedMessageCount,
    hasMoreOlder,
    isLoadingOlder,
    loadOlder,
  }
}
