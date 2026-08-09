import { computed, ref, watch, type Ref } from 'vue'
import {
  getSendErrorMessage,
  sendConversationMessage,
  sendDirectMessage,
} from '@/features/send-message/model/send-message'
import { messageApi } from '@/entities/message/api/message.api'
import type { Message, MessageAttachment } from '@/entities/message/model/types'
import type { PendingPeer } from '@/features/conversation-selection/model/conversation-selection.store'

type ComposerTarget =
  | { kind: 'conversation'; conversationId: string }
  | { kind: 'pending-peer'; peer: PendingPeer }
  | null

export type PendingAttachment = {
  localId: string
  file: File
  previewUrl?: string
  uploading: boolean
  attachment?: MessageAttachment
  error?: string
}

export function useMessageComposer(options: {
  selectedConversationId: Ref<string | null>
  pendingPeer: Ref<PendingPeer | null>
  replyTo: Ref<Message | null>
  onConversationOpened: (conversationId: string) => void
  onSent?: () => void
}) {
  const draft = ref('')
  const isSending = ref(false)
  const sendError = ref<string | null>(null)
  const pendingAttachments = ref<PendingAttachment[]>([])

  const target = computed<ComposerTarget>(() => {
    if (options.pendingPeer.value) {
      return { kind: 'pending-peer', peer: options.pendingPeer.value }
    }

    if (options.selectedConversationId.value) {
      return { kind: 'conversation', conversationId: options.selectedConversationId.value }
    }

    return null
  })

  const hasReadyAttachments = computed(() =>
    pendingAttachments.value.some((item) => item.attachment && !item.uploading && !item.error),
  )

  const canSend = computed(
    () =>
      (draft.value.trim().length > 0 || hasReadyAttachments.value)
      && !isSending.value
      && target.value !== null
      && !pendingAttachments.value.some((item) => item.uploading),
  )

  watch([options.selectedConversationId, options.pendingPeer], () => {
    draft.value = ''
    sendError.value = null
    options.replyTo.value = null
    clearAttachments()
  })

  function clearReply() {
    options.replyTo.value = null
  }

  function clearAttachments() {
    for (const item of pendingAttachments.value) {
      if (item.previewUrl) URL.revokeObjectURL(item.previewUrl)
    }
    pendingAttachments.value = []
  }

  function removeAttachment(localId: string) {
    const item = pendingAttachments.value.find((entry) => entry.localId === localId)
    if (item?.previewUrl) URL.revokeObjectURL(item.previewUrl)
    pendingAttachments.value = pendingAttachments.value.filter((entry) => entry.localId !== localId)
  }

  async function addFiles(files: FileList | File[]) {
    const conversationId =
      target.value?.kind === 'conversation' ? target.value.conversationId : null
    if (!conversationId) {
      sendError.value = 'Вложения доступны после создания чата'
      return
    }

    const list = Array.from(files)
    for (const file of list) {
      const localId = crypto.randomUUID()
      const previewUrl = file.type.startsWith('image/') || file.type.startsWith('video/')
        ? URL.createObjectURL(file)
        : undefined
      const pending: PendingAttachment = {
        localId,
        file,
        previewUrl,
        uploading: true,
      }
      pendingAttachments.value = [...pendingAttachments.value, pending]

      try {
        const { attachment } = await messageApi.uploadAttachment(conversationId, file)
        pendingAttachments.value = pendingAttachments.value.map((entry) =>
          entry.localId === localId
            ? { ...entry, uploading: false, attachment }
            : entry,
        )
      } catch {
        pendingAttachments.value = pendingAttachments.value.map((entry) =>
          entry.localId === localId
            ? { ...entry, uploading: false, error: 'Не удалось загрузить файл' }
            : entry,
        )
      }
    }
  }

  async function send(): Promise<void> {
    const rawContent = draft.value.trim()
    const activeTarget = target.value
    const attachmentIds = pendingAttachments.value
      .map((item) => item.attachment?.id)
      .filter((id): id is string => Boolean(id))

    if ((!rawContent && attachmentIds.length === 0) || isSending.value || !activeTarget) {
      return
    }
    if (pendingAttachments.value.some((item) => item.uploading)) {
      return
    }

    const replyToId = options.replyTo.value?.id

    isSending.value = true
    sendError.value = null

    try {
      if (activeTarget.kind === 'pending-peer') {
        if (attachmentIds.length > 0) {
          sendError.value = 'Сначала отправьте текст, чтобы создать диалог'
          return
        }
        const { conversationId } = await sendDirectMessage({
          userId: activeTarget.peer.id,
          content: rawContent,
        })
        draft.value = ''
        options.replyTo.value = null
        options.onConversationOpened(conversationId)
        options.onSent?.()
        return
      }

      await sendConversationMessage(
        activeTarget.conversationId,
        rawContent,
        replyToId,
        attachmentIds.length ? attachmentIds : undefined,
      )
      draft.value = ''
      options.replyTo.value = null
      clearAttachments()
      options.onSent?.()
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
    pendingAttachments,
    send,
    clearReply,
    addFiles,
    removeAttachment,
    clearAttachments,
  }
}
