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
      v-else-if="name === 'check'"
      d="M5 12l4 4L19 6"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
    <path
      v-else-if="name === 'edit'"
      d="M4 20h4l10.5-10.5a2.1 2.1 0 0 0-3-3L5 17v3"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
    <g v-else-if="name === 'username'">
      <circle cx="12" cy="12" r="4" stroke="currentColor" stroke-width="2" />
      <path
        d="M16 8v4a3 3 0 0 1-6 0V8"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
      />
    </g>
    <path
      v-else-if="name === 'privacy'"
      d="M12 3l8 4v6c0 5.25-3.5 8.5-8 9.5-4.5-1-8-4.5-8-9.5V7l8-4z"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
    <g v-else-if="name === 'theme'">
      <path
        d="M12 3a9 9 0 1 0 0 18V3z"
        stroke="currentColor"
        stroke-width="2"
        stroke-linejoin="round"
      />
      <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2" />
    </g>
    <g v-else-if="name === 'password'">
      <rect
        x="5"
        y="11"
        width="14"
        height="10"
        rx="2"
        stroke="currentColor"
        stroke-width="2"
      />
      <path
        d="M8 11V8a4 4 0 0 1 8 0v3"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
      />
    </g>
  </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { AppIconName } from '@/shared/ui/app-icon.types'

const props = withDefaults(
  defineProps<{
    name: AppIconName
    size?: 'sm' | 'md'
  }>(),
  {
    size: 'md',
  },
)

const sizeClass = computed(() => `app-icon--${props.size}`)
</script>
