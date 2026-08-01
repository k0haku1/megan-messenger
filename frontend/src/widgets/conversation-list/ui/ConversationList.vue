<template>
  <aside class="conversation-sidebar" aria-label="Список чатов">
    <header class="conversation-sidebar__header">
      <div class="conversation-sidebar__title">
        <h1>{{ folderTitle }}</h1>
        <BaseIconButton label="Новый чат" variant="ghost"><span class="action-glyph">✎</span></BaseIconButton>
      </div>
      <BaseSearchInput v-model="search" placeholder="Поиск" />
    </header>

    <div class="conversation-list">
      <div v-if="syncState?.status === 'error'" class="sync-notice">
        <span>Нет соединения</span>
        <small>Показаны сохранённые данные</small>
      </div>

      <button
        v-for="(conversation, index) in filteredConversations"
        :key="conversation.id"
        class="conversation-item"
        :class="{ 'is-active': selectedConversationId === conversation.id }"
        type="button"
        @click="navigation.select(conversation.id)"
      >
        <BaseAvatar :name="conversation.title ?? 'Личный диалог'" :color="index" size="lg" />
        <span class="conversation-item__body">
          <span class="conversation-item__row">
            <strong>{{ conversation.title ?? 'Личный диалог' }}</strong>
            <time>{{ conversation.createdAt ? new Date(conversation.createdAt).toLocaleDateString('ru', { day: 'numeric', month: 'short' }) : '' }}</time>
          </span>
          <span class="conversation-item__preview">
            {{ conversation.type === 'group' ? 'Групповой чат' : 'Личные сообщения' }}
          </span>
        </span>
      </button>

      <div v-if="filteredConversations.length === 0" class="conversation-list__empty">
        <span class="conversation-list__empty-icon">⌕</span>
        <strong>{{ search ? 'Ничего не найдено' : 'Здесь пока пусто' }}</strong>
        <span>{{ search ? 'Попробуйте другой запрос' : 'Чаты появятся после синхронизации' }}</span>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { db } from '@/shared/model/db'
import { useLiveQuery } from '@/shared/lib/use-live-query'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import BaseSearchInput from '@/shared/ui/BaseSearchInput.vue'

const conversations = useLiveQuery(() => db.conversations.orderBy('createdAt').reverse().toArray(), [])
const syncState = useLiveQuery(() => db.syncStates.get('conversations'), undefined)
const navigation = useConversationSelectionStore()
const { selectedConversationId, activeFolder } = storeToRefs(navigation)
const search = ref('')

const folderTitle = computed(() => ({ all: 'Все чаты', personal: 'Личные', groups: 'Группы', work: 'Работа' }[activeFolder.value] ?? 'Папка'))
const filteredConversations = computed(() => conversations.value.filter((conversation) => {
  const matchesFolder = activeFolder.value === 'all'
    || activeFolder.value === 'personal' && conversation.type === 'dm'
    || activeFolder.value === 'groups' && conversation.type === 'group'
    || activeFolder.value === 'work'
  const matchesSearch = (conversation.title ?? 'Личный диалог').toLocaleLowerCase('ru').includes(search.value.toLocaleLowerCase('ru'))
  return matchesFolder && matchesSearch
}))
</script>
