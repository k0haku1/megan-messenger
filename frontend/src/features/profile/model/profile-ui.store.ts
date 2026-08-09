import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

export const PROFILE_VIEWS = [
  'main',
  'edit',
  'password',
  'theme',
  'username',
  'privacy',
  'privacy-find',
  'privacy-dm',
] as const
export type ProfileView = (typeof PROFILE_VIEWS)[number]

const VIEW_TITLES: Record<ProfileView, string> = {
  main: 'Мой профиль',
  edit: 'Редактировать профиль',
  password: 'Облачный пароль',
  theme: 'Тема оформления',
  username: 'Username',
  privacy: 'Конфиденциальность',
  'privacy-find': 'Кто может найти меня',
  'privacy-dm': 'Кто может написать мне',
}

const BACK_TARGET: Partial<Record<ProfileView, ProfileView>> = {
  edit: 'main',
  password: 'main',
  theme: 'main',
  username: 'edit',
  privacy: 'main',
  'privacy-find': 'privacy',
  'privacy-dm': 'privacy',
}

export const useProfileUiStore = defineStore('profile-ui', () => {
  const isOpen = ref(false)
  const view = ref<ProfileView>('main')

  const title = computed(() => VIEW_TITLES[view.value])
  const canGoBack = computed(() => view.value !== 'main')

  function open(initialView: ProfileView = 'main') {
    view.value = initialView
    isOpen.value = true
  }

  function close() {
    isOpen.value = false
    view.value = 'main'
  }

  function openView(next: ProfileView) {
    view.value = next
  }

  function back() {
    view.value = BACK_TARGET[view.value] ?? 'main'
  }

  return {
    isOpen,
    view,
    title,
    canGoBack,
    open,
    close,
    openView,
    back,
  }
})
