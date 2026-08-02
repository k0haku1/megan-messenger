<template>
  <div ref="menuRef" class="new-chat-menu">
    <BaseIconButton label="Новый чат" variant="ghost" :active="menuOpen" @click="menuOpen = !menuOpen">
      <AppIcon name="plus" size="sm" />
    </BaseIconButton>

    <div v-if="menuOpen" class="new-chat-menu__panel" role="menu" aria-label="Создание чата">
      <button class="new-chat-menu__item" type="button" role="menuitem" @click="openDirect()">
        <AppIcon name="edit" size="sm" />
        <span>Личное сообщение</span>
      </button>
      <button class="new-chat-menu__item" type="button" role="menuitem" @click="openCreateGroup()">
        <AppIcon name="folder-groups" size="sm" />
        <span>Создать группу</span>
      </button>
      <button class="new-chat-menu__item" type="button" role="menuitem" @click="openJoinGroup()">
        <AppIcon name="link" size="sm" />
        <span>Вступить по ссылке</span>
      </button>
    </div>

    <CreateGroupModal :open="createGroupOpen" @close="createGroupOpen = false" />
    <JoinGroupModal :open="joinGroupOpen" @close="joinGroupOpen = false" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onClickOutside } from '@vueuse/core'
import CreateGroupModal from '@/features/create-group/ui/CreateGroupModal.vue'
import JoinGroupModal from '@/features/join-group/ui/JoinGroupModal.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'

const emit = defineEmits<{ focusSearch: [] }>()

const menuRef = ref<HTMLElement | null>(null)
const menuOpen = ref(false)
const createGroupOpen = ref(false)
const joinGroupOpen = ref(false)

onClickOutside(menuRef, () => {
  menuOpen.value = false
})

function openDirect() {
  menuOpen.value = false
  emit('focusSearch')
}

function openCreateGroup() {
  menuOpen.value = false
  createGroupOpen.value = true
}

function openJoinGroup() {
  menuOpen.value = false
  joinGroupOpen.value = true
}
</script>

<style scoped>
.new-chat-menu {
  position: relative;
}
.new-chat-menu__panel {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  z-index: 20;
  display: flex;
  min-width: 220px;
  flex-direction: column;
  padding: 6px;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-panel);
}
.new-chat-menu__item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  text-align: left;
  color: var(--color-text-primary);
  background: transparent;
  border: 0;
  border-radius: var(--radius-sm);
  cursor: pointer;
}
.new-chat-menu__item:hover {
  background: var(--color-bg-hover);
}
</style>
