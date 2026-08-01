<template>
  <form class="auth-form" @submit.prevent="flow.submitPassword()">
    <BaseTextField
      v-model="flow.password"
      label="Пароль"
      type="password"
      placeholder="••••"
      autocomplete="current-password"
      :error="flow.fieldErrors.password ?? flow.fieldErrors.currentPassword"
      :disabled="flow.loading"
    />
    <p v-if="flow.formError" class="auth-form__error">{{ flow.formError }}</p>
    <BaseButton type="submit" block :loading="flow.loading">Продолжить</BaseButton>
    <BaseButton variant="ghost" block :disabled="flow.loading" @click="flow.goBack()">Назад</BaseButton>
  </form>
</template>

<script setup lang="ts">
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'
import { useAuthFlowStore } from '../model/auth-flow.store'

const flow = useAuthFlowStore()
</script>
