<template>
  <form class="auth-form" @submit.prevent="flow.submitPhone()">
    <BaseTextField
      v-model="flow.phone"
      label="Телефон"
      placeholder="+79991234567"
      type="tel"
      inputmode="tel"
      autocomplete="tel"
      :error="flow.fieldErrors.phone"
      :disabled="flow.loading"
    />
    <p v-if="flow.formError" class="auth-form__error">{{ flow.formError }}</p>
    <BaseButton type="submit" block :loading="flow.loading">Получить код</BaseButton>
  </form>
</template>

<script setup lang="ts">
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'
import { useAuthFlowStore } from '../model/auth-flow.store'

const flow = useAuthFlowStore()
</script>
