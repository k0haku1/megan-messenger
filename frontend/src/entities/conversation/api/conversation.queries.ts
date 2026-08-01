import { useQuery } from '@tanstack/vue-query'
import { computed, watch } from 'vue'
import { conversationApi } from './conversation.api'
import { replaceConversations } from './conversation.repository'
import { db } from '@/shared/model/db'
import type { SyncState } from '@/shared/model/db'
import { useSessionStore } from '@/entities/session/model/session.store'

const conversationsKey = ['conversations'] as const

async function setSyncState(state: SyncState): Promise<void> { await db.syncStates.put(state) }

/** Fetches only into Dexie. Query data is intentionally not exposed to the UI. */
export function useConversationsSync() {
  const session = useSessionStore()
  const enabled = computed(() => session.isAuthenticated && !session.needsUsername)

  const query = useQuery({
    queryKey: conversationsKey,
    queryFn: conversationApi.list,
    enabled,
  })

  watch(query.fetchStatus, (fetchStatus) => {
    if (fetchStatus === 'fetching') void setSyncState({ key: 'conversations', status: 'loading', updatedAt: Date.now() })
  }, { immediate: true })
  watch(query.data, (data) => { if (data) void replaceConversations(data.conversations).then(() => setSyncState({ key: 'conversations', status: 'success', updatedAt: Date.now() })) }, { immediate: true })
  watch(query.error, (error) => {
    if (error) void setSyncState({ key: 'conversations', status: 'error', updatedAt: Date.now(), error: error instanceof Error ? error.message : 'Unknown error' })
  }, { immediate: true })
}
