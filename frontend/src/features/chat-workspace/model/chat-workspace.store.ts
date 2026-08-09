import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import {
  clientToWorkspacePoint,
  computeDropPreview,
  createId,
  createSingleColumn,
  createTile,
  extractTile,
  findTileLocation,
  isPointInBounds,
  layoutColumnSashes,
  layoutRowSashes,
  layoutWorkspace,
  normalizeWeights,
} from './geometry'
import type {
  ChatColumn,
  ChatTile,
  DropPreview,
  Rect,
  WorkspaceBounds,
} from './types'

function emptyBounds(): WorkspaceBounds {
  return { width: 0, height: 0, left: 0, top: 0 }
}

function cloneColumns(columns: ChatColumn[]): ChatColumn[] {
  return columns.map((column) => ({
    ...column,
    tiles: column.tiles.map((tile) => ({ ...tile })),
  }))
}

export const useChatWorkspaceStore = defineStore('chat-workspace', () => {
  // New key: old floating pane shape is incompatible with tiling columns.
  const columns = useLocalStorage<ChatColumn[]>('megan:chat-columns', [])
  const activeTileId = useLocalStorage<string | null>('megan:active-tile', null)
  const workspaceBounds = ref<WorkspaceBounds>(emptyBounds())
  const enterFromByTileId = ref<Record<string, Rect>>({})
  const listDragActive = ref(false)
  const dropTargetHot = ref(false)
  const dropPreview = ref<DropPreview | null>(null)
  /** Tile currently being rearranged (excluded from drop hit-testing). */
  const draggingTileId = ref<string | null>(null)
  const tileDragGhost = ref<{
    title: string
    clientX: number
    clientY: number
  } | null>(null)

  const displayColumns = computed(() => columns.value)

  const laidOutTiles = computed(() =>
    layoutWorkspace(displayColumns.value, workspaceBounds.value),
  )

  const columnSashes = computed(() =>
    layoutColumnSashes(laidOutTiles.value, displayColumns.value, workspaceBounds.value),
  )

  const rowSashes = computed(() => layoutRowSashes(laidOutTiles.value))

  const activeLaidOut = computed(
    () => laidOutTiles.value.find((tile) => tile.tileId === activeTileId.value) ?? null,
  )

  const openConversationIds = computed(() =>
    columns.value.flatMap((column) => column.tiles.map((tile) => tile.conversationId)),
  )

  const activeConversationId = computed(
    () => activeLaidOut.value?.conversationId ?? null,
  )

  function setWorkspaceBounds(bounds: WorkspaceBounds): void {
    workspaceBounds.value = bounds
  }

  function setColumns(next: ChatColumn[]): void {
    columns.value = next
  }

  function focusTile(tileId: string): void {
    const exists = columns.value.some((column) =>
      column.tiles.some((tile) => tile.id === tileId),
    )
    if (!exists) return
    activeTileId.value = tileId
  }

  function focusConversation(conversationId: string): boolean {
    const location = findTileLocation(columns.value, conversationId)
    if (!location) return false
    activeTileId.value = location.tile.id
    return true
  }

  /** Click semantics: one chat occupies the full workspace. */
  function openExclusive(conversationId: string): string {
    const existing = findTileLocation(columns.value, conversationId)
    if (existing && columns.value.length === 1 && columns.value[0]!.tiles.length === 1) {
      activeTileId.value = existing.tile.id
      return existing.tile.id
    }

    if (existing) {
      // Already tiled somewhere — just focus, keep layout.
      activeTileId.value = existing.tile.id
      return existing.tile.id
    }

    const column = createSingleColumn(conversationId)
    setColumns([column])
    activeTileId.value = column.tiles[0]!.id
    return column.tiles[0]!.id
  }

  function insertTileAtPreview(
    baseColumns: ChatColumn[],
    tile: ChatTile,
    preview: DropPreview,
  ): ChatColumn[] {
    const next = cloneColumns(baseColumns)

    if (next.length === 0) {
      return [
        {
          id: createId('col'),
          weight: 1,
          tiles: [{ ...tile, weight: 1 }],
        },
      ]
    }

    const targetColumn = next[preview.columnIndex]
    if (!targetColumn) {
      return [
        {
          id: createId('col'),
          weight: 1,
          tiles: [{ ...tile, weight: 1 }],
        },
      ]
    }

    if (preview.zone === 'left' || preview.zone === 'right') {
      const stolen = targetColumn.weight / 2
      targetColumn.weight = Math.max(stolen, 0.0001)
      const inserted: ChatColumn = {
        id: createId('col'),
        weight: stolen,
        tiles: [{ ...tile, weight: 1 }],
      }
      const insertAt = preview.zone === 'left' ? preview.columnIndex : preview.columnIndex + 1
      next.splice(insertAt, 0, inserted)
      return normalizeWeights(next)
    }

    const targetTile = targetColumn.tiles[preview.tileIndex]
    if (!targetTile) {
      return normalizeWeights([
        ...next,
        { id: createId('col'), weight: 1, tiles: [{ ...tile, weight: 1 }] },
      ])
    }

    const stolen = targetTile.weight / 2
    targetTile.weight = Math.max(stolen, 0.0001)
    const moved = { ...tile, weight: stolen }
    const insertAt = preview.zone === 'top' ? preview.tileIndex : preview.tileIndex + 1
    targetColumn.tiles.splice(insertAt, 0, moved)
    targetColumn.tiles = normalizeWeights(targetColumn.tiles)
    return next
  }

  function applyDropPreview(
    conversationId: string,
    preview: DropPreview,
    fromViewportRect?: Rect,
  ): string {
    const existing = findTileLocation(columns.value, conversationId)
    if (existing) {
      // Existing tile: rearrange instead of no-op focus.
      return moveTileWithPreview(existing.tile.id, preview, fromViewportRect)
    }

    const bounds = workspaceBounds.value
    const newTile = createTile(conversationId, 1)
    const tileId = newTile.id

    if (fromViewportRect && bounds.width > 0) {
      enterFromByTileId.value = {
        ...enterFromByTileId.value,
        [tileId]: {
          x: fromViewportRect.x - bounds.left,
          y: fromViewportRect.y - bounds.top,
          w: fromViewportRect.w,
          h: fromViewportRect.h,
        },
      }
    }

    if (columns.value.length === 0) {
      setColumns([
        {
          id: createId('col'),
          weight: 1,
          tiles: [newTile],
        },
      ])
      activeTileId.value = tileId
      return tileId
    }

    setColumns(insertTileAtPreview(columns.value, newTile, preview))
    activeTileId.value = tileId
    return tileId
  }

  function moveTileWithPreview(
    tileId: string,
    preview: DropPreview,
    fromViewportRect?: Rect,
  ): string {
    const extracted = extractTile(cloneColumns(columns.value), tileId)
    if (!extracted) return tileId

    const bounds = workspaceBounds.value
    if (fromViewportRect && bounds.width > 0) {
      enterFromByTileId.value = {
        ...enterFromByTileId.value,
        [tileId]: {
          x: fromViewportRect.x - bounds.left,
          y: fromViewportRect.y - bounds.top,
          w: fromViewportRect.w,
          h: fromViewportRect.h,
        },
      }
    }

    setColumns(insertTileAtPreview(extracted.columns, extracted.tile, preview))
    activeTileId.value = tileId
    return tileId
  }

  function clearEnterFrom(tileId: string): void {
    if (!(tileId in enterFromByTileId.value)) return
    const next = { ...enterFromByTileId.value }
    delete next[tileId]
    enterFromByTileId.value = next
  }

  function closeTile(tileId: string): void {
    const next = cloneColumns(columns.value)
    let removed = false

    for (let columnIndex = 0; columnIndex < next.length; columnIndex += 1) {
      const column = next[columnIndex]!
      const tileIndex = column.tiles.findIndex((tile) => tile.id === tileId)
      if (tileIndex < 0) continue

      column.tiles.splice(tileIndex, 1)
      removed = true

      if (column.tiles.length === 0) {
        next.splice(columnIndex, 1)
      } else {
        column.tiles = normalizeWeights(column.tiles)
      }
      break
    }

    if (!removed) return

    const normalized = next.length > 0 ? normalizeWeights(next) : []
    setColumns(normalized)
    clearEnterFrom(tileId)

    if (activeTileId.value !== tileId) return

    const flat = normalized.flatMap((column) => column.tiles)
    activeTileId.value = flat[flat.length - 1]?.id ?? null
  }

  function closeByConversation(conversationId: string): void {
    const location = findTileLocation(columns.value, conversationId)
    if (location) closeTile(location.tile.id)
  }

  function clearAll(): void {
    setColumns([])
    activeTileId.value = null
    enterFromByTileId.value = {}
    dropPreview.value = null
    draggingTileId.value = null
    tileDragGhost.value = null
  }

  function updateDropPreview(clientX: number, clientY: number): void {
    const bounds = workspaceBounds.value
    if (!isPointInBounds(clientX, clientY, bounds)) {
      dropPreview.value = null
      dropTargetHot.value = false
      return
    }

    dropTargetHot.value = true
    const point = clientToWorkspacePoint(clientX, clientY, bounds)
    dropPreview.value = computeDropPreview(columns.value, bounds, point, {
      excludeTileId: draggingTileId.value ?? undefined,
    })
  }

  function clearDropPreview(): void {
    dropPreview.value = null
    dropTargetHot.value = false
  }

  function beginTileDrag(tileId: string, title: string, clientX: number, clientY: number): void {
    draggingTileId.value = tileId
    tileDragGhost.value = { title, clientX, clientY }
    updateDropPreview(clientX, clientY)
  }

  function updateTileDragGhost(clientX: number, clientY: number): void {
    if (!tileDragGhost.value) return
    tileDragGhost.value = { ...tileDragGhost.value, clientX, clientY }
    updateDropPreview(clientX, clientY)
  }

  function endTileDrag(): void {
    draggingTileId.value = null
    tileDragGhost.value = null
    clearDropPreview()
  }

  function dropConversationAt(
    conversationId: string,
    clientX: number,
    clientY: number,
    fromViewportRect?: Rect,
  ): boolean {
    const bounds = workspaceBounds.value
    if (!isPointInBounds(clientX, clientY, bounds)) {
      clearDropPreview()
      return false
    }

    const point = clientToWorkspacePoint(clientX, clientY, bounds)
    const preview =
      dropPreview.value
      ?? computeDropPreview(columns.value, bounds, point, {
        excludeTileId: draggingTileId.value ?? undefined,
      })
    clearDropPreview()
    if (!preview) return false

    applyDropPreview(conversationId, preview, fromViewportRect)
    return true
  }

  function dropTileAt(
    tileId: string,
    clientX: number,
    clientY: number,
    fromViewportRect?: Rect,
  ): boolean {
    const bounds = workspaceBounds.value
    if (!isPointInBounds(clientX, clientY, bounds)) {
      endTileDrag()
      return false
    }

    const point = clientToWorkspacePoint(clientX, clientY, bounds)
    const preview =
      dropPreview.value
      ?? computeDropPreview(columns.value, bounds, point, { excludeTileId: tileId })

    endTileDrag()
    if (!preview) return false

    moveTileWithPreview(tileId, preview, fromViewportRect)
    return true
  }

  function resizeColumnPair(leftIndex: number, deltaPx: number): void {
    const next = cloneColumns(columns.value)
    const left = next[leftIndex]
    const right = next[leftIndex + 1]
    if (!left || !right) return

    const width = workspaceBounds.value.width
    if (width <= 0) return

    const pairWeight = left.weight + right.weight
    const pairPx = width * pairWeight
    if (pairPx <= 0) return

    const leftPx = (left.weight / pairWeight) * pairPx + deltaPx
    const ratio = clampRatio(leftPx / pairPx)
    left.weight = pairWeight * ratio
    right.weight = pairWeight * (1 - ratio)
    setColumns(next)
  }

  function resizeRowPair(columnIndex: number, topTileIndex: number, deltaPx: number): void {
    const next = cloneColumns(columns.value)
    const column = next[columnIndex]
    if (!column) return
    const top = column.tiles[topTileIndex]
    const bottom = column.tiles[topTileIndex + 1]
    if (!top || !bottom) return

    const height = workspaceBounds.value.height
    if (height <= 0) return

    const pairWeight = top.weight + bottom.weight
    const pairPx = height * pairWeight
    if (pairPx <= 0) return

    const topPx = (top.weight / pairWeight) * pairPx + deltaPx
    const ratio = clampRatio(topPx / pairPx)
    top.weight = pairWeight * ratio
    bottom.weight = pairWeight * (1 - ratio)
    column.tiles = [...column.tiles]
    setColumns(next)
  }

  function syncFromSelection(conversationId: string | null): void {
    if (!conversationId) return
    if (findTileLocation(columns.value, conversationId)) {
      focusConversation(conversationId)
      return
    }
    openExclusive(conversationId)
  }

  function setListDragActive(active: boolean): void {
    listDragActive.value = active
    if (!active) clearDropPreview()
  }

  return {
    columns,
    activeTileId,
    laidOutTiles,
    columnSashes,
    rowSashes,
    activeLaidOut,
    openConversationIds,
    activeConversationId,
    workspaceBounds,
    enterFromByTileId,
    listDragActive,
    dropTargetHot,
    dropPreview,
    draggingTileId,
    tileDragGhost,
    setWorkspaceBounds,
    openExclusive,
    openMaximized: openExclusive,
    focusTile,
    focusPane: focusTile,
    focusConversation,
    closeTile,
    closePane: closeTile,
    closeByConversation,
    clearAll,
    clearEnterFrom,
    updateDropPreview,
    clearDropPreview,
    beginTileDrag,
    updateTileDragGhost,
    endTileDrag,
    dropConversationAt,
    dropTileAt,
    resizeColumnPair,
    resizeRowPair,
    syncFromSelection,
    setListDragActive,
    setDropTargetHot: (hot: boolean) => {
      dropTargetHot.value = hot
    },
    isPointInWorkspace: (clientX: number, clientY: number) =>
      isPointInBounds(clientX, clientY, workspaceBounds.value),
  }
})

function clampRatio(value: number): number {
  return Math.min(0.85, Math.max(0.15, value))
}
