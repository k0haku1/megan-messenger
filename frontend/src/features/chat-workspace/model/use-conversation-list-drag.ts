import { onScopeDispose, ref, type Ref } from 'vue'
import { DRAG_THRESHOLD_PX } from './geometry'
import type { Rect } from './types'

export interface ConversationDragState {
  conversationId: string
  title: string
  clientX: number
  clientY: number
  originRect: Rect
  active: boolean
}

/**
 * Pointer-based drag from conversation list with threshold.
 * Click (no drag) vs drag-to-workspace are distinguished by movement.
 */
export function useConversationListDrag(options: {
  onClick: (conversationId: string) => void
  onDrop: (payload: {
    conversationId: string
    clientX: number
    clientY: number
    originRect: Rect
  }) => void
  onCancel?: () => void
}) {
  const drag = ref<ConversationDragState | null>(null)
  let pointerId: number | null = null
  let startX = 0
  let startY = 0

  function onPointerMove(event: PointerEvent): void {
    if (!drag.value || pointerId !== event.pointerId) return

    const dx = event.clientX - startX
    const dy = event.clientY - startY
    if (!drag.value.active && Math.hypot(dx, dy) >= DRAG_THRESHOLD_PX) {
      drag.value = { ...drag.value, active: true }
      event.preventDefault()
    }

    if (!drag.value.active) return

    drag.value = {
      ...drag.value,
      clientX: event.clientX,
      clientY: event.clientY,
    }
  }

  function finish(event: PointerEvent): void {
    if (!drag.value || pointerId !== event.pointerId) return
    const snapshot = drag.value
    cleanup()

    if (!snapshot.active) {
      options.onClick(snapshot.conversationId)
      return
    }

    options.onDrop({
      conversationId: snapshot.conversationId,
      clientX: event.clientX,
      clientY: event.clientY,
      originRect: snapshot.originRect,
    })
  }

  function cleanup(): void {
    window.removeEventListener('pointermove', onPointerMove)
    window.removeEventListener('pointerup', finish)
    window.removeEventListener('pointercancel', cancel)
    pointerId = null
    drag.value = null
  }

  function cancel(event?: PointerEvent): void {
    if (event && pointerId !== null && event.pointerId !== pointerId) return
    cleanup()
    options.onCancel?.()
  }

  function startDrag(payload: {
    event: PointerEvent
    conversationId: string
    title: string
    originEl: HTMLElement
  }): void {
    if (payload.event.button !== 0) return

    const rect = payload.originEl.getBoundingClientRect()
    pointerId = payload.event.pointerId
    startX = payload.event.clientX
    startY = payload.event.clientY
    drag.value = {
      conversationId: payload.conversationId,
      title: payload.title,
      clientX: payload.event.clientX,
      clientY: payload.event.clientY,
      originRect: {
        x: rect.left,
        y: rect.top,
        w: rect.width,
        h: rect.height,
      },
      active: false,
    }

    window.addEventListener('pointermove', onPointerMove)
    window.addEventListener('pointerup', finish)
    window.addEventListener('pointercancel', cancel)
  }

  onScopeDispose(cleanup)

  return {
    drag: drag as Ref<ConversationDragState | null>,
    startDrag,
    cancelDrag: cancel,
  }
}
