<template>
  <aside class="folder-sidebar" aria-label="Папки чатов">
    <div class="folder-sidebar__brand" aria-label="Megan Messenger">M</div>
    <nav class="folder-sidebar__nav">
      <BaseIconButton
        v-for="folder in CHAT_FOLDERS"
        :key="folder.id"
        :label="folder.label"
        :active="!isProjectsRoute && navigation.activeFolder === folder.id"
        @click="openFolder(folder.id)"
      >
        <AppIcon :name="folder.icon" size="sm" />
      </BaseIconButton>
      <span class="folder-sidebar__divider" />
      <BaseIconButton
        label="Проекты"
        :active="isProjectsRoute"
        @click="router.push({ name: 'projects' })"
      >
        <AppIcon name="book" size="sm" />
      </BaseIconButton>
      <BaseIconButton label="Добавить папку">
        <AppIcon name="plus" size="sm" />
      </BaseIconButton>
    </nav>

    <div ref="profileRef" class="folder-sidebar__profile">
      <BaseIconButton
        label="Профиль"
        :active="menuOpen"
        @click="menuOpen = !menuOpen"
      >
        <BaseAvatar :name="displayName" :src="session.user?.avatarUrl" size="sm" :color="avatarColor" />
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
          <span class="profile-menu__icon">
            <AppIcon name="user" size="sm" />
          </span>
          <span>Мой профиль</span>
        </button>
        <button
          class="profile-menu__item"
          type="button"
          role="menuitem"
          :disabled="loggingOut"
          @click="onLogout"
        >
          <span class="profile-menu__icon">
            <AppIcon name="logout" size="sm" />
          </span>
          <span>Выйти</span>
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { onClickOutside } from '@vueuse/core'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import { CHAT_FOLDERS, type FolderId } from '@/shared/config/folders'
import { useSessionStore } from '@/entities/session/model/session.store'
import { getUserAvatarColor, getUserDisplayName } from '@/entities/session/lib/display'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import { useProfileUiStore } from '@/features/profile/model/profile-ui.store'

const navigation = useConversationSelectionStore()
const session = useSessionStore()
const profileUi = useProfileUiStore()
const router = useRouter()
const route = useRoute()

const isProjectsRoute = computed(() => route.path.startsWith('/projects'))

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

function openFolder(folderId: FolderId) {
  navigation.selectFolder(folderId)
  if (isProjectsRoute.value) {
    void router.push({ name: 'messenger' })
  }
}

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
