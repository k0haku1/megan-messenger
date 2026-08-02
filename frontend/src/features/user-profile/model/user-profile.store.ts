import { ref } from 'vue'
import { defineStore } from 'pinia'
import { userApi } from '@/entities/user/api/user.api'
import { normalizeUsernameQuery } from '@/entities/user/lib/username'
import type { PublicUser } from '@/entities/user/model/types'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'

export const useUserProfileStore = defineStore('user-profile', () => {
  const isOpen = ref(false)
  const username = ref('')
  const profile = ref<PublicUser | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  async function open(rawUsername: string) {
    const normalized = normalizeUsernameQuery(rawUsername)
    if (!normalized) return

    username.value = normalized
    isOpen.value = true
    isLoading.value = true
    error.value = null
    profile.value = null

    try {
      profile.value = await userApi.getByUsername(normalized)
    } catch (err) {
      error.value = getApiErrorMessage(err, { fallback: 'Пользователь не найден' })
    } finally {
      isLoading.value = false
    }
  }

  function close() {
    isOpen.value = false
    username.value = ''
    profile.value = null
    error.value = null
    isLoading.value = false
  }

  return {
    isOpen,
    username,
    profile,
    isLoading,
    error,
    open,
    close,
  }
})
