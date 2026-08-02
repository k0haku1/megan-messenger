<template>
  <label class="search-field">
    <AppIcon class="search-field__icon" name="search" size="sm" />
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
import AppIcon from '@/shared/ui/AppIcon.vue'

defineProps<{ modelValue: string; placeholder?: string }>()
defineEmits<{ 'update:modelValue': [value: string] }>()

const inputRef = ref<HTMLInputElement | null>(null)

function focus() {
  inputRef.value?.focus()
}

defineExpose({ focus })
</script>
