<template>
  <form class="auth-form" @submit.prevent="flow.submitUsername()">
    <BaseTextField
      v-model="flow.username"
      label="Username"
      placeholder="sev"
      autocomplete="username"
      :error="flow.fieldErrors.username"
      :disabled="flow.loading"
    />
    <p class="auth-form__hint">Латиница и цифры, от 3 до 32 символов</p>
    <p v-if="flow.formError" class="auth-form__error">{{ flow.formError }}</p>
    <BaseButton type="submit" block :loading="flow.loading">Сохранить</BaseButton>
  </form>
</template>

<script setup lang="ts">
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'
import { useAuthFlowStore } from '../model/auth-flow.store'

const flow = useAuthFlowStore()
</script>
