<template>
  <button
    class="icon-button"
    :class="[`icon-button--${variant}`, { 'is-active': active }]"
    type="button"
    :aria-label="label"
    :title="label"
    :disabled="disabled"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    label: string;
    active?: boolean;
    disabled?: boolean;
    variant?: "default" | "primary" | "ghost";
  }>(),
  {
    active: false,
    disabled: false,
    variant: "default",
  },
);
</script>
