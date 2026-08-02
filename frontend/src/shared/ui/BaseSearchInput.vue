<template>
  <label class="search-field">
    <span class="search-field__icon" aria-hidden="true">⌕</span>
    <input
      ref="inputRef"
      :value="modelValue"
      type="search"
      :placeholder="placeholder ?? 'Поиск'"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
  </label>
</template>

<script setup lang="ts">
import { ref } from 'vue'

defineProps<{ modelValue: string; placeholder?: string }>()
defineEmits<{ 'update:modelValue': [value: string] }>()

const inputRef = ref<HTMLInputElement | null>(null)

function focus() {
  inputRef.value?.focus()
}

defineExpose({ focus })
</script>
