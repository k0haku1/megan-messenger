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
          <BaseIconButton
            v-if="header.type === 'group'"
            label="Информация о группе"
            @click="groupInfoOpen = true"
          >
            <AppIcon name="menu" />
          </BaseIconButton>
          <BaseIconButton v-else label="Информация о чате">
            <AppIcon name="menu" />
          </BaseIconButton>
        </div>
      </header>

      <GroupInfoSheet
        :open="groupInfoOpen"
        :conversation-id="selectedConversationId!"
        :title="header.title"
        :slug="conversation?.slug"
        @close="groupInfoOpen = false"
      />

      <CaptureDecisionModal
        :open="captureOpen"
        :message="captureMessage"
        :projects="linkedProjects"
        @close="closeCapture()"
        @saved="onDecisionSaved"
      />

      <ForwardMessageModal
        :open="forwardOpen"
        :message="forwardMessage"
        :current-conversation-id="selectedConversationId"
        @close="closeForward()"
        @forwarded="onForwarded"
      />

      <MessageContextMenu
        :open="menuOpen"
        :anchor="menuAnchor"
        :items="menuItems"
        @close="closeMenu()"
        @select="handleMenuAction"
      />

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
            <article class="message-bubble" @contextmenu.prevent="openMenu($event, message)">
              <p>{{ message.content }}</p>
              <time>{{ formatMessageTime(message.createdAt) }}</time>
            </article>
          </div>
        </div>
      </div>

      <div class="composer-wrap">
        <div v-if="replyTo" class="composer-reply">
          <div class="composer-reply__meta">
            <strong>Ответ для {{ replyTo.sender.username }}</strong>
            <BaseIconButton label="Отменить ответ" @click="clearReply()">
              <AppIcon name="close" size="sm" />
            </BaseIconButton>
          </div>
          <p>{{ replyTo.content }}</p>
        </div>

        <footer class="composer">
          <BaseIconButton label="Прикрепить файл">
            <AppIcon name="attach" />
          </BaseIconButton>
          <textarea
            v-model="draft"
            rows="1"
            :placeholder="replyTo ? 'Ответ...' : 'Сообщение'"
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
        <p v-if="actionNotice" class="composer-notice">{{ actionNotice }}</p>
      </div>
    </template>

    <div v-else class="chat-placeholder">
      <div class="chat-placeholder__mark">M</div>
      <strong>Выберите чат</strong>
      <span>Откройте диалог слева, чтобы начать общение</span>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import CaptureDecisionModal from '@/features/capture-decision/ui/CaptureDecisionModal.vue'
import ForwardMessageModal from '@/features/forward-message/ui/ForwardMessageModal.vue'
import GroupInfoSheet from '@/features/group-info/ui/GroupInfoSheet.vue'
import MessageContextMenu, {
  type MessageContextMenuAnchor,
  type MessageContextMenuItem,
} from '@/features/message-context-menu/ui/MessageContextMenu.vue'
import { storeToRefs } from 'pinia'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import { getConversationAvatarName, getConversationTitle } from '@/entities/conversation/lib/display'
import { useConversationProjects } from '@/entities/conversation/lib/use-conversation-projects'
import { useConversation } from '@/entities/conversation/lib/use-conversation'
import { hideMessageForMe } from '@/entities/message/api/hidden-message.repository'
import type { Message } from '@/entities/message/model/types'
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
const groupInfoOpen = ref(false)
const captureOpen = ref(false)
const captureMessage = ref<Message | null>(null)
const forwardOpen = ref(false)
const forwardMessage = ref<Message | null>(null)
const actionNotice = ref('')
const replyTo = ref<Message | null>(null)

const menuOpen = ref(false)
const menuAnchor = ref<MessageContextMenuAnchor | null>(null)
const menuMessage = ref<Message | null>(null)

const conversation = useConversation(selectedConversationId)
const visibleMessages = useConversationMessages(selectedConversationId)
const currentUserId = computed(() => session.user?.id)

const isGroupChat = computed(() => conversation.value?.type === 'group')
const linkedProjectsQuery = useConversationProjects(selectedConversationId, isGroupChat)
const linkedProjects = computed(() => linkedProjectsQuery.data.value ?? [])
const canCaptureDecision = computed(() => isGroupChat.value && linkedProjects.value.length > 0)

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

const menuItems = computed<MessageContextMenuItem[]>(() => {
  const items: MessageContextMenuItem[] = [
    { id: 'reply', label: 'Ответить' },
    { id: 'forward', label: 'Переслать' },
  ]

  if (canCaptureDecision.value) {
    items.push({ id: 'capture', label: 'Зафиксировать решение' })
  }

  items.push({ id: 'delete', label: 'Удалить для себя', danger: true })

  return items
})

const { draft, isSending, sendError, canSend, send, clearReply } = useMessageComposer({
  selectedConversationId,
  pendingPeer,
  replyTo,
  onConversationOpened: (conversationId) => navigation.select(conversationId),
})

function showNotice(text: string) {
  actionNotice.value = text
  window.setTimeout(() => {
    actionNotice.value = ''
  }, 4000)
}

function openMenu(event: MouseEvent, message: Message) {
  const bubble = event.currentTarget as HTMLElement
  const rect = bubble.getBoundingClientRect()
  const own = isOwnMessage(message, currentUserId.value)

  menuMessage.value = message
  menuAnchor.value = {
    top: rect.top,
    left: rect.left,
    right: rect.right,
    bottom: rect.bottom,
    align: own ? 'end' : 'start',
  }
  menuOpen.value = true
}

function closeMenu() {
  menuOpen.value = false
  menuAnchor.value = null
  menuMessage.value = null
}

function handleMenuAction(actionId: string) {
  const message = menuMessage.value
  if (!message) return

  switch (actionId) {
    case 'reply':
      replyTo.value = message
      break
    case 'forward':
      forwardMessage.value = message
      forwardOpen.value = true
      break
    case 'capture':
      captureMessage.value = message
      captureOpen.value = true
      break
    case 'delete':
      void hideMessageForMe(message.id)
      showNotice('Сообщение скрыто для вас')
      break
  }
}

function closeCapture() {
  captureOpen.value = false
  captureMessage.value = null
}

function closeForward() {
  forwardOpen.value = false
  forwardMessage.value = null
}

function onDecisionSaved(projectId: string) {
  const project = linkedProjects.value.find((item) => item.id === projectId)
  showNotice(project ? `Решение сохранено в «${project.name}»` : 'Решение сохранено')
}

function onForwarded() {
  showNotice('Сообщение переслано')
}
</script>

<style scoped>
.composer-wrap {
  width: min(calc(100% - 32px), 780px);
  margin: 0 auto 14px;
}
.composer-wrap :deep(.composer) {
  width: 100%;
  margin: 0;
}
.composer-reply {
  margin-bottom: 8px;
  padding: 8px 10px;
  background: var(--color-bg-subtle);
  border-left: 3px solid var(--color-accent);
  border-radius: var(--radius-sm);
}
.composer-reply__meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 4px;
}
.composer-reply__meta strong {
  font-size: var(--font-size-sm);
}
.composer-reply p {
  margin: 0;
  overflow: hidden;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  text-overflow: ellipsis;
  white-space: nowrap;
}
.composer-error,
.composer-notice {
  margin: 8px 0 0;
  font-size: 0.875rem;
}
.composer-error {
  color: var(--color-danger, #e5484d);
}
.composer-notice {
  color: var(--color-accent);
}
</style>
