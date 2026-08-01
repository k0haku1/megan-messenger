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
    <div class="folder-sidebar__profile">
      <BaseIconButton label="Профиль"><BaseAvatar name="Алексей" size="sm" :color="4" /></BaseIconButton>
    </div>
  </aside>
</template>

<script setup lang="ts">
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'

const navigation = useConversationSelectionStore()
const folders = [
  { id: 'all', icon: '◉', label: 'Все чаты' },
  { id: 'personal', icon: '♙', label: 'Личные' },
  { id: 'groups', icon: '♟', label: 'Группы' },
  { id: 'work', icon: '▣', label: 'Работа' },
] as const
</script>
