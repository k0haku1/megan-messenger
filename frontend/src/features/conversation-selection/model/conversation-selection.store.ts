import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import { useChatWorkspaceStore } from '@/features/chat-workspace/model/chat-workspace.store'

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

  function syncSelectedFromWorkspace(): void {
    const workspace = useChatWorkspaceStore()
    selectedConversationId.value = workspace.activeConversationId
  }

  function openPendingPeer(peer: PendingPeer): void {
    const workspace = useChatWorkspaceStore()
    workspace.clearAll()
    pendingPeer.value = peer
    selectedConversationId.value = null
  }

  function select(conversationId: string | null): void {
    const workspace = useChatWorkspaceStore()
    pendingPeer.value = null

    if (!conversationId) {
      workspace.clearAll()
      selectedConversationId.value = null
      return
    }

    workspace.openMaximized(conversationId)
    selectedConversationId.value = conversationId
  }

  function closeConversation(conversationId: string): void {
    const workspace = useChatWorkspaceStore()
    workspace.closeByConversation(conversationId)
    pendingPeer.value = null
    syncSelectedFromWorkspace()
  }

  function clearPendingPeer(): void {
    pendingPeer.value = null
  }

  function selectFolder(folderId: string): void {
    activeFolder.value = folderId
  }

  /** Focus an already-open pane without forcing maximized layout. */
  function focusConversation(conversationId: string): void {
    const workspace = useChatWorkspaceStore()
    pendingPeer.value = null
    if (workspace.focusConversation(conversationId)) {
      selectedConversationId.value = conversationId
      return
    }
    select(conversationId)
  }

  return {
    selectedConversationId,
    activeFolder,
    pendingPeer,
    select,
    closeConversation,
    openPendingPeer,
    clearPendingPeer,
    selectFolder,
    focusConversation,
    syncSelectedFromWorkspace,
  }
})
