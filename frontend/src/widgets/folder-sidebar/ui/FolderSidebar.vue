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
import { useRouter } from 'vue-router'
import { onClickOutside } from '@vueuse/core'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import { useSessionStore } from '@/entities/session/model/session.store'
import { getUserAvatarColor, getUserDisplayName } from '@/entities/session/lib/display'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import { useProfileUiStore } from '@/features/profile/model/profile-ui.store'

const navigation = useConversationSelectionStore()
const session = useSessionStore()
const profileUi = useProfileUiStore()
const router = useRouter()

const menuOpen = ref(false)
const loggingOut = ref(false)
const profileRef = ref<HTMLElement | null>(null)

onClickOutside(profileRef, () => {
  menuOpen.value = false
})

const displayName = computed(() =>
  getUserDisplayName(session.user?.username, session.user?.phone),
)
const avatarColor = computed(() => getUserAvatarColor(displayName.value))

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
    navigation.select(null)
    await session.logout()
    await router.replace({ name: 'auth' })
  } finally {
    loggingOut.value = false
  }
}
</script>
