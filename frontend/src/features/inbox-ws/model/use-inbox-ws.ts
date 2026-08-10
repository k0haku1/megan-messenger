import { onScopeDispose, watch } from 'vue'
import {
  applyConversationActivity,
  type ConversationActivity,
} from '@/entities/conversation/lib/apply-conversation-activity'
import { useSessionStore } from '@/entities/session/model/session.store'
import { getInboxWsUrl } from '@/shared/api/ws-url'

const RECONNECT_DELAY_MS = 3_000

type InboxWsEvent = {
  type: 'conversation'
  conversation: ConversationActivity
}

/**
 * User-scoped inbox socket: keeps the chat list in sync without full refetches.
 */
export function useInboxWs(): void {
  let socket: WebSocket | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let active = false

  function disconnect(): void {
    active = false
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

  function scheduleReconnect(): void {
    if (reconnectTimer || !active) return
    reconnectTimer = setTimeout(() => {
      reconnectTimer = null
      if (active) connect()
    }, RECONNECT_DELAY_MS)
  }

  function connect(): void {
    disconnect()
    const session = useSessionStore()
    const token = session.accessToken
    if (!token || !session.isAuthenticated) return

    active = true
    socket = new WebSocket(getInboxWsUrl(token))

    socket.onmessage = (event) => {
      try {
        const payload = JSON.parse(String(event.data)) as InboxWsEvent
        if (payload.type !== 'conversation' || !payload.conversation) return

        const session = useSessionStore()
        applyConversationActivity(payload.conversation, {
          currentUserId: session.user?.id,
        })
      } catch {
        // Ignore non-JSON frames.
      }
    }

    socket.onclose = () => {
      if (active) scheduleReconnect()
    }
  }

  const session = useSessionStore()
  watch(
    () => [session.isAuthenticated, session.accessToken] as const,
    ([isAuthenticated, token]) => {
      if (isAuthenticated && token) {
        connect()
        return
      }
      disconnect()
    },
    { immediate: true },
  )

  onScopeDispose(disconnect)
}
