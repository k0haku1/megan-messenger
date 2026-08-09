<template>
  <section
    ref="rootEl"
    class="chat-workspace"
    aria-label="Рабочая область чатов"
    :class="{ 'is-drop-target': dropTargetHot }"
  >
    <template v-if="hasContent">
      <ChatPaneFrame
        v-for="tile in laidOutTiles"
        :key="tile.tileId"
        :tile-id="tile.tileId"
        :rect="tile.rect"
        :active="tile.tileId === activeTileId"
        :title="paneTitle(tile.conversationId)"
        :solo="laidOutTiles.length === 1"
        :dragging="tile.tileId === draggingTileId"
        @activate="activateTile(tile)"
      >
        <ChatWindow
          :conversation-id="tile.conversationId"
          :active="tile.tileId === activeTileId"
          embedded
          :rearrangeable="canRearrange"
          @close="closeTile(tile)"
          @header-drag="beginTilePointerDrag(tile, paneTitle(tile.conversationId), $event)"
        />
      </ChatPaneFrame>

      <ChatPaneFrame
        v-if="pendingPeer && laidOutTiles.length === 0"
        tile-id="pending"
        :rect="fullRect"
        :active="true"
        :title="pendingPeer.username"
        :show-chrome="false"
        @activate="() => undefined"
      >
        <ChatWindow :active="true" embedded />
      </ChatPaneFrame>

      <button
        v-for="sash in columnSashes"
        :key="`col-sash-${sash.index}`"
        type="button"
        class="chat-sash chat-sash--column"
        :style="{ left: `${sash.x}px`, top: `${sash.y}px`, height: `${sash.h}px` }"
        aria-label="Изменить ширину колонок"
        @pointerdown="beginColumnSashDrag(sash.index, $event)"
      />

      <button
        v-for="sash in rowSashes"
        :key="`row-sash-${sash.columnIndex}-${sash.tileIndex}`"
        type="button"
        class="chat-sash chat-sash--row"
        :style="{ left: `${sash.x}px`, top: `${sash.y}px`, width: `${sash.w}px` }"
        aria-label="Изменить высоту чатов"
        @pointerdown="beginRowSashDrag(sash.columnIndex, sash.tileIndex, $event)"
      />
    </template>

    <div v-else class="chat-placeholder">
      <div class="chat-placeholder__mark">M</div>
      <strong>Выберите чат</strong>
      <span>Откройте диалог слева или перетащите чат сюда</span>
    </div>

    <div
      v-if="dropPreview"
      class="chat-drop-preview"
      aria-hidden="true"
    >
      <div
        v-if="dropPreview.remainRect"
        class="chat-drop-preview__remain"
        :style="rectStyle(dropPreview.remainRect)"
      />
      <div
        class="chat-drop-preview__new"
        :style="rectStyle(dropPreview.newRect)"
      >
        <span>{{ dropLabel }}</span>
      </div>
    </div>

    <Teleport to="body">
      <div
        v-if="tileDragGhost"
        class="conversation-drag-ghost"
        :style="{
          transform: `translate3d(${tileDragGhost.clientX - 28}px, ${tileDragGhost.clientY - 28}px, 0)`,
        }"
      >
        <span>{{ tileDragGhost.title }}</span>
      </div>
    </Teleport>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, onScopeDispose, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useConversations } from '@/entities/conversation/api/conversation.queries'
import { getConversationTitle } from '@/entities/conversation/lib/display'
import { useChatWorkspaceStore } from '@/features/chat-workspace/model/chat-workspace.store'
import type { Rect } from '@/features/chat-workspace/model/types'
import { usePaneInteractions } from '@/features/chat-workspace/model/use-pane-interactions'
import { useTileDrag } from '@/features/chat-workspace/model/use-tile-drag'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import ChatWindow from '@/widgets/chat-window/ui/ChatWindow.vue'
import ChatPaneFrame from './ChatPaneFrame.vue'

const workspace = useChatWorkspaceStore()
const selection = useConversationSelectionStore()
const {
  laidOutTiles,
  activeTileId,
  dropTargetHot,
  dropPreview,
  columnSashes,
  rowSashes,
  draggingTileId,
  tileDragGhost,
} = storeToRefs(workspace)
const { pendingPeer, selectedConversationId } = storeToRefs(selection)
const { conversations } = useConversations()
const {
  activateTile,
  closeTile,
  beginColumnSashDrag,
  beginRowSashDrag,
} = usePaneInteractions()
const { beginTilePointerDrag } = useTileDrag()

const rootEl = ref<HTMLElement | null>(null)
let resizeObserver: ResizeObserver | null = null

const hasContent = computed(
  () => laidOutTiles.value.length > 0 || pendingPeer.value !== null,
)

const canRearrange = computed(() => laidOutTiles.value.length > 1)
const dropLabel = computed(() => {
  const preview = dropPreview.value
  const moving = Boolean(draggingTileId.value)
  if (!preview) return moving ? 'Переместить' : 'Новый чат'
  if (!preview.remainRect) return 'На всю область'
  const prefix = moving ? 'Сюда' : null
  switch (preview.zone) {
    case 'left':
      return prefix ? `${prefix} · слева` : 'Слева · 50%'
    case 'right':
      return prefix ? `${prefix} · справа` : 'Справа · 50%'
    case 'top':
      return prefix ? `${prefix} · сверху` : 'Сверху · 50%'
    case 'bottom':
      return prefix ? `${prefix} · снизу` : 'Снизу · 50%'
  }
  return moving ? 'Переместить' : 'Новый чат'
})

const fullRect = computed<Rect>(() => ({
  x: 0,
  y: 0,
  w: Math.max(workspace.workspaceBounds.width, 0),
  h: Math.max(workspace.workspaceBounds.height, 0),
}))

function paneTitle(conversationId: string): string {
  const conversation = conversations.value.find((item) => item.id === conversationId)
  return conversation ? getConversationTitle(conversation) : 'Чат'
}

function rectStyle(rect: Rect) {
  return {
    left: `${rect.x}px`,
    top: `${rect.y}px`,
    width: `${rect.w}px`,
    height: `${rect.h}px`,
  }
}

function publishBounds(): void {
  const el = rootEl.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  workspace.setWorkspaceBounds({
    width: rect.width,
    height: rect.height,
    left: rect.left,
    top: rect.top,
  })
}

function ensurePaneForSelection(): void {
  if (pendingPeer.value) return
  const selected = selectedConversationId.value
  if (!selected) return
  workspace.syncFromSelection(selected)
}

onMounted(() => {
  publishBounds()
  ensurePaneForSelection()
  resizeObserver = new ResizeObserver(() => publishBounds())
  if (rootEl.value) resizeObserver.observe(rootEl.value)
  window.addEventListener('resize', publishBounds)
})

onScopeDispose(() => {
  resizeObserver?.disconnect()
  window.removeEventListener('resize', publishBounds)
})

watch(selectedConversationId, () => ensurePaneForSelection())
</script>
