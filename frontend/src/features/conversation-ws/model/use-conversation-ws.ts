import { onScopeDispose, type Ref, watch } from 'vue'
import { removeMessage, upsertMessage } from '@/entities/message/api/message.repository'
import type { Message } from '@/entities/message/model/types'
import { useSessionStore } from '@/entities/session/model/session.store'
import { queryClient } from '@/shared/api/query-client'
import { queryKeys } from '@/shared/api/query-keys'
import { getConversationWsUrl } from '@/shared/api/ws-url'

const RECONNECT_DELAY_MS = 3_000

export function useConversationWs(conversationId: Ref<string | null>): void {
  let socket: WebSocket | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let activeConversationId: string | null = null

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

  function connect(id: string): void {
    disconnect()

    const session = useSessionStore()
    const token = session.accessToken
    if (!token) return

    activeConversationId = id
    socket = new WebSocket(getConversationWsUrl(id, token))

    socket.onmessage = (event) => {
      try {
        const message = JSON.parse(String(event.data)) as Message
        void upsertMessageFromEvent(message, session.user?.id).then(() => {
          void queryClient.invalidateQueries({ queryKey: queryKeys.conversations })
        })
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

async function upsertMessageFromEvent(message: Message, currentUserId?: string): Promise<void> {
	if (message.deletedAt) {
		await removeMessage(message.id)
		return
	}
  if (message.reactionUpdatedBy && message.reactionUpdatedBy !== currentUserId) {
    const previous = await db.messages.get(message.id)
    if (previous) {
      const ownReactions = new Map((previous.reactions ?? []).map((reaction) => [reaction.emoji, reaction.reactedByMe]))
      message.reactions = (message.reactions ?? []).map((reaction) => ({
        ...reaction,
        reactedByMe: ownReactions.get(reaction.emoji) ?? false,
      }))
    }
  }
  await upsertMessage(message)
}
