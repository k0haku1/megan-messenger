import { onScopeDispose } from 'vue'
import { DRAG_THRESHOLD_PX } from './geometry'
import type { LaidOutTile, Rect } from './types'
import { useChatWorkspaceStore } from './chat-workspace.store'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'

/**
 * Drag an already-open tile (chrome / chat header) to retile it elsewhere.
 */
export function useTileDrag() {
  const workspace = useChatWorkspaceStore()
  const selection = useConversationSelectionStore()

  let pointerId: number | null = null
  let startX = 0
  let startY = 0
  let active = false
  let tileId: string | null = null
  let title = ''
  let originViewportRect: Rect | null = null
  let captureEl: HTMLElement | null = null

  function cleanupListeners(): void {
    window.removeEventListener('pointermove', onPointerMove)
    window.removeEventListener('pointerup', onPointerUp)
    window.removeEventListener('pointercancel', onPointerCancel)

    if (captureEl && pointerId !== null) {
      try {
        if (captureEl.hasPointerCapture(pointerId)) {
          captureEl.releasePointerCapture(pointerId)
        }
      } catch {
        // ignore
      }
    }

    pointerId = null
    tileId = null
    title = ''
    originViewportRect = null
    captureEl = null
    active = false
  }

  function onPointerMove(event: PointerEvent): void {
    if (!tileId || pointerId !== event.pointerId) return

    const dx = event.clientX - startX
    const dy = event.clientY - startY

    if (!active && Math.hypot(dx, dy) >= DRAG_THRESHOLD_PX) {
      active = true
      workspace.beginTileDrag(tileId, title, event.clientX, event.clientY)
      event.preventDefault()
    }

    if (!active) return
    workspace.updateTileDragGhost(event.clientX, event.clientY)
  }

  function onPointerUp(event: PointerEvent): void {
    if (!tileId || pointerId !== event.pointerId) return
    const id = tileId
    const rect = originViewportRect
    const wasActive = active
    cleanupListeners()

    if (!wasActive) {
      workspace.endTileDrag()
      return
    }

    const dropped = workspace.dropTileAt(id, event.clientX, event.clientY, rect ?? undefined)
    if (!dropped) {
      workspace.endTileDrag()
      return
    }

    const location = workspace.laidOutTiles.find((tile) => tile.tileId === id)
    if (location) {
      selection.selectedConversationId = location.conversationId
      selection.clearPendingPeer()
    }
  }

  function onPointerCancel(event: PointerEvent): void {
    if (pointerId !== null && event.pointerId !== pointerId) return
    cleanupListeners()
    workspace.endTileDrag()
  }

  function beginTilePointerDrag(
    tile: LaidOutTile,
    tileTitle: string,
    event: PointerEvent,
  ): void {
    if (event.button !== 0) return
    const target = event.target as HTMLElement | null
    if (target?.closest('button, a, textarea, input, [data-no-pane-drag]')) return

    event.preventDefault()
    event.stopPropagation()

    workspace.focusTile(tile.tileId)
    selection.selectedConversationId = tile.conversationId
    selection.clearPendingPeer()

    const bounds = workspace.workspaceBounds
    const current = workspace.laidOutTiles.find((item) => item.tileId === tile.tileId)
    const rect = current?.rect ?? tile.rect

    pointerId = event.pointerId
    startX = event.clientX
    startY = event.clientY
    active = false
    tileId = tile.tileId
    title = tileTitle
    originViewportRect = {
      x: bounds.left + rect.x,
      y: bounds.top + rect.y,
      w: rect.w,
      h: rect.h,
    }

    const el = event.currentTarget
    if (el instanceof HTMLElement) {
      captureEl = el
      try {
        el.setPointerCapture(event.pointerId)
      } catch {
        // ignore — window listeners still work
      }
    }

    window.addEventListener('pointermove', onPointerMove)
    window.addEventListener('pointerup', onPointerUp)
    window.addEventListener('pointercancel', onPointerCancel)
  }

  onScopeDispose(() => {
    cleanupListeners()
    workspace.endTileDrag()
  })

  return {
    beginTilePointerDrag,
    /** @deprecated alias */
    beginChromeDrag: beginTilePointerDrag,
  }
}
