<template>
  <span class="folder-badge" :class="`folder-badge--${size}`">{{ initial }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { getFolderInitial } from '@/entities/folder/lib/display'

const props = withDefaults(
  defineProps<{
    name: string
    size?: 'sm' | 'md'
  }>(),
  { size: 'sm' },
)

const initial = computed(() => getFolderInitial(props.name))
</script>

<style scoped>
.folder-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border-radius: 8px;
  background: var(--color-primary);
  color: #fff;
  font-weight: 700;
  line-height: 1;
}

.folder-badge--sm {
  width: 22px;
  height: 22px;
  font-size: 0.75rem;
}

.folder-badge--md {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  font-size: 0.875rem;
}
</style>
