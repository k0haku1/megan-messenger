import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useLocalStorage } from '@vueuse/core'

export interface PendingPeer {
  id: string
  username: string
  avatarUrl?: string
}

/** Ephemeral UI state. Messages remain in Dexie. */
export const useConversationSelectionStore = defineStore('conversation-selection', () => {
  const selectedConversationId = useLocalStorage<string | null>(
    'megan:selected-conversation',
    null,
  )
  const activeFolder = useLocalStorage('megan:active-folder', 'all')
  const pendingPeer = ref<PendingPeer | null>(null)

  function openPendingPeer(peer: PendingPeer): void {
    pendingPeer.value = peer
    selectedConversationId.value = null
  }

  function select(conversationId: string | null): void {
    pendingPeer.value = null
    selectedConversationId.value = conversationId
  }

  function clearPendingPeer(): void {
    pendingPeer.value = null
  }

  function selectFolder(folderId: string): void {
    activeFolder.value = folderId
  }

  return {
    selectedConversationId,
    activeFolder,
    pendingPeer,
    select,
    openPendingPeer,
    clearPendingPeer,
    selectFolder,
  }
})
