<template>
  <svg
    class="app-icon"
    :class="sizeClass"
    viewBox="0 0 24 24"
    fill="none"
    aria-hidden="true"
  >
    <path
      v-if="name === 'chevron-left'"
      d="M15 6l-6 6 6 6"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
    <path
      v-else-if="name === 'close'"
      d="M6 6l12 12M18 6L6 18"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
    />
    <path
      v-else-if="name === 'edit'"
      d="M4 20h4l10.5-10.5a2.1 2.1 0 0 0-3-3L5 17v3"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
  </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    name: 'chevron-left' | 'close' | 'check' | 'edit'
    size?: 'sm' | 'md'
  }>(),
  {
    size: 'md',
  },
)

const sizeClass = computed(() => `app-icon--${props.size}`)
</script>
