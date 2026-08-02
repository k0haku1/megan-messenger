<template>
  <div class="profile-view">
    <section class="profile-view__section">
      <p class="profile-modal__hint">
        Дополнительная защита после SMS. Если пароль задан, при входе его нужно ввести после кода.
      </p>

      <p v-if="session.user?.hasPassword" class="profile-modal__status is-on">Пароль установлен</p>
      <p v-else class="profile-modal__status">Пароль не задан</p>

      <form class="auth-form" @submit.prevent="onSavePassword">
        <BaseTextField
          v-if="session.user?.hasPassword"
          v-model="currentPassword"
          label="Текущий пароль"
          type="password"
          autocomplete="current-password"
          :error="fieldErrors.currentPassword"
          :disabled="loading"
        />
        <BaseTextField
          v-model="password"
          :label="session.user?.hasPassword ? 'Новый пароль' : 'Пароль'"
          type="password"
          autocomplete="new-password"
          :error="fieldErrors.password"
          :disabled="loading"
        />
        <p v-if="formError" class="auth-form__error">{{ formError }}</p>
        <p v-if="successMessage" class="profile-modal__success">{{ successMessage }}</p>
        <BaseButton type="submit" block :loading="loading">
          {{ session.user?.hasPassword ? 'Сменить пароль' : 'Установить пароль' }}
        </BaseButton>
      </form>

      <form
        v-if="session.user?.hasPassword"
        class="auth-form profile-modal__remove"
        @submit.prevent="onRemovePassword"
      >
        <BaseTextField
          v-model="removePasswordValue"
          label="Текущий пароль для удаления"
          type="password"
          autocomplete="current-password"
          :error="removeFieldErrors.currentPassword"
          :disabled="loading"
        />
        <BaseButton type="submit" variant="danger" block :loading="loading">
          Удалить пароль
        </BaseButton>
      </form>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'
import { ApiError } from '@/shared/api/http'
import { authApi } from '@/entities/session/api/auth.api'
import { useSessionStore } from '@/entities/session/model/session.store'
import { useProfileUiStore } from '../model/profile-ui.store'

const session = useSessionStore()
const profileUi = useProfileUiStore()

const password = ref('')
const currentPassword = ref('')
const removePasswordValue = ref('')
const loading = ref(false)
const formError = ref('')
const successMessage = ref('')
const fieldErrors = ref<Record<string, string>>({})
const removeFieldErrors = ref<Record<string, string>>({})

function resetForm() {
  password.value = ''
  currentPassword.value = ''
  removePasswordValue.value = ''
  formError.value = ''
  successMessage.value = ''
  fieldErrors.value = {}
  removeFieldErrors.value = {}
}

watch(
  () => [profileUi.isOpen, profileUi.view] as const,
  ([open, view]) => {
    if (open && view === 'password') resetForm()
  },
  { immediate: true },
)

function captureError(error: unknown, target: 'save' | 'remove') {
  if (error instanceof ApiError) {
    if (error.fields) {
      if (target === 'save') fieldErrors.value = error.fields
      else removeFieldErrors.value = error.fields
      formError.value = ''
      return
    }
    formError.value = error.message || 'Не удалось сохранить'
    return
  }
  formError.value = 'Не удалось сохранить'
}

async function onSavePassword() {
  formError.value = ''
  successMessage.value = ''
  fieldErrors.value = {}
  loading.value = true
  try {
    if (password.value.trim().length < 4) {
      fieldErrors.value = { password: 'Минимум 4 символа' }
      return
    }
    if (session.user?.hasPassword) {
      await authApi.setPassword(password.value, currentPassword.value)
    } else {
      await authApi.setPassword(password.value)
    }
    await session.refreshProfile()
    password.value = ''
    currentPassword.value = ''
    successMessage.value = 'Пароль сохранён. При следующем входе после SMS потребуется пароль.'
  } catch (error) {
    captureError(error, 'save')
  } finally {
    loading.value = false
  }
}

async function onRemovePassword() {
  formError.value = ''
  successMessage.value = ''
  removeFieldErrors.value = {}
  loading.value = true
  try {
    await authApi.removePassword(removePasswordValue.value)
    await session.refreshProfile()
    removePasswordValue.value = ''
    successMessage.value = 'Облачный пароль удалён'
  } catch (error) {
    captureError(error, 'remove')
  } finally {
    loading.value = false
  }
}
</script>
