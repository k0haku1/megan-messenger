<template>
  <button class="profile-nav-row" type="button" @click="$emit('select')">
    <span v-if="icon" class="profile-nav-row__icon">
      <AppIcon :name="icon" size="sm" />
    </span>
    <span class="profile-nav-row__body">
      <strong>{{ label }}</strong>
      <small v-if="description">{{ description }}</small>
    </span>
    <span v-if="value" class="profile-nav-row__value">{{ value }}</span>
    <span class="profile-nav-row__chevron" aria-hidden="true">
      <AppIcon name="chevron-right" size="sm" />
    </span>
  </button>
</template>

<script setup lang="ts">
import AppIcon from '@/shared/ui/AppIcon.vue'
import type { AppIconName } from '@/shared/ui/app-icon.types'

defineProps<{
  icon?: AppIconName
  label: string
  description?: string
  value?: string
}>()

defineEmits<{ select: [] }>()
</script>
