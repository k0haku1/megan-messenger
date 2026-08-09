import { computed, ref, watch, type Ref } from 'vue'
import {
  getSendErrorMessage,
  sendConversationMessage,
  sendDirectMessage,
} from '@/features/send-message/model/send-message'
import type { Message } from '@/entities/message/model/types'
import type { PendingPeer } from '@/features/conversation-selection/model/conversation-selection.store'

type ComposerTarget =
  | { kind: 'conversation'; conversationId: string }
  | { kind: 'pending-peer'; peer: PendingPeer }
  | null

export function useMessageComposer(options: {
  selectedConversationId: Ref<string | null>
  pendingPeer: Ref<PendingPeer | null>
  replyTo: Ref<Message | null>
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
    options.replyTo.value = null
  })

  function clearReply() {
    options.replyTo.value = null
  }

  async function send(): Promise<void> {
    const rawContent = draft.value.trim()
    const activeTarget = target.value
    if (!rawContent || isSending.value || !activeTarget) {
      return
    }

    const replyToId = options.replyTo.value?.id

    isSending.value = true
    sendError.value = null

    try {
      if (activeTarget.kind === 'pending-peer') {
        const { conversationId } = await sendDirectMessage({
          userId: activeTarget.peer.id,
          content,
        })
        draft.value = ''
        options.replyTo.value = null
        options.onConversationOpened(conversationId)
        return
      }

      await sendConversationMessage(activeTarget.conversationId, rawContent, replyToId)
      draft.value = ''
      options.replyTo.value = null
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
    clearReply,
  }
}
