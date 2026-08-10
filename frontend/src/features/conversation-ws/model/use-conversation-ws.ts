import { onScopeDispose, type Ref, watch } from 'vue'
import { conversationApi } from '@/entities/conversation/api/conversation.api'
import { removeMessage, upsertMessage } from '@/entities/message/api/message.repository'
import type { Message } from '@/entities/message/model/types'
import { useReadWatermarkStore } from '@/features/message-read/model/read-watermark.store'
import { useSessionStore } from '@/entities/session/model/session.store'
import { getConversationWsUrl } from '@/shared/api/ws-url'
import { db } from '@/shared/model/db'

const RECONNECT_DELAY_MS = 3_000

type ConversationWsEvent =
  | { type: 'message'; message: Message }
  | { type: 'read'; read: { conversationId: string; userId: string; lastReadAt: string } }

export function useConversationWs(conversationId: Ref<string | null>): void {
  let socket: WebSocket | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let activeConversationId: string | null = null
  const watermarks = useReadWatermarkStore()

  function disconnect(): void {
    activeConversationId = null
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (socket) {
      socket.onopen = null
      socket.onmessage = null
      socket.onclose = null
      socket.onerror = null
      socket.close()
      socket = null
    }
  }

  function scheduleReconnect(conversationIdToReconnect: string): void {
    if (reconnectTimer || activeConversationId !== conversationIdToReconnect) {
      return
    }

    reconnectTimer = setTimeout(() => {
      reconnectTimer = null
      if (activeConversationId === conversationIdToReconnect) {
        connect(conversationIdToReconnect)
      }
    }, RECONNECT_DELAY_MS)
  }

  async function refreshOthersReadAt(id: string): Promise<void> {
    try {
      const state = await conversationApi.readState(id)
      watermarks.setOthersReadAt(id, state.othersReadAt)
    } catch {
      // Keep previous watermark until the next sync.
    }
  }

  function connect(id: string): void {
    disconnect()

    const session = useSessionStore()
    const token = session.accessToken
    if (!token) return

    activeConversationId = id
    socket = new WebSocket(getConversationWsUrl(id, token))

    socket.onmessage = (event) => {
      try {
        const payload = JSON.parse(String(event.data)) as ConversationWsEvent | Message
        if (isEnvelope(payload)) {
          if (payload.type === 'message' && payload.message) {
            void upsertMessageFromEvent(payload.message, session.user?.id)
            return
          }
          if (payload.type === 'read' && payload.read) {
            if (payload.read.userId !== session.user?.id) {
              void refreshOthersReadAt(payload.read.conversationId)
            }
            return
          }
          return
        }

        void upsertMessageFromEvent(payload, session.user?.id)
      } catch {
        // Ignore non-JSON frames such as protocol-level control messages.
      }
    }

    socket.onclose = () => {
      if (activeConversationId === id) {
        scheduleReconnect(id)
      }
    }
  }

  watch(
    conversationId,
    (id) => {
      if (id) {
        connect(id)
        return
      }
      disconnect()
    },
    { immediate: true },
  )

  onScopeDispose(disconnect)
}

function isEnvelope(payload: ConversationWsEvent | Message): payload is ConversationWsEvent {
  return typeof payload === 'object' && payload !== null && 'type' in payload
}

async function upsertMessageFromEvent(message: Message, currentUserId?: string): Promise<void> {
  if (message.deletedAt) {
    await removeMessage(message.id)
    return
  }
  if (message.reactionUpdatedBy && message.reactionUpdatedBy !== currentUserId) {
    const previous = await db.messages.get(message.id)
    if (previous) {
      const ownReactions = new Map(
        (previous.reactions ?? []).map((reaction) => [reaction.emoji, reaction.reactedByMe]),
      )
      message.reactions = (message.reactions ?? []).map((reaction) => ({
        ...reaction,
        reactedByMe: ownReactions.get(reaction.emoji) ?? false,
      }))
    }
  }
  await upsertMessage(message)
}
