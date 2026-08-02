import { ref } from 'vue'
import { userApi } from '@/entities/user/api/user.api'
import { useSessionStore } from '@/entities/session/model/session.store'
import type { DMPolicy } from '@/entities/user/model/types'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'

export function usePrivacySettings() {
  const session = useSessionStore()
  const isSaving = ref(false)
  const error = ref('')

  async function save(usernameSearchable: boolean, dmPolicy: DMPolicy): Promise<boolean> {
    error.value = ''
    isSaving.value = true

    try {
      const result = await userApi.updatePrivacy({ usernameSearchable, dmPolicy })
      if (session.user) {
        session.user.usernameSearchable = result.usernameSearchable
        session.user.dmPolicy = result.dmPolicy
      }
      return true
    } catch (err) {
      error.value = getApiErrorMessage(err, { fallback: 'Не удалось сохранить настройки' })
      return false
    } finally {
      isSaving.value = false
    }
  }

  return {
    session,
    isSaving,
    error,
    save,
  }
}
