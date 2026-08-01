<template>
  <label class="text-field" :class="{ 'is-invalid': Boolean(error) }">
    <span v-if="label" class="text-field__label">{{ label }}</span>
    <input
      :value="modelValue"
      :type="type"
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      :inputmode="inputmode"
      :maxlength="maxlength"
      :disabled="disabled"
      :aria-invalid="Boolean(error)"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <span v-if="error" class="text-field__error">{{ error }}</span>
  </label>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: string
    label?: string
    placeholder?: string
    type?: string
    autocomplete?: string
    inputmode?: 'text' | 'tel' | 'numeric' | 'email'
    maxlength?: number
    disabled?: boolean
    error?: string
  }>(),
  {
    type: 'text',
    disabled: false,
  },
)

defineEmits<{ 'update:modelValue': [value: string] }>()
</script>
