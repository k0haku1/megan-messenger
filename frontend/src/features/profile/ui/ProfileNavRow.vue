<template>
  <button class="profile-nav-row" type="button" @click="$emit('select')">
    <span v-if="icon" class="profile-nav-row__icon" aria-hidden="true">{{ icon }}</span>
    <span class="profile-nav-row__body">
      <strong>{{ label }}</strong>
      <small v-if="description">{{ description }}</small>
    </span>
    <span v-if="value" class="profile-nav-row__value">{{ value }}</span>
    <span class="profile-nav-row__chevron" aria-hidden="true">›</span>
  </button>
</template>

<script setup lang="ts">
defineProps<{
  icon?: string
  label: string
  description?: string
  value?: string
}>()

defineEmits<{ select: [] }>()
</script>
