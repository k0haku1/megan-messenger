import { ref, watch } from 'vue'
import { defineStore } from 'pinia'
import { useDebounceFn } from '@vueuse/core'
import { userApi } from '@/entities/user/api/user.api'
import { normalizeUsernameQuery } from '@/entities/user/lib/username'
import type { UserSearchResult } from '@/entities/user/model/types'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'

export const useGlobalSearchStore = defineStore('global-search', () => {
  const query = ref('')
  const userResults = ref<UserSearchResult[]>([])
  const isSearching = ref(false)
  const searchError = ref<string | null>(null)

  const runSearch = useDebounceFn(async (rawQuery: string) => {
    const normalized = normalizeUsernameQuery(rawQuery)
    if (normalized.length < 2) {
      userResults.value = []
      searchError.value = null
      isSearching.value = false
      return
    }

    isSearching.value = true
    searchError.value = null

    try {
      const response = await userApi.search(normalized)
      userResults.value = response.users
    } catch (error) {
      userResults.value = []
      searchError.value = getApiErrorMessage(error, { fallback: 'Не удалось выполнить поиск' })
    } finally {
      isSearching.value = false
    }
  }, 300)

  watch(query, (value) => {
    void runSearch(value)
  })

  function clear() {
    query.value = ''
    userResults.value = []
    searchError.value = null
    isSearching.value = false
  }

  return {
    query,
    userResults,
    isSearching,
    searchError,
    clear,
  }
})
