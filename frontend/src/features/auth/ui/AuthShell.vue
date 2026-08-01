<template>
  <section class="auth-shell">
    <div class="auth-shell__panel">
      <div class="auth-shell__brand" aria-hidden="true">M</div>
      <header class="auth-shell__header">
        <h1>{{ title }}</h1>
        <p>{{ subtitle }}</p>
      </header>
      <slot />
    </div>
  </section>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  subtitle: string
}>()
</script>
