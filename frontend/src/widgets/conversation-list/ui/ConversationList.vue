<template>
  <aside class="conversation-sidebar" aria-label="Список чатов">
    <header class="conversation-sidebar__header">
      <div class="conversation-sidebar__title">
        <h1>{{ folderTitle }}</h1>
        <NewChatMenu @focus-search="focusSearch()" />
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
        :class="{
          'is-active': selectedConversationId === conversation.id,
          'is-dragging': drag?.conversationId === conversation.id && drag.active,
        }"
        type="button"
        @pointerdown="onItemPointerDown($event, conversation)"
      >
        <BaseAvatar
          :name="getConversationAvatarName(conversation)"
          :src="getConversationAvatarUrl(conversation)"
          :color="index"
          size="lg"
        />
        <span class="conversation-item__body">
          <span class="conversation-item__row">
            <strong>{{ getConversationTitle(conversation) }}</strong>
            <time>{{ formatConversationListTime(getConversationActivityAt(conversation)) }}</time>
          </span>
          <span class="conversation-item__preview">
            {{ getConversationPreview(conversation) }}
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

    <Teleport to="body">
      <div
        v-if="drag?.active"
        class="conversation-drag-ghost"
        :style="{
          transform: `translate3d(${drag.clientX - 28}px, ${drag.clientY - 28}px, 0)`,
        }"
      >
        <BaseAvatar :name="drag.title" size="lg" />
        <span>{{ drag.title }}</span>
      </div>
    </Teleport>
  </aside>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useConversations } from '@/entities/conversation/api/conversation.queries'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import {
  formatConversationListTime,
  getConversationActivityAt,
  getConversationAvatarName,
  getConversationAvatarUrl,
  getConversationPreview,
  getConversationTitle,
} from '@/entities/conversation/lib/display'
import type { Conversation } from '@/entities/conversation/model/types'
import { normalizeUsernameQuery } from '@/entities/user/lib/username'
import { useChatWorkspaceStore } from '@/features/chat-workspace/model/chat-workspace.store'
import { useConversationListDrag } from '@/features/chat-workspace/model/use-conversation-list-drag'
import GlobalSearchResults from '@/features/global-search/ui/GlobalSearchResults.vue'
import NewChatMenu from '@/features/new-chat/ui/NewChatMenu.vue'
import { useGlobalSearchStore } from '@/features/global-search/model/global-search.store'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import BaseSearchInput from '@/shared/ui/BaseSearchInput.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import { CHAT_FOLDER_TITLES } from '@/shared/config/folders'
import type { FolderId } from '@/shared/config/folders'

const { conversations, isError } = useConversations()
const navigation = useConversationSelectionStore()
const workspace = useChatWorkspaceStore()
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

const { drag, startDrag } = useConversationListDrag({
  onClick: (conversationId) => navigation.select(conversationId),
  onDrop: ({ conversationId, clientX, clientY, originRect }) => {
    workspace.setListDragActive(false)
    const dropped = workspace.dropConversationAt(
      conversationId,
      clientX,
      clientY,
      originRect,
    )
    if (dropped) {
      navigation.focusConversation(conversationId)
      return
    }
    navigation.select(conversationId)
  },
  onCancel: () => workspace.setListDragActive(false),
})

watch(
  () => drag.value?.active === true,
  (active) => {
    workspace.setListDragActive(Boolean(active))
  },
)

watch(
  () => {
    if (!drag.value?.active) return null
    return {
      x: drag.value.clientX,
      y: drag.value.clientY,
    }
  },
  (point) => {
    if (!point) {
      workspace.clearDropPreview()
      return
    }
    workspace.updateDropPreview(point.x, point.y)
  },
)

function onItemPointerDown(event: PointerEvent, conversation: Conversation): void {
  const originEl = event.currentTarget
  if (!(originEl instanceof HTMLElement)) return
  startDrag({
    event,
    conversationId: conversation.id,
    title: getConversationTitle(conversation),
    originEl,
  })
}

function focusSearch() {
  searchInputRef.value?.focus()
}
</script>
