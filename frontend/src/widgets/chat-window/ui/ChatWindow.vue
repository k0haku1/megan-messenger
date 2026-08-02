<template>
  <section class="chat-window" aria-label="Переписка">
    <template v-if="header">
      <header class="chat-header">
        <BaseAvatar :name="header.avatarName" :color="header.type === 'group' ? 2 : 0" />
        <div class="chat-header__meta">
          <strong>{{ header.title }}</strong>
          <span>{{ header.type === 'group' ? 'групповой чат' : 'был(а) недавно' }}</span>
        </div>
        <div class="chat-header__actions">
          <BaseIconButton label="Поиск по сообщениям">
            <AppIcon name="search" />
          </BaseIconButton>
          <BaseIconButton label="Информация о чате">
            <AppIcon name="menu" />
          </BaseIconButton>
        </div>
      </header>

      <div class="message-area">
        <div v-if="visibleMessages.length === 0" class="chat-empty">
          <AppIcon class="chat-empty__icon" name="mail" />
          <strong>{{ header.isPending ? 'Напишите первое сообщение' : 'Начните общение' }}</strong>
          <span>{{
            header.isPending
              ? 'Диалог появится в списке после отправки'
              : 'Сообщения этого чата появятся здесь'
          }}</span>
        </div>
        <div v-else class="message-stack">
          <div
            v-for="(message, index) in visibleMessages"
            :key="message.id"
            class="message-item"
            :class="{ 'is-own': isOwnMessage(message, currentUserId) }"
          >
            <span v-if="shouldShowSenderName(message, index, visibleMessages, currentUserId)" class="message-item__sender">
              {{ message.sender.username }}
            </span>
            <article class="message-bubble">
              <p>{{ message.content }}</p>
              <time>{{ formatMessageTime(message.createdAt) }}</time>
            </article>
          </div>
        </div>
      </div>

      <footer class="composer">
        <BaseIconButton label="Прикрепить файл">
          <AppIcon name="attach" />
        </BaseIconButton>
        <textarea
          v-model="draft"
          rows="1"
          placeholder="Сообщение"
          aria-label="Сообщение"
          :disabled="isSending"
          @keydown.enter.exact.prevent="send()"
        />
        <BaseIconButton label="Эмодзи">
          <AppIcon name="emoji" />
        </BaseIconButton>
        <BaseIconButton label="Отправить" variant="primary" :disabled="!canSend" @click="send()">
          <AppIcon name="send" />
        </BaseIconButton>
      </footer>

      <p v-if="sendError" class="composer-error">{{ sendError }}</p>
    </template>

    <div v-else class="chat-placeholder">
      <div class="chat-placeholder__mark">M</div>
      <strong>Выберите чат</strong>
      <span>Откройте диалог слева, чтобы начать общение</span>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import { getConversationAvatarName, getConversationTitle } from '@/entities/conversation/lib/display'
import { useConversation } from '@/entities/conversation/lib/use-conversation'
import {
  formatMessageTime,
  isOwnMessage,
  shouldShowSenderName,
} from '@/entities/message/lib/display'
import { useConversationMessages } from '@/entities/message/lib/use-conversation-messages'
import { useMessageComposer } from '@/features/compose-message/model/use-message-composer'
import { useMessagesSync } from '@/features/message-sync/model/use-messages-sync'
import { useConversationWs } from '@/features/conversation-ws/model/use-conversation-ws'
import { useSessionStore } from '@/entities/session/model/session.store'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'

const navigation = useConversationSelectionStore()
const session = useSessionStore()
const { selectedConversationId, pendingPeer } = storeToRefs(navigation)

const conversation = useConversation(selectedConversationId)
const visibleMessages = useConversationMessages(selectedConversationId)
const currentUserId = computed(() => session.user?.id)

useMessagesSync(selectedConversationId)
useConversationWs(selectedConversationId)

const header = computed(() => {
  if (conversation.value) {
    return {
      type: conversation.value.type,
      title: getConversationTitle(conversation.value),
      avatarName: getConversationAvatarName(conversation.value),
      isPending: false,
    }
  }

  if (pendingPeer.value) {
    return {
      type: 'dm' as const,
      title: pendingPeer.value.username,
      avatarName: pendingPeer.value.username,
      isPending: true,
    }
  }

  return null
})

const { draft, isSending, sendError, canSend, send } = useMessageComposer({
  selectedConversationId,
  pendingPeer,
  onConversationOpened: (conversationId) => navigation.select(conversationId),
})
</script>

<style scoped>
.composer-error {
  margin: 0;
  padding: 0.5rem 1rem 1rem;
  color: var(--color-danger, #e5484d);
  font-size: 0.875rem;
}
</style>
