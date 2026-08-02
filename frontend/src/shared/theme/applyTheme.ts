import { THEME_META_COLORS, THEME_STORAGE_KEY } from './constants'
import type { ResolvedTheme, ThemeMode } from './types'

const isThemeMode = (value: string | null): value is ThemeMode =>
  value === 'light' || value === 'dark' || value === 'system'

export const getStoredThemePreference = (): ThemeMode => {
  if (typeof localStorage === 'undefined') return 'system'

  const saved = localStorage.getItem(THEME_STORAGE_KEY)
  return isThemeMode(saved) ? saved : 'system'
}

export const setStoredThemePreference = (mode: ThemeMode): void => {
  localStorage.setItem(THEME_STORAGE_KEY, mode)
}

export const resolveTheme = (preference: ThemeMode): ResolvedTheme => {
  if (preference === 'system') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  }

  return preference
}

const updateThemeColorMeta = (resolved: ResolvedTheme): void => {
  const meta = document.querySelector('meta[name="theme-color"]')
  if (!meta) return

  meta.setAttribute('content', THEME_META_COLORS[resolved])
}

export const applyResolvedTheme = (resolved: ResolvedTheme): void => {
  document.documentElement.dataset.theme = resolved
  document.documentElement.style.colorScheme = resolved
  updateThemeColorMeta(resolved)
}

export const applyThemePreference = (preference: ThemeMode): ResolvedTheme => {
  const resolved = resolveTheme(preference)
  applyResolvedTheme(resolved)
  return resolved
}
