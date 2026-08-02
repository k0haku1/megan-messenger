import { onScopeDispose, type Ref, watch } from 'vue'
import { upsertMessage } from '@/entities/message/api/message.repository'
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
        void upsertMessage(message).then(() => {
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
