<template>
  <form class="auth-form" @submit.prevent="flow.submitUsername()">
    <BaseTextField
      v-model="flow.username"
      label="Username"
      placeholder="megan_user"
      autocomplete="username"
      :error="flow.fieldErrors.username"
      :disabled="flow.loading"
    />
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
