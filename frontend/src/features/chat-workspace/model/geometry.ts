import type {
  ChatColumn,
  ChatTile,
  ColumnSash,
  DropPreview,
  LaidOutTile,
  Point,
  Rect,
  RowSash,
  SplitZone,
  WorkspaceBounds,
} from './types'

export const MIN_PANE_WIDTH = 280
export const MIN_PANE_HEIGHT = 220
export const TILE_GAP = 4
export const DRAG_THRESHOLD_PX = 6
export const DROP_EDGE_RATIO = 0.5

export function createId(prefix: string): string {
  return `${prefix}-${crypto.randomUUID()}`
}

export function clamp(value: number, min: number, max: number): number {
  if (max < min) return min
  return Math.min(max, Math.max(min, value))
}

export function isPointInBounds(
  clientX: number,
  clientY: number,
  bounds: WorkspaceBounds | null,
): boolean {
  if (!bounds || bounds.width <= 0 || bounds.height <= 0) return false
  return (
    clientX >= bounds.left
    && clientX <= bounds.left + bounds.width
    && clientY >= bounds.top
    && clientY <= bounds.top + bounds.height
  )
}

export function clientToWorkspacePoint(
  clientX: number,
  clientY: number,
  bounds: WorkspaceBounds,
): Point {
  return {
    x: clientX - bounds.left,
    y: clientY - bounds.top,
  }
}

export function totalWeight(items: { weight: number }[]): number {
  return items.reduce((sum, item) => sum + Math.max(item.weight, 0.0001), 0)
}

export function normalizeWeights<T extends { weight: number }>(items: T[]): T[] {
  if (items.length === 0) return items
  const sum = totalWeight(items)
  return items.map((item) => ({
    ...item,
    weight: Math.max(item.weight, 0.0001) / sum,
  }))
}

function distributeSizes(weights: number[], total: number, gap: number, minSize: number): number[] {
  const count = weights.length
  if (count === 0) return []
  if (count === 1) return [total]

  const gapTotal = gap * (count - 1)
  const available = Math.max(0, total - gapTotal)
  const weightSum = weights.reduce((sum, weight) => sum + weight, 0) || 1

  const raw = weights.map((weight) => (available * weight) / weightSum)
  const sizes = raw.map((size) => Math.max(minSize, size))
  let overflow = sizes.reduce((sum, size) => sum + size, 0) - available

  if (overflow <= 0) return sizes

  // Shrink panes that are above min to absorb overflow.
  while (overflow > 0.5) {
    const flexible = sizes
      .map((size, index) => ({ size, index }))
      .filter((item) => item.size > minSize + 0.5)
    if (flexible.length === 0) break

    const share = overflow / flexible.length
    for (const item of flexible) {
      const reducible = item.size - minSize
      const delta = Math.min(reducible, share)
      sizes[item.index] -= delta
      overflow -= delta
    }
  }

  return sizes
}

export function layoutWorkspace(
  columns: ChatColumn[],
  bounds: Pick<WorkspaceBounds, 'width' | 'height'>,
  gap = TILE_GAP,
): LaidOutTile[] {
  if (columns.length === 0 || bounds.width <= 0 || bounds.height <= 0) return []

  const colWidths = distributeSizes(
    columns.map((column) => column.weight),
    bounds.width,
    gap,
    MIN_PANE_WIDTH,
  )

  const result: LaidOutTile[] = []
  let x = 0

  columns.forEach((column, columnIndex) => {
    const width = colWidths[columnIndex] ?? MIN_PANE_WIDTH
    const rowHeights = distributeSizes(
      column.tiles.map((tile) => tile.weight),
      bounds.height,
      gap,
      MIN_PANE_HEIGHT,
    )

    let y = 0
    column.tiles.forEach((tile, tileIndex) => {
      const height = rowHeights[tileIndex] ?? MIN_PANE_HEIGHT
      result.push({
        columnId: column.id,
        tileId: tile.id,
        columnIndex,
        tileIndex,
        conversationId: tile.conversationId,
        rect: { x, y, w: width, h: height },
      })
      y += height + (tileIndex < column.tiles.length - 1 ? gap : 0)
    })

    x += width + (columnIndex < columns.length - 1 ? gap : 0)
  })

  return result
}

export function layoutColumnSashes(
  laidOut: LaidOutTile[],
  columns: ChatColumn[],
  bounds: Pick<WorkspaceBounds, 'height'>,
  gap = TILE_GAP,
): ColumnSash[] {
  if (columns.length < 2) return []

  const sashes: ColumnSash[] = []
  for (let index = 0; index < columns.length - 1; index += 1) {
    const leftTiles = laidOut.filter((tile) => tile.columnIndex === index)
    const rightEdge = leftTiles.reduce((max, tile) => Math.max(max, tile.rect.x + tile.rect.w), 0)
    sashes.push({
      index,
      x: rightEdge,
      y: 0,
      h: bounds.height,
    })
    void gap
  }
  return sashes
}

export function layoutRowSashes(laidOut: LaidOutTile[], gap = TILE_GAP): RowSash[] {
  const sashes: RowSash[] = []
  const byColumn = new Map<number, LaidOutTile[]>()

  for (const tile of laidOut) {
    const list = byColumn.get(tile.columnIndex) ?? []
    list.push(tile)
    byColumn.set(tile.columnIndex, list)
  }

  for (const [columnIndex, tiles] of byColumn) {
    const ordered = [...tiles].sort((a, b) => a.tileIndex - b.tileIndex)
    for (let index = 0; index < ordered.length - 1; index += 1) {
      const top = ordered[index]!
      sashes.push({
        columnIndex,
        tileIndex: index,
        x: top.rect.x,
        y: top.rect.y + top.rect.h,
        w: top.rect.w,
      })
      void gap
    }
  }

  return sashes
}

export function zoneFromPoint(localX: number, localY: number, width: number, height: number): SplitZone {
  const nx = width <= 0 ? 0.5 : clamp(localX / width, 0, 1)
  const ny = height <= 0 ? 0.5 : clamp(localY / height, 0, 1)

  // Prefer horizontal splits (Niri columns) unless clearly in top/bottom band.
  if (ny < 0.28) return 'top'
  if (ny > 0.72) return 'bottom'
  if (nx < DROP_EDGE_RATIO) return 'left'
  return 'right'
}

export function splitRect(rect: Rect, zone: SplitZone): { newRect: Rect; remainRect: Rect } {
  const halfW = rect.w / 2
  const halfH = rect.h / 2

  switch (zone) {
    case 'left':
      return {
        newRect: { x: rect.x, y: rect.y, w: halfW, h: rect.h },
        remainRect: { x: rect.x + halfW, y: rect.y, w: halfW, h: rect.h },
      }
    case 'right':
      return {
        newRect: { x: rect.x + halfW, y: rect.y, w: halfW, h: rect.h },
        remainRect: { x: rect.x, y: rect.y, w: halfW, h: rect.h },
      }
    case 'top':
      return {
        newRect: { x: rect.x, y: rect.y, w: rect.w, h: halfH },
        remainRect: { x: rect.x, y: rect.y + halfH, w: rect.w, h: halfH },
      }
    case 'bottom':
      return {
        newRect: { x: rect.x, y: rect.y + halfH, w: rect.w, h: halfH },
        remainRect: { x: rect.x, y: rect.y, w: rect.w, h: halfH },
      }
  }
}

export function createSingleColumn(conversationId: string): ChatColumn {
  return {
    id: createId('col'),
    weight: 1,
    tiles: [
      {
        id: createId('tile'),
        conversationId,
        weight: 1,
      },
    ],
  }
}

export function createTile(conversationId: string, weight = 1): ChatTile {
  return {
    id: createId('tile'),
    conversationId,
    weight,
  }
}

export function findTileLocation(
  columns: ChatColumn[],
  conversationId: string,
): { columnIndex: number; tileIndex: number; tile: ChatTile } | null {
  for (let columnIndex = 0; columnIndex < columns.length; columnIndex += 1) {
    const column = columns[columnIndex]!
    const tileIndex = column.tiles.findIndex((tile) => tile.conversationId === conversationId)
    if (tileIndex >= 0) {
      return { columnIndex, tileIndex, tile: column.tiles[tileIndex]! }
    }
  }
  return null
}

export function findTileLocationById(
  columns: ChatColumn[],
  tileId: string,
): { columnIndex: number; tileIndex: number; tile: ChatTile; column: ChatColumn } | null {
  for (let columnIndex = 0; columnIndex < columns.length; columnIndex += 1) {
    const column = columns[columnIndex]!
    const tileIndex = column.tiles.findIndex((tile) => tile.id === tileId)
    if (tileIndex >= 0) {
      return { columnIndex, tileIndex, tile: column.tiles[tileIndex]!, column }
    }
  }
  return null
}

/** Remove a tile and collapse empty columns. Does not normalize sibling weights beyond the affected column. */
export function extractTile(
  columns: ChatColumn[],
  tileId: string,
): { columns: ChatColumn[]; tile: ChatTile } | null {
  const next = columns.map((column) => ({
    ...column,
    tiles: column.tiles.map((tile) => ({ ...tile })),
  }))

  for (let columnIndex = 0; columnIndex < next.length; columnIndex += 1) {
    const column = next[columnIndex]!
    const tileIndex = column.tiles.findIndex((tile) => tile.id === tileId)
    if (tileIndex < 0) continue

    const [tile] = column.tiles.splice(tileIndex, 1)
    if (!tile) return null

    if (column.tiles.length === 0) {
      next.splice(columnIndex, 1)
    } else {
      column.tiles = normalizeWeights(column.tiles)
    }

    return {
      columns: next.length > 0 ? normalizeWeights(next) : [],
      tile,
    }
  }

  return null
}

export function computeDropPreview(
  columns: ChatColumn[],
  bounds: Pick<WorkspaceBounds, 'width' | 'height'>,
  point: Point,
  options?: { excludeTileId?: string },
): DropPreview | null {
  if (bounds.width <= 0 || bounds.height <= 0) return null

  let layoutColumns = columns
  if (options?.excludeTileId) {
    const extracted = extractTile(columns, options.excludeTileId)
    layoutColumns = extracted?.columns ?? columns
  }

  if (layoutColumns.length === 0) {
    const full = { x: 0, y: 0, w: bounds.width, h: bounds.height }
    return {
      zone: 'right',
      newRect: full,
      remainRect: null,
      columnIndex: 0,
      tileIndex: 0,
    }
  }

  const laidOut = layoutWorkspace(layoutColumns, bounds)
  const hit = laidOut.find(
    (tile) =>
      point.x >= tile.rect.x
      && point.x <= tile.rect.x + tile.rect.w
      && point.y >= tile.rect.y
      && point.y <= tile.rect.y + tile.rect.h,
  ) ?? laidOut[laidOut.length - 1]

  if (!hit) return null

  const zone = zoneFromPoint(point.x - hit.rect.x, point.y - hit.rect.y, hit.rect.w, hit.rect.h)
  const { newRect, remainRect } = splitRect(hit.rect, zone)

  return {
    zone,
    newRect,
    remainRect,
    columnIndex: hit.columnIndex,
    tileIndex: hit.tileIndex,
  }
}
