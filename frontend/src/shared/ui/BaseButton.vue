<template>
  <button
    class="base-button"
    :class="[`base-button--${variant}`, { 'is-block': block, 'is-loading': loading }]"
    :type="type"
    :disabled="disabled || loading"
  >
    <span class="base-button__label"><slot /></span>
  </button>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'primary' | 'ghost' | 'danger'
    type?: 'button' | 'submit'
    block?: boolean
    loading?: boolean
    disabled?: boolean
  }>(),
  {
    variant: 'primary',
    type: 'button',
    block: false,
    loading: false,
    disabled: false,
  },
)
</script>
