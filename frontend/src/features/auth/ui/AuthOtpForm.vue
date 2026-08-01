<template>
  <form class="auth-form" @submit.prevent="flow.submitOtp()">
    <BaseTextField
      v-model="flow.code"
      label="Код из SMS"
      placeholder="123456"
      inputmode="numeric"
      autocomplete="one-time-code"
      :maxlength="6"
      :error="flow.fieldErrors.code"
      :disabled="flow.loading"
    />
    <p v-if="flow.formError" class="auth-form__error">{{ flow.formError }}</p>
    <BaseButton type="submit" block :loading="flow.loading">Войти</BaseButton>
    <div class="auth-form__row">
      <BaseButton variant="ghost" :disabled="flow.loading" @click="flow.goBack()">Назад</BaseButton>
      <BaseButton variant="ghost" :disabled="flow.loading" @click="flow.resendCode()">Отправить снова</BaseButton>
    </div>
  </form>
</template>

<script setup lang="ts">
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'
import { useAuthFlowStore } from '../model/auth-flow.store'

const flow = useAuthFlowStore()
</script>
