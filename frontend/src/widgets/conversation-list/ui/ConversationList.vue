<template>
  <aside class="conversation-sidebar" aria-label="Список чатов">
    <header class="conversation-sidebar__header">
      <div class="conversation-sidebar__title">
        <h1>{{ folderTitle }}</h1>
        <BaseIconButton label="Новый чат" variant="ghost" @click="focusSearch()">
          <AppIcon name="edit" size="sm" />
        </BaseIconButton>
      </div>
      <BaseSearchInput
        ref="searchInputRef"
        v-model="searchStore.query"
        placeholder="Поиск или @username"
      />
    </header>

    <GlobalSearchResults />

    <div class="conversation-list">
      <div v-if="isError" class="sync-notice">
        <span>Нет соединения</span>
        <small>Не удалось обновить список чатов</small>
      </div>

      <button
        v-for="(conversation, index) in filteredConversations"
        :key="conversation.id"
        class="conversation-item"
        :class="{ 'is-active': selectedConversationId === conversation.id }"
        type="button"
        @click="navigation.select(conversation.id)"
      >
        <BaseAvatar :name="getConversationAvatarName(conversation)" :color="index" size="lg" />
        <span class="conversation-item__body">
          <span class="conversation-item__row">
            <strong>{{ getConversationTitle(conversation) }}</strong>
            <time>{{
              conversation.createdAt
                ? new Date(conversation.createdAt).toLocaleDateString('ru', {
                    day: 'numeric',
                    month: 'short',
                  })
                : ''
            }}</time>
          </span>
          <span class="conversation-item__preview">
            {{ conversation.type === 'group' ? 'Групповой чат' : 'Личные сообщения' }}
          </span>
        </span>
      </button>

      <div
        v-if="filteredConversations.length === 0 && !showUserResults"
        class="conversation-list__empty"
      >
        <AppIcon class="conversation-list__empty-icon" name="search" />
        <strong>{{ searchStore.query ? 'Чаты не найдены' : 'Здесь пока пусто' }}</strong>
        <span>{{
          searchStore.query ? 'Попробуйте @username' : 'Чаты появятся после синхронизации'
        }}</span>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useConversations } from '@/entities/conversation/api/conversation.queries'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import {
  getConversationAvatarName,
  getConversationTitle,
} from '@/entities/conversation/lib/display'
import { normalizeUsernameQuery } from '@/entities/user/lib/username'
import GlobalSearchResults from '@/features/global-search/ui/GlobalSearchResults.vue'
import { useGlobalSearchStore } from '@/features/global-search/model/global-search.store'
import AppIcon from '@/shared/ui/AppIcon.vue'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import BaseSearchInput from '@/shared/ui/BaseSearchInput.vue'
import { CHAT_FOLDER_TITLES } from '@/shared/config/folders'
import type { FolderId } from '@/shared/config/folders'

const { conversations, isError } = useConversations()
const navigation = useConversationSelectionStore()
const searchStore = useGlobalSearchStore()
const { selectedConversationId, activeFolder } = storeToRefs(navigation)
const searchInputRef = ref<{ focus: () => void } | null>(null)

const folderTitle = computed(
  () => CHAT_FOLDER_TITLES[activeFolder.value as FolderId] ?? 'Папка',
)

const showUserResults = computed(() => normalizeUsernameQuery(searchStore.query).length >= 2)

const filteredConversations = computed(() =>
  conversations.value.filter((conversation) => {
    const matchesFolder =
      activeFolder.value === 'all'
      || (activeFolder.value === 'personal' && conversation.type === 'dm')
      || (activeFolder.value === 'groups' && conversation.type === 'group')
      || activeFolder.value === 'work'

    const title = getConversationTitle(conversation).toLocaleLowerCase('ru')
    const query = searchStore.query.toLocaleLowerCase('ru').trim()
    const matchesSearch = !query || title.includes(query.replace(/^@+/, ''))

    return matchesFolder && matchesSearch
  }),
)

function focusSearch() {
  searchInputRef.value?.focus()
}
</script>
