import { computed, ref, watch, type Ref } from 'vue'
import {
  getSendErrorMessage,
  sendConversationMessage,
  sendDirectMessage,
} from '@/features/send-message/model/send-message'
import type { PendingPeer } from '@/features/conversation-selection/model/conversation-selection.store'

type ComposerTarget =
  | { kind: 'conversation'; conversationId: string }
  | { kind: 'pending-peer'; peer: PendingPeer }
  | null

export function useMessageComposer(options: {
  selectedConversationId: Ref<string | null>
  pendingPeer: Ref<PendingPeer | null>
  onConversationOpened: (conversationId: string) => void
}) {
  const draft = ref('')
  const isSending = ref(false)
  const sendError = ref<string | null>(null)

  const target = computed<ComposerTarget>(() => {
    if (options.pendingPeer.value) {
      return { kind: 'pending-peer', peer: options.pendingPeer.value }
    }

    if (options.selectedConversationId.value) {
      return { kind: 'conversation', conversationId: options.selectedConversationId.value }
    }

    return null
  })

  const canSend = computed(
    () => draft.value.trim().length > 0 && !isSending.value && target.value !== null,
  )

  watch([options.selectedConversationId, options.pendingPeer], () => {
    draft.value = ''
    sendError.value = null
  })

  async function send(): Promise<void> {
    const content = draft.value.trim()
    const activeTarget = target.value
    if (!content || isSending.value || !activeTarget) {
      return
    }

    isSending.value = true
    sendError.value = null

    try {
      if (activeTarget.kind === 'pending-peer') {
        const { conversationId } = await sendDirectMessage({
          userId: activeTarget.peer.id,
          content,
        })
        draft.value = ''
        options.onConversationOpened(conversationId)
        return
      }

      await sendConversationMessage(activeTarget.conversationId, content)
      draft.value = ''
    } catch (error) {
      sendError.value = getSendErrorMessage(error)
    } finally {
      isSending.value = false
    }
  }

  return {
    draft,
    isSending,
    sendError,
    canSend,
    send,
  }
}
