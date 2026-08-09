import type { LaidOutTile } from './types'
import { useChatWorkspaceStore } from './chat-workspace.store'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'

export function usePaneInteractions() {
  const workspace = useChatWorkspaceStore()
  const selection = useConversationSelectionStore()

  function activateTile(tile: LaidOutTile): void {
    workspace.focusTile(tile.tileId)
    selection.selectedConversationId = tile.conversationId
    selection.clearPendingPeer()
  }

  function closeTile(tile: LaidOutTile): void {
    selection.closeConversation(tile.conversationId)
  }

  function beginColumnSashDrag(sashIndex: number, event: PointerEvent): void {
    if (event.button !== 0) return
    event.preventDefault()
    event.stopPropagation()

    let lastX = event.clientX
    const pointerId = event.pointerId

    function onMove(moveEvent: PointerEvent): void {
      if (moveEvent.pointerId !== pointerId) return
      const dx = moveEvent.clientX - lastX
      lastX = moveEvent.clientX
      if (dx !== 0) workspace.resizeColumnPair(sashIndex, dx)
    }

    function onUp(upEvent: PointerEvent): void {
      if (upEvent.pointerId !== pointerId) return
      window.removeEventListener('pointermove', onMove)
      window.removeEventListener('pointerup', onUp)
      window.removeEventListener('pointercancel', onUp)
    }

    window.addEventListener('pointermove', onMove)
    window.addEventListener('pointerup', onUp)
    window.addEventListener('pointercancel', onUp)
  }

  function beginRowSashDrag(
    columnIndex: number,
    tileIndex: number,
    event: PointerEvent,
  ): void {
    if (event.button !== 0) return
    event.preventDefault()
    event.stopPropagation()

    let lastY = event.clientY
    const pointerId = event.pointerId

    function onMove(moveEvent: PointerEvent): void {
      if (moveEvent.pointerId !== pointerId) return
      const dy = moveEvent.clientY - lastY
      lastY = moveEvent.clientY
      if (dy !== 0) workspace.resizeRowPair(columnIndex, tileIndex, dy)
    }

    function onUp(upEvent: PointerEvent): void {
      if (upEvent.pointerId !== pointerId) return
      window.removeEventListener('pointermove', onMove)
      window.removeEventListener('pointerup', onUp)
      window.removeEventListener('pointercancel', onUp)
    }

    window.addEventListener('pointermove', onMove)
    window.addEventListener('pointerup', onUp)
    window.addEventListener('pointercancel', onUp)
  }

  return {
    activateTile,
    closeTile,
    beginColumnSashDrag,
    beginRowSashDrag,
  }
}
