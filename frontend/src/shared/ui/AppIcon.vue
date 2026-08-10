<template>
  <svg
    class="app-icon"
    :class="[sizeClass, { 'app-icon--send': name === 'send' }]"
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
      v-else-if="name === 'chevron-right'"
      d="M9 6l6 6-6 6"
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
      d="M5.5 12.5l3.2 3.2L18 6.5"
      stroke="currentColor"
      stroke-width="2.2"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
    <g
      v-else-if="name === 'checks'"
      stroke="currentColor"
      stroke-width="2.2"
      stroke-linecap="round"
      stroke-linejoin="round"
      fill="none"
    >
      <path d="M2.8 12.6l3 3L12.2 8" />
      <path d="M9.2 13.2l3.1 3.1L21.2 7.2" />
    </g>
    <path
      v-else-if="name === 'edit'"
      d="M4 20h4l10.5-10.5a2.1 2.1 0 0 0-3-3L5 17v3"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
    <g v-else-if="name === 'search'">
      <circle cx="11" cy="11" r="6" stroke="currentColor" stroke-width="2" />
      <path d="M16 16l4 4" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
    </g>
    <g v-else-if="name === 'menu'">
      <circle cx="12" cy="6" r="1.5" fill="currentColor" />
      <circle cx="12" cy="12" r="1.5" fill="currentColor" />
      <circle cx="12" cy="18" r="1.5" fill="currentColor" />
    </g>
    <g v-else-if="name === 'emoji'">
      <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2" />
      <path d="M8.5 10h.01M15.5 10h.01" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
      <path
        d="M8.5 14a4 4 0 0 0 7 0"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
      />
    </g>
    <path
      v-else-if="name === 'send'"
      d="M5 12l14-7-7 14-2-5-5-2z"
      stroke="currentColor"
      stroke-width="2"
      stroke-linejoin="round"
    />
    <path
      v-else-if="name === 'attach'"
      d="M8 12l7.5-7.5a3 3 0 0 1 4.24 4.24L9.5 19a4 4 0 0 1-5.66-5.66L14 3"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
    <path
      v-else-if="name === 'plus'"
      d="M12 5v14M5 12h14"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
    />
    <g v-else-if="name === 'user'">
      <circle cx="12" cy="8" r="4" stroke="currentColor" stroke-width="2" />
      <path
        d="M5 20c1.5-3.5 4.5-5 7-5s5.5 1.5 7 5"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
      />
    </g>
    <g v-else-if="name === 'logout'">
      <path
        d="M10 7V5a2 2 0 0 1 2-2h7v18h-7a2 2 0 0 1-2-2v-2"
        stroke="currentColor"
        stroke-width="2"
        stroke-linejoin="round"
      />
      <path
        d="M14 12H4m0 0l3-3M4 12l3 3"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </g>
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
      <rect x="5" y="11" width="14" height="10" rx="2" stroke="currentColor" stroke-width="2" />
      <path
        d="M8 11V8a4 4 0 0 1 8 0v3"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
      />
    </g>
    <g v-else-if="name === 'folder-all'">
      <rect x="4" y="6" width="16" height="13" rx="2" stroke="currentColor" stroke-width="2" />
      <path d="M4 10h16" stroke="currentColor" stroke-width="2" />
    </g>
    <g v-else-if="name === 'folder-personal'">
      <circle cx="12" cy="10" r="3" stroke="currentColor" stroke-width="2" />
      <path
        d="M7 18c.8-2.2 2.6-3.5 5-3.5s4.2 1.3 5 3.5"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
      />
    </g>
    <g v-else-if="name === 'folder-groups'">
      <circle cx="9" cy="10" r="2.5" stroke="currentColor" stroke-width="2" />
      <circle cx="15" cy="10" r="2.5" stroke="currentColor" stroke-width="2" />
      <path
        d="M5 18c.6-2 2-3 4-3m6 3c.6-2 2-3 4-3"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
      />
    </g>
    <g v-else-if="name === 'folder-work'">
      <rect x="4" y="8" width="16" height="11" rx="2" stroke="currentColor" stroke-width="2" />
      <path d="M9 8V6a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2" stroke="currentColor" stroke-width="2" />
    </g>
    <g v-else-if="name === 'mail'">
      <rect x="4" y="6" width="16" height="12" rx="2" stroke="currentColor" stroke-width="2" />
      <path d="M4 8l8 5 8-5" stroke="currentColor" stroke-width="2" stroke-linejoin="round" />
    </g>
    <g v-else-if="name === 'book'">
      <path
        d="M6 5h11a2 2 0 0 1 2 2v13H8a2 2 0 0 0-2 2V5z"
        stroke="currentColor"
        stroke-width="2"
        stroke-linejoin="round"
      />
      <path d="M8 5v15" stroke="currentColor" stroke-width="2" />
    </g>
    <path
      v-else-if="name === 'link'"
      d="M10 14a4 4 0 0 0 5.66 0l2.34-2.34a4 4 0 0 0-5.66-5.66L11 7.34M14 10a4 4 0 0 0-5.66 0L6 12.34a4 4 0 0 0 5.66 5.66L13 16.66"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
    />
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
