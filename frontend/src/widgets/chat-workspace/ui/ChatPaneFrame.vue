<template>
  <article
    class="chat-pane"
    :class="{
      'is-active': active,
      'is-entering': entering,
      'is-solo': solo,
      'is-dragging-source': dragging,
    }"
    :style="paneStyle"
    :aria-label="title"
    @pointerdown="onPanePointerDown"
  >
    <div class="chat-pane__body">
      <slot />
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import type { Rect } from '@/features/chat-workspace/model/types'
import { useChatWorkspaceStore } from '@/features/chat-workspace/model/chat-workspace.store'

const props = withDefaults(
  defineProps<{
    tileId: string
    rect: Rect
    active: boolean
    title: string
    solo?: boolean
    dragging?: boolean
  }>(),
  {
    solo: false,
    dragging: false,
  },
)

const emit = defineEmits<{
  activate: []
}>()

const workspace = useChatWorkspaceStore()
const entering = ref(false)
const enterFrom = ref<Rect | null>(null)

const paneStyle = computed(() => {
  const from = enterFrom.value
  if (entering.value && from) {
    return {
      left: `${from.x}px`,
      top: `${from.y}px`,
      width: `${from.w}px`,
      height: `${from.h}px`,
      opacity: 0.75,
    }
  }

  return {
    left: `${props.rect.x}px`,
    top: `${props.rect.y}px`,
    width: `${props.rect.w}px`,
    height: `${props.rect.h}px`,
  }
})

function onPanePointerDown(event: PointerEvent): void {
  const target = event.target as HTMLElement | null
  // Header drag is handled inside ChatWindow; ignore here.
  if (target?.closest('[data-tile-drag-handle]')) return
  emit('activate')
}

function armEnterAnimation(): void {
  const from = workspace.enterFromByTileId[props.tileId]
  if (!from) return
  enterFrom.value = from
  entering.value = true
  workspace.clearEnterFrom(props.tileId)

  void nextTick(() => {
    requestAnimationFrame(() => {
      enterFrom.value = null
      window.setTimeout(() => {
        entering.value = false
      }, 200)
    })
  })
}

onMounted(armEnterAnimation)
watch(
  () => workspace.enterFromByTileId[props.tileId],
  (value) => {
    if (value) armEnterAnimation()
  },
)
</script>
