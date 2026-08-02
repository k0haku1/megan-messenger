<template>
  <div class="profile-view">
    <p class="modal-sheet__hint">{{ USERNAME_HINT }}</p>

    <form class="auth-form" @submit.prevent="submit()">
      <BaseTextField
        v-model="username"
        label="Username"
        placeholder="megan_user"
        autocomplete="username"
        :error="fieldError"
      />
      <p v-if="success" class="modal-sheet__success">{{ success }}</p>
      <BaseButton type="submit" class="is-block" :disabled="isSaving">
        {{ isSaving ? 'Сохраняем…' : 'Сохранить username' }}
      </BaseButton>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { userApi } from '@/entities/user/api/user.api'
import { USERNAME_HINT } from '@/entities/user/lib/username'
import { useSessionStore } from '@/entities/session/model/session.store'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'
import { ApiError } from '@/shared/api/http'

const session = useSessionStore()
const username = ref(session.user?.username ?? '')
const fieldError = ref<string | undefined>()
const success = ref('')
const isSaving = ref(false)

async function submit() {
  fieldError.value = undefined
  success.value = ''
  isSaving.value = true

  try {
    const result = await userApi.updateUsername(username.value.trim())
    if (session.user) {
      session.user.username = result.username
    }
    success.value = `Username обновлён: @${result.username}`
  } catch (error) {
    if (error instanceof ApiError && error.fields?.username) {
      fieldError.value = error.fields.username
    } else {
      fieldError.value = error instanceof ApiError ? error.message : 'Не удалось сохранить username'
    }
  } finally {
    isSaving.value = false
  }
}
</script>
