export interface Rect {
  x: number
  y: number
  w: number
  h: number
}

export interface Point {
  x: number
  y: number
}

export interface WorkspaceBounds {
  width: number
  height: number
  left: number
  top: number
}

/** Leaf chat inside a column (vertical stack, Niri-style). */
export interface ChatTile {
  id: string
  conversationId: string
  weight: number
}

/** Horizontal column of the workspace strip. */
export interface ChatColumn {
  id: string
  weight: number
  tiles: ChatTile[]
}

export type SplitZone = 'left' | 'right' | 'top' | 'bottom'

export interface LaidOutTile {
  columnId: string
  tileId: string
  columnIndex: number
  tileIndex: number
  conversationId: string
  rect: Rect
}

export interface DropPreview {
  zone: SplitZone
  /** Exact area the new chat will occupy */
  newRect: Rect
  /** Remaining area of the tile being split (null when filling empty workspace) */
  remainRect: Rect | null
  columnIndex: number
  tileIndex: number
}

export interface ColumnSash {
  index: number
  x: number
  y: number
  h: number
}

export interface RowSash {
  columnIndex: number
  tileIndex: number
  x: number
  y: number
  w: number
}
