<template>
  <aside class="folder-sidebar" aria-label="Папки чатов">
    <div class="folder-sidebar__brand" aria-label="Megan Messenger">M</div>
    <nav class="folder-sidebar__nav">
      <BaseIconButton
        v-for="folder in folders"
        :key="folder.id"
        :label="folder.label"
        :active="navigation.activeFolder === folder.id"
        @click="navigation.selectFolder(folder.id)"
      >
        <span class="folder-icon" aria-hidden="true">{{ folder.icon }}</span>
      </BaseIconButton>
      <span class="folder-sidebar__divider" />
      <BaseIconButton label="Добавить папку"><span class="folder-icon folder-icon--plus">＋</span></BaseIconButton>
    </nav>

    <div ref="profileRef" class="folder-sidebar__profile">
      <BaseIconButton
        label="Профиль"
        :active="menuOpen"
        @click="menuOpen = !menuOpen"
      >
        <BaseAvatar :name="displayName" size="sm" :color="avatarColor" />
      </BaseIconButton>

      <div
        v-if="menuOpen"
        class="profile-menu"
        role="menu"
        aria-label="Меню профиля"
      >
        <button
          class="profile-menu__item"
          type="button"
          role="menuitem"
          @click="onOpenProfile"
        >
          <span class="profile-menu__icon" aria-hidden="true">☺</span>
          <span>Мой профиль</span>
        </button>
        <button
          class="profile-menu__item"
          type="button"
          role="menuitem"
          :disabled="loggingOut"
          @click="onLogout"
        >
          <span class="profile-menu__icon" aria-hidden="true">⎋</span>
          <span>Выйти</span>
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { onClickOutside } from '@vueuse/core'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import { useSessionStore } from '@/entities/session/model/session.store'
import { useProfileUiStore } from '@/features/profile/model/profile-ui.store'
import { router } from '@/app/router'

const navigation = useConversationSelectionStore()
const session = useSessionStore()
const profileUi = useProfileUiStore()

const menuOpen = ref(false)
const loggingOut = ref(false)
const profileRef = ref<HTMLElement | null>(null)

onClickOutside(profileRef, () => {
  menuOpen.value = false
})

const displayName = computed(() => session.user?.username || session.user?.phone || 'User')
const avatarColor = computed(() => {
  const seed = displayName.value
  let hash = 0
  for (let i = 0; i < seed.length; i += 1) hash = (hash + seed.charCodeAt(i) * (i + 1)) % 6
  return hash
})

const folders = [
  { id: 'all', icon: '◉', label: 'Все чаты' },
  { id: 'personal', icon: '♙', label: 'Личные' },
  { id: 'groups', icon: '♟', label: 'Группы' },
  { id: 'work', icon: '▣', label: 'Работа' },
] as const

function onOpenProfile() {
  menuOpen.value = false
  profileUi.open()
}

async function onLogout() {
  if (loggingOut.value) return
  loggingOut.value = true
  menuOpen.value = false
  try {
    await session.logout()
    await router.replace({ name: 'auth' })
  } finally {
    loggingOut.value = false
  }
}
</script>
