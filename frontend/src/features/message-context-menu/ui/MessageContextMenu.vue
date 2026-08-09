<template>
  <Teleport to="body">
    <div
      v-if="open"
      ref="menuRef"
      class="message-context-menu"
      role="menu"
      @contextmenu.prevent
    >
      <button
        v-for="item in items"
        :key="item.id"
        type="button"
        class="message-context-menu__item"
        :class="{ 'is-danger': item.danger, 'is-disabled': item.disabled }"
        role="menuitem"
        :disabled="item.disabled"
        @click="select(item.id)"
      >
        {{ item.label }}
      </button>
      <div v-if="reactionDetails.length" class="message-context-menu__divider" />
      <div v-if="reactionDetails.length" class="message-context-menu__reactions" aria-label="Реакции">
        <div v-for="reaction in reactionDetails" :key="reaction.emoji" class="message-context-menu__reaction">
          <span>{{ reaction.emoji }}</span>
          <span>{{ reaction.users.join(', ') }}</span>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

export interface MessageContextMenuItem {
  id: string
  label: string
  danger?: boolean
  disabled?: boolean
}

export interface MessageContextMenuAnchor {
  top: number
  left: number
  right: number
  bottom: number
  align: 'start' | 'end'
}

export interface MessageContextMenuReactionDetail {
  emoji: string
  users: string[]
}

const props = defineProps<{
  open: boolean
  anchor: MessageContextMenuAnchor | null
  items: MessageContextMenuItem[]
  reactionDetails?: MessageContextMenuReactionDetail[]
}>()

const emit = defineEmits<{ close: []; select: [actionId: string] }>()

const menuRef = ref<HTMLElement | null>(null)
const reactionDetails = computed(() => props.reactionDetails ?? [])

watch(
  () => [props.open, props.anchor] as const,
  async ([isOpen]) => {
    if (!isOpen || !props.anchor) return
    await nextTick()
    positionFromAnchor()
  },
)

function positionFromAnchor() {
  const menu = menuRef.value
  const anchor = props.anchor
  if (!menu || !anchor) return

  const padding = 8
  const gap = 6
  const rect = menu.getBoundingClientRect()

  let x = anchor.align === 'end' ? anchor.right - rect.width : anchor.left
  let y = anchor.bottom + gap

  if (y + rect.height > window.innerHeight - padding) {
    y = anchor.top - rect.height - gap
  }

  if (x + rect.width > window.innerWidth - padding) {
    x = window.innerWidth - rect.width - padding
  }
  if (x < padding) {
    x = padding
  }
  if (y < padding) {
    y = padding
  }

  menu.style.left = `${x}px`
  menu.style.top = `${y}px`
}

function select(actionId: string) {
  emit('select', actionId)
  emit('close')
}

function onPointerDown(event: MouseEvent) {
  if (!props.open) return
  const menu = menuRef.value
  if (menu?.contains(event.target as Node)) return
  emit('close')
}

function onKeyDown(event: KeyboardEvent) {
  if (props.open && event.key === 'Escape') {
    emit('close')
  }
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

function onScrollClose() {
  if (props.open) emit('close')
}
</script>

<style scoped>
.message-context-menu {
  position: fixed;
  z-index: 1200;
  min-width: 200px;
  padding: 4px;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: 0 8px 24px rgb(15 23 42 / 18%);
}
.message-context-menu__item {
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
.message-context-menu__item:hover:not(:disabled) {
  background: var(--color-bg-subtle);
}
.message-context-menu__item.is-danger {
  color: var(--color-danger);
}
.message-context-menu__item.is-disabled,
.message-context-menu__item:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
.message-context-menu__divider { height: 1px; margin: 4px 6px; background: var(--color-border); }
.message-context-menu__reactions { padding: 4px 8px 6px; }
.message-context-menu__reaction { display: grid; grid-template-columns: 24px minmax(0, 1fr); gap: 6px; padding: 3px 0; color: var(--color-text-secondary); font-size: var(--font-size-sm); }
.message-context-menu__reaction span:last-child { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
