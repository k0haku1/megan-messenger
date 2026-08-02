<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="modal-sheet"
      role="dialog"
      aria-modal="true"
      :aria-labelledby="titleId"
      @click.self="emit('close')"
    >
      <div class="modal-sheet__panel" :class="panelClass">
        <header class="modal-sheet__header">
          <div class="modal-sheet__header-start">
            <slot name="header-start" />
            <h2 :id="titleId" class="modal-sheet__title">{{ title }}</h2>
          </div>
          <BaseIconButton label="Закрыть" variant="ghost" @click="emit('close')">
            <AppIcon name="close" />
          </BaseIconButton>
        </header>

        <slot />
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'

const props = withDefaults(
  defineProps<{
    open: boolean
    title: string
    titleId: string
    scrollable?: boolean
    compact?: boolean
  }>(),
  {
    scrollable: false,
    compact: false,
  },
)

const emit = defineEmits<{ close: [] }>()

const panelClass = computed(() => ({
  'modal-sheet__panel--scroll': props.scrollable,
  'modal-sheet__panel--compact': props.compact,
}))
</script>
