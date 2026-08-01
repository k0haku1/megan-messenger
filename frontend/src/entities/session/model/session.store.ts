import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { authApi } from '../api/auth.api'
import type { AuthResult, SessionUser } from './types'

const ACCESS_TOKEN_KEY = 'mm_access_token'

function readStoredToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY)
}

export const useSessionStore = defineStore('session', () => {
  const accessToken = ref<string | null>(readStoredToken())
  const user = ref<SessionUser | null>(null)
  const challengeToken = ref<string | null>(null)
  const pendingPhone = ref('')
  const bootstrapped = ref(false)

  const isAuthenticated = computed(() => Boolean(accessToken.value))
  const needsUsername = computed(() => Boolean(accessToken.value && user.value && !user.value.onboardingComplete))
  const needsPassword = computed(() => Boolean(challengeToken.value))
  const isReady = computed(() => bootstrapped.value)

  function persistToken(token: string | null) {
    accessToken.value = token
    if (token) localStorage.setItem(ACCESS_TOKEN_KEY, token)
    else localStorage.removeItem(ACCESS_TOKEN_KEY)
  }

  function applyAuthResult(result: AuthResult) {
    if (result.needPassword && result.challengeToken) {
      challengeToken.value = result.challengeToken
      persistToken(null)
      user.value = null
      return
    }

    challengeToken.value = null
    if (result.accessToken) persistToken(result.accessToken)
  }

  async function bootstrap() {
    if (!accessToken.value) {
      bootstrapped.value = true
      return
    }

    try {
      user.value = await authApi.getMe()
    } catch {
      persistToken(null)
      user.value = null
    } finally {
      bootstrapped.value = true
    }
  }

  async function refreshProfile() {
    if (!accessToken.value) {
      user.value = null
      return
    }
    user.value = await authApi.getMe()
  }

  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // ignore network errors on logout
    }
    persistToken(null)
    user.value = null
    challengeToken.value = null
    pendingPhone.value = ''
  }

  function setPendingPhone(phone: string) {
    pendingPhone.value = phone
  }

  function clearChallenge() {
    challengeToken.value = null
  }

  function setUser(next: SessionUser) {
    user.value = next
  }

  return {
    accessToken,
    user,
    challengeToken,
    pendingPhone,
    bootstrapped,
    isAuthenticated,
    needsUsername,
    needsPassword,
    isReady,
    bootstrap,
    refreshProfile,
    applyAuthResult,
    persistToken,
    logout,
    setPendingPhone,
    clearChallenge,
    setUser,
  }
})
