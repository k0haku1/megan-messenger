<template>
  <Teleport to="body">
    <div
      v-if="open"
      ref="menuRef"
      class="context-menu"
      role="menu"
      @contextmenu.prevent
    >
      <button
        v-for="item in items"
        :key="item.id"
        type="button"
        class="context-menu__item"
        :class="{ 'is-danger': item.danger, 'is-disabled': item.disabled }"
        role="menuitem"
        :disabled="item.disabled"
        @click="select(item.id)"
      >
        {{ item.label }}
      </button>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { positionFixedMenu, type MenuPoint } from '@/shared/lib/position-fixed-menu'

export interface ContextMenuItem {
  id: string
  label: string
  danger?: boolean
  disabled?: boolean
}

const props = defineProps<{
  open: boolean
  point: MenuPoint | null
  items: ContextMenuItem[]
}>()

const emit = defineEmits<{ close: []; select: [actionId: string] }>()

const menuRef = ref<HTMLElement | null>(null)

watch(
  () => [props.open, props.point] as const,
  async ([isOpen, point]) => {
    if (!isOpen || !point) return
    await nextTick()
    if (menuRef.value) positionFixedMenu(menuRef.value, point)
  },
)

function select(actionId: string) {
  emit('select', actionId)
  emit('close')
}

function onPointerDown(event: MouseEvent) {
  if (!props.open) return
  if (menuRef.value?.contains(event.target as Node)) return
  emit('close')
}

function onKeyDown(event: KeyboardEvent) {
  if (props.open && event.key === 'Escape') emit('close')
}

function onScrollClose() {
  if (props.open) emit('close')
}

onMounted(() => {
  document.addEventListener('mousedown', onPointerDown)
  document.addEventListener('keydown', onKeyDown)
  window.addEventListener('scroll', onScrollClose, true)
})

onUnmounted(() => {
  document.removeEventListener('mousedown', onPointerDown)
  document.removeEventListener('keydown', onKeyDown)
  window.removeEventListener('scroll', onScrollClose, true)
})
</script>

<style scoped>
.context-menu {
  position: fixed;
  z-index: 1200;
  min-width: 200px;
  padding: 4px;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: 0 8px 24px rgb(15 23 42 / 18%);
}

.context-menu__item {
  display: block;
  width: 100%;
  padding: 8px 12px;
  text-align: left;
  font: inherit;
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
  background: transparent;
  border: 0;
  border-radius: var(--radius-sm);
  cursor: pointer;
}

.context-menu__item:hover:not(:disabled) {
  background: var(--color-bg-subtle);
}

.context-menu__item.is-danger {
  color: var(--color-danger);
}

.context-menu__item.is-disabled,
.context-menu__item:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
</style>
