<template>
  <aside class="folder-sidebar" aria-label="Папки чатов">
    <div class="folder-sidebar__brand" aria-label="Megan Messenger">M</div>
    <nav class="folder-sidebar__nav">
      <BaseIconButton
        :label="ALL_FOLDER.label"
        :active="!isProjectsRoute && navigation.activeFolder === ALL_FOLDER.id"
        @click="openFolder(ALL_FOLDER.id)"
      >
        <AppIcon name="folder-all" size="sm" />
      </BaseIconButton>

      <BaseIconButton
        v-for="folder in folders"
        :key="folder.id"
        :label="folder.name"
        :active="!isProjectsRoute && navigation.activeFolder === folder.id"
        @click="openFolder(folder.id)"
        @contextmenu.prevent="openFolderMenu($event, folder)"
      >
        <FolderBadge :name="folder.name" size="sm" />
      </BaseIconButton>

      <span class="folder-sidebar__divider" />

      <BaseIconButton
        label="Проекты"
        :active="isProjectsRoute"
        @click="router.push({ name: 'projects' })"
      >
        <AppIcon name="book" size="sm" />
      </BaseIconButton>

      <BaseIconButton label="Добавить папку" @click="showCreateModal = true">
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

  <ContextMenu
    :open="!!folderMenu"
    :point="folderMenu?.point ?? null"
    :items="folderMenuItems"
    @close="folderMenu = null"
    @select="onFolderMenuSelect"
  />

  <FolderFormModal :open="showCreateModal" mode="create" @close="showCreateModal = false" />
  <FolderFormModal
    :open="!!renameTarget"
    mode="rename"
    :folder-id="renameTarget?.id"
    :initial-name="renameTarget?.name"
    @close="renameTarget = null"
  />
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { onClickOutside } from '@vueuse/core'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import ContextMenu from '@/shared/ui/ContextMenu.vue'
import FolderBadge from '@/entities/folder/ui/FolderBadge.vue'
import { useSessionStore } from '@/entities/session/model/session.store'
import { getUserAvatarColor, getUserDisplayName } from '@/entities/session/lib/display'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import { useProfileUiStore } from '@/features/profile/model/profile-ui.store'
import { useFolders, useDeleteFolder } from '@/entities/folder/api/folder.queries'
import type { ChatFolder } from '@/entities/folder/model/types'
import FolderFormModal from '@/features/create-folder/ui/FolderFormModal.vue'
import { ALL_FOLDER } from '@/shared/config/folders'
import type { MenuPoint } from '@/shared/lib/position-fixed-menu'

const navigation = useConversationSelectionStore()
const session = useSessionStore()
const profileUi = useProfileUiStore()
const router = useRouter()
const route = useRoute()

const isProjectsRoute = computed(() => route.path.startsWith('/projects'))

const menuOpen = ref(false)
const loggingOut = ref(false)
const profileRef = ref<HTMLElement | null>(null)
const showCreateModal = ref(false)
const renameTarget = ref<ChatFolder | null>(null)

const { folders } = useFolders()
const { mutateAsync: deleteFolder } = useDeleteFolder()

const folderMenu = ref<{ point: MenuPoint; folder: ChatFolder } | null>(null)

const folderMenuItems = [
  { id: 'rename', label: 'Переименовать' },
  { id: 'delete', label: 'Удалить', danger: true },
]

onClickOutside(profileRef, () => { menuOpen.value = false })

const displayName = computed(() =>
  getUserDisplayName(session.user?.username, session.user?.phone),
)
const avatarColor = computed(() => getUserAvatarColor(displayName.value))

function openFolder(folderId: string) {
  navigation.selectFolder(folderId)
  if (isProjectsRoute.value) {
    void router.push({ name: 'messenger' })
  }
}

function openFolderMenu(event: MouseEvent, folder: ChatFolder) {
  folderMenu.value = {
    point: { x: event.clientX, y: event.clientY },
    folder,
  }
}

function onFolderMenuSelect(actionId: string) {
  if (!folderMenu.value) return
  const { folder } = folderMenu.value
  folderMenu.value = null

  if (actionId === 'rename') {
    renameTarget.value = folder
    return
  }

  if (actionId === 'delete') {
    void removeFolder(folder)
  }
}

async function removeFolder(folder: ChatFolder) {
  if (navigation.activeFolder === folder.id) {
    navigation.selectFolder(ALL_FOLDER.id)
  }
  await deleteFolder(folder.id)
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
