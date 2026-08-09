import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { authApi } from '../api/auth.api'
import type { AuthResult, SessionUser } from './types'
import { getAccessToken, setAccessToken, subscribeAccessToken } from '@/shared/api/access-token'
import { clearLocalCache } from '@/shared/lib/clear-local-cache'

export const useSessionStore = defineStore('session', () => {
  const accessToken = ref<string | null>(getAccessToken())
  const user = ref<SessionUser | null>(null)
  const challengeToken = ref<string | null>(null)
  const pendingPhone = ref('')
  const bootstrapped = ref(false)

  subscribeAccessToken((token) => {
    accessToken.value = token
    if (!token) user.value = null
  })

  const isAuthenticated = computed(() => Boolean(accessToken.value))
  const needsUsername = computed(() => Boolean(accessToken.value && user.value && !user.value.onboardingComplete))
  const needsPassword = computed(() => Boolean(challengeToken.value))
  const isReady = computed(() => bootstrapped.value)

  function persistToken(token: string | null) {
    setAccessToken(token)
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
    await clearLocalCache()
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
