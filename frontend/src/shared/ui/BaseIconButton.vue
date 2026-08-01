<template>
  <button
    class="icon-button"
    :class="[`icon-button--${variant}`, { 'is-active': active }]"
    type="button"
    :aria-label="label"
    :title="label"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    label: string;
    active?: boolean;
    variant?: "default" | "primary" | "ghost";
  }>(),
  {
    active: false,
    variant: "default",
  },
);
</script>
