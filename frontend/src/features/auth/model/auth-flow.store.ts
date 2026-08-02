import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { ApiError } from '@/shared/api/http'
import { authApi } from '@/entities/session/api/auth.api'
import { normalizeUsernameQuery, USERNAME_HINT } from '@/entities/user/lib/username'
import { useSessionStore } from '@/entities/session/model/session.store'
import { router } from '@/app/router'

export type AuthStep = 'phone' | 'otp' | 'password' | 'username'

function normalizePhoneInput(value: string): string {
  const trimmed = value.trim()
  if (!trimmed) return ''
  const digits = trimmed.replace(/[^\d+]/g, '')
  if (digits.startsWith('+')) return digits
  if (digits.startsWith('8') && digits.length === 11) return `+7${digits.slice(1)}`
  if (digits.startsWith('7') && digits.length === 11) return `+${digits}`
  return digits.startsWith('+') ? digits : `+${digits}`
}

function resolvePostAuthPath(): string {
  const redirect = router.currentRoute.value.query.redirect
  if (typeof redirect === 'string' && redirect.startsWith('/') && redirect !== '/auth') {
    return redirect
  }
  return '/'
}

export const useAuthFlowStore = defineStore('auth-flow', () => {
  const session = useSessionStore()
  const step = ref<AuthStep>('phone')
  const phone = ref('')
  const code = ref('')
  const password = ref('')
  const username = ref('')
  const loading = ref(false)
  const formError = ref('')
  const fieldErrors = ref<Record<string, string>>({})

  const title = computed(() => {
    switch (step.value) {
      case 'phone':
        return 'Вход в Megan'
      case 'otp':
        return 'Код из SMS'
      case 'password':
        return 'Облачный пароль'
      case 'username':
        return 'Выберите username'
    }
  })

  const subtitle = computed(() => {
    switch (step.value) {
      case 'phone':
        return 'Введите номер телефона в формате +79991234567'
      case 'otp':
        return `Мы отправили код на ${session.pendingPhone || phone.value}`
      case 'password':
        return 'Для этого аккаунта включена дополнительная защита'
      case 'username':
        return USERNAME_HINT
    }
  })

  function resetErrors() {
    formError.value = ''
    fieldErrors.value = {}
  }

  function captureError(error: unknown) {
    if (error instanceof ApiError) {
      if (error.fields) {
        fieldErrors.value = error.fields
        formError.value = ''
        return
      }
      formError.value = error.message || 'Не удалось выполнить запрос'
      return
    }
    formError.value = 'Не удалось выполнить запрос'
  }

  async function enterMessenger() {
    await router.replace(resolvePostAuthPath())
  }

  function syncStepFromSession() {
    if (session.needsPassword) {
      step.value = 'password'
      return
    }
    if (session.needsUsername) {
      step.value = 'username'
      return
    }
    if (session.pendingPhone && !session.isAuthenticated) {
      step.value = 'otp'
      return
    }
    step.value = 'phone'
  }

  async function submitPhone() {
    resetErrors()
    loading.value = true
    try {
      const normalized = normalizePhoneInput(phone.value)
      if (!/^\+[1-9]\d{7,14}$/.test(normalized)) {
        fieldErrors.value = { phone: 'Введите номер в формате +79991234567' }
        return
      }
      await authApi.startPhone(normalized)
      phone.value = normalized
      session.setPendingPhone(normalized)
      code.value = ''
      step.value = 'otp'
    } catch (error) {
      captureError(error)
    } finally {
      loading.value = false
    }
  }

  async function submitOtp() {
    resetErrors()
    loading.value = true
    try {
      const result = await authApi.verifyPhone(session.pendingPhone || phone.value, code.value.trim())
      session.applyAuthResult(result)
      if (result.needPassword) {
        step.value = 'password'
        return
      }
      await session.refreshProfile()
      if (session.needsUsername) {
        step.value = 'username'
        return
      }
      await enterMessenger()
    } catch (error) {
      captureError(error)
    } finally {
      loading.value = false
    }
  }

  async function submitPassword() {
    resetErrors()
    loading.value = true
    try {
      if (!session.challengeToken) {
        formError.value = 'Сессия проверки пароля истекла. Запросите код снова.'
        step.value = 'phone'
        return
      }
      const result = await authApi.verifyPassword(session.challengeToken, password.value)
      session.applyAuthResult(result)
      await session.refreshProfile()
      if (session.needsUsername) {
        step.value = 'username'
        return
      }
      await enterMessenger()
    } catch (error) {
      captureError(error)
    } finally {
      loading.value = false
    }
  }

  async function submitUsername() {
    resetErrors()
    loading.value = true
    try {
      const result = await authApi.completeUsername(normalizeUsernameQuery(username.value))
      session.applyAuthResult(result)
      if (result.user) session.setUser(result.user)
      else await session.refreshProfile()
      await enterMessenger()
    } catch (error) {
      captureError(error)
    } finally {
      loading.value = false
    }
  }

  async function resendCode() {
    resetErrors()
    loading.value = true
    try {
      const target = session.pendingPhone || phone.value
      await authApi.startPhone(target)
      formError.value = ''
    } catch (error) {
      captureError(error)
    } finally {
      loading.value = false
    }
  }

  function goBack() {
    resetErrors()
    if (step.value === 'otp') {
      step.value = 'phone'
      code.value = ''
      return
    }
    if (step.value === 'password') {
      session.clearChallenge()
      step.value = 'otp'
      password.value = ''
    }
  }

  return {
    step,
    phone,
    code,
    password,
    username,
    loading,
    formError,
    fieldErrors,
    title,
    subtitle,
    submitPhone,
    submitOtp,
    submitPassword,
    submitUsername,
    resendCode,
    goBack,
    syncStepFromSession,
  }
})
