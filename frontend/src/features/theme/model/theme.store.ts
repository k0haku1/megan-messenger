import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import {
  applyThemePreference,
  getStoredThemePreference,
  resolveTheme,
  setStoredThemePreference,
} from '@/shared/theme/applyTheme'
import { THEME_LABELS } from '@/shared/theme/constants'
import type { ThemeMode } from '@/shared/theme/types'

let systemThemeListener: ((event: MediaQueryListEvent) => void) | null = null

export const useThemeStore = defineStore('theme', () => {
  const preference = ref<ThemeMode>(getStoredThemePreference())
  const resolved = computed(() => resolveTheme(preference.value))
  const label = computed(() => THEME_LABELS[preference.value])

  function bindSystemThemeListener() {
    if (typeof window === 'undefined' || systemThemeListener) return

    const media = window.matchMedia('(prefers-color-scheme: dark)')
    systemThemeListener = () => {
      if (preference.value === 'system') {
        applyThemePreference('system')
      }
    }
    media.addEventListener('change', systemThemeListener)
  }

  function init() {
    applyThemePreference(preference.value)
    bindSystemThemeListener()
  }

  function setPreference(mode: ThemeMode) {
    preference.value = mode
    setStoredThemePreference(mode)
    applyThemePreference(mode)
  }

  return {
    preference,
    resolved,
    label,
    init,
    setPreference,
  }
})
