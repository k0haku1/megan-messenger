<template>
  <label class="search-field">
    <span class="search-field__icon" aria-hidden="true">⌕</span>
    <input
      :value="modelValue"
      type="search"
      :placeholder="placeholder ?? 'Поиск'"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
  </label>
</template>

<script setup lang="ts">
defineProps<{ modelValue: string; placeholder?: string }>()
defineEmits<{ 'update:modelValue': [value: string] }>()
</script>
