<template>
  <span
    class="avatar"
    :class="[`avatar--${size}`, `avatar--color-${color % 6}`]"
    aria-hidden="true"
  >
    <img v-if="src" :src="src" :alt="name" />
    <span v-else>{{ initials }}</span>
  </span>
</template>

<script setup lang="ts">
import { computed } from "vue";

const props = withDefaults(
  defineProps<{
    name: string;
    src?: string;
    size?: "sm" | "md" | "lg" | "xl";
    color?: number;
  }>(),
  {
    size: "md",
    color: 0,
  },
);

const initials = computed(() =>
  props.name
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part[0])
    .join("")
    .toUpperCase(),
);
</script>
