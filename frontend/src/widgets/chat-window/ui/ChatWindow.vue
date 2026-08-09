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
        :reaction-details="menuReactionDetails"
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
            :class="{ 'is-own': isOwnMessage(message, currentUserId), 'is-deleting': deletingMessageIds.has(message.id) }"
          >
            <span v-if="shouldShowSenderName(message, index, visibleMessages, currentUserId)" class="message-item__sender">
              {{ message.sender.username }}
            </span>
            <article
              class="message-bubble"
              :class="{ 'is-context-pressed': contextPressedMessageId === message.id }"
              @contextmenu.prevent="openMenu($event, message)"
              @mouseenter="scheduleReactionPicker(message.id)"
              @mouseleave="hideReactionPicker()"
            >
              <template v-if="repliedMessage(message)">
                <div class="message-bubble__reply-preview">
                  <strong>{{ repliedMessage(message)?.sender.username }}</strong>
                  <span>{{ replyPreview(repliedMessage(message)!) }}</span>
                </div>
              </template>
              <span v-else-if="message.replyToId" class="message-bubble__reference">↩ Исходное сообщение недоступно</span>
              <span v-if="message.forwardedFromId" class="message-bubble__reference">↪ Пересланное сообщение</span>
              <p>{{ message.content }}</p>
              <div v-if="message.reactions?.length" class="message-bubble__reactions">
                <button
                  v-for="reaction in message.reactions"
                  :key="reaction.emoji"
                  type="button"
                  :class="{ 'is-active': reaction.reactedByMe, 'is-appearing': appearingReactionKey === reactionKey(message.id, reaction.emoji) }"
                  @click="toggleReaction(message, reaction.emoji, reaction.reactedByMe)"
                >
                  {{ reaction.emoji }} <small>{{ reaction.count }}</small>
                </button>
              </div>
              <div
                v-if="reactionPickerMessageId === message.id && !message.deletedAt"
                class="message-reaction-picker"
                @mouseenter="cancelReactionPickerHide()"
                @mouseleave="hideReactionPicker()"
              >
                <button
                  v-for="emoji in reactionEmojis"
                  :key="emoji"
                  type="button"
                  :class="{ 'is-active': hasReaction(message, emoji) }"
                  :aria-label="`Поставить реакцию ${emoji}`"
                  @click="toggleReaction(message, emoji, hasReaction(message, emoji))"
                >{{ emoji }}</button>
                <button
                  type="button"
                  aria-label="Выбрать другую реакцию"
                  @click="toggleEmojiPicker(message.id)"
                >＋</button>
                <emoji-picker
                  v-if="emojiPickerMessageId === message.id"
                  class="message-reaction-picker__emoji-picker"
                  @emoji-click="onEmojiClick(message, $event)"
                />
              </div>
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
import 'emoji-picker-element'
import { computed, ref } from 'vue'
import CaptureDecisionModal from '@/features/capture-decision/ui/CaptureDecisionModal.vue'
import ForwardMessageModal from '@/features/forward-message/ui/ForwardMessageModal.vue'
import GroupInfoSheet from '@/features/group-info/ui/GroupInfoSheet.vue'
import MessageContextMenu, {
  type MessageContextMenuAnchor,
  type MessageContextMenuItem,
  type MessageContextMenuReactionDetail,
} from '@/features/message-context-menu/ui/MessageContextMenu.vue'
import { storeToRefs } from 'pinia'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import { getConversationAvatarName, getConversationTitle } from '@/entities/conversation/lib/display'
import { useConversationProjects } from '@/entities/conversation/lib/use-conversation-projects'
import { useConversation } from '@/entities/conversation/lib/use-conversation'
import { hideMessageForMe } from '@/entities/message/api/hidden-message.repository'
import { messageApi } from '@/entities/message/api/message.api'
import { removeMessage, upsertMessage } from '@/entities/message/api/message.repository'
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
const reactionPickerMessageId = ref<string | null>(null)
const emojiPickerMessageId = ref<string | null>(null)
const deletingMessageIds = ref(new Set<string>())
const contextPressedMessageId = ref<string | null>(null)
const appearingReactionKey = ref<string | null>(null)
let reactionPickerTimer: ReturnType<typeof window.setTimeout> | null = null
const reactionEmojis = ['👍', '❤️', '😂']

const conversation = useConversation(selectedConversationId)
const visibleMessages = useConversationMessages(selectedConversationId)
const currentUserId = computed(() => session.user?.id)
const messagesById = computed(() => new Map(visibleMessages.value.map((message) => [message.id, message])))

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

  if (menuMessage.value && isOwnMessage(menuMessage.value, currentUserId.value) && !menuMessage.value.deletedAt) {
    items.push({ id: 'delete-everyone', label: 'Удалить для всех', danger: true })
  }

  return items
})

const menuReactionDetails = computed<MessageContextMenuReactionDetail[]>(() =>
  (menuMessage.value?.reactions ?? [])
    .map((reaction) => ({
      emoji: reaction.emoji,
      users: (reaction.reactedBy ?? []).map((user) => user.username).filter(Boolean),
    }))
    .filter((reaction) => reaction.users.length > 0),
)

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
  contextPressedMessageId.value = message.id
  window.setTimeout(() => {
    if (contextPressedMessageId.value === message.id) contextPressedMessageId.value = null
  }, 180)
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

async function handleMenuAction(actionId: string) {
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
      await deleteForMe(message)
      break
    case 'delete-everyone':
      await deleteForEveryone(message)
      break
  }
}

function repliedMessage(message: Message): Message | undefined {
  return message.replyToId ? messagesById.value.get(message.replyToId) : undefined
}

function replyPreview(message: Message): string {
  return message.deletedAt ? 'Сообщение удалено' : message.content
}

function hasReaction(message: Message, emoji: string): boolean {
  return message.reactions?.some((reaction) => reaction.emoji === emoji && reaction.reactedByMe) ?? false
}

function scheduleReactionPicker(messageId: string) {
  if (reactionPickerTimer) window.clearTimeout(reactionPickerTimer)
  reactionPickerTimer = window.setTimeout(() => {
    reactionPickerMessageId.value = messageId
    emojiPickerMessageId.value = null
    reactionPickerTimer = null
  }, 1500)
}

function hideReactionPicker() {
  if (reactionPickerTimer) window.clearTimeout(reactionPickerTimer)
  reactionPickerTimer = window.setTimeout(() => {
    reactionPickerMessageId.value = null
    emojiPickerMessageId.value = null
    reactionPickerTimer = null
  }, 100)
}

function cancelReactionPickerHide() {
  if (reactionPickerTimer) {
    window.clearTimeout(reactionPickerTimer)
    reactionPickerTimer = null
  }
}

function toggleEmojiPicker(messageId: string) {
  emojiPickerMessageId.value = emojiPickerMessageId.value === messageId ? null : messageId
}

function onEmojiClick(message: Message, event: Event) {
  const emoji = (event as CustomEvent<{ unicode: string }>).detail.unicode
  if (!emoji) return
  void toggleReaction(message, emoji, hasReaction(message, emoji))
  emojiPickerMessageId.value = null
}

async function deleteForMe(message: Message) {
  const conversationId = selectedConversationId.value
  if (!conversationId) return
  try {
    deletingMessageIds.value = new Set([...deletingMessageIds.value, message.id])
    await Promise.all([messageApi.hideForMe(conversationId, message.id), delay(260)])
    await hideMessageForMe(message.id)
  } catch {
    removeDeletingState(message.id)
    showNotice('Не удалось удалить сообщение')
  }
}

async function deleteForEveryone(message: Message) {
  const conversationId = selectedConversationId.value
  if (!conversationId) return
  try {
    deletingMessageIds.value = new Set([...deletingMessageIds.value, message.id])
    await Promise.all([messageApi.deleteForEveryone(conversationId, message.id), delay(260)])
    await removeMessage(message.id)
  } catch {
    removeDeletingState(message.id)
    showNotice('Не удалось удалить сообщение для всех')
  }
}

async function toggleReaction(message: Message, emoji: string, active: boolean) {
  const conversationId = selectedConversationId.value
  if (!conversationId || message.deletedAt) return
  try {
    const result = await messageApi.react(conversationId, message.id, emoji, !active)
    await upsertMessage(result.message)
    if (!active) {
      appearingReactionKey.value = reactionKey(message.id, emoji)
      window.setTimeout(() => { appearingReactionKey.value = null }, 280)
    }
  } catch {
    showNotice('Не удалось изменить реакцию')
  }
}

function reactionKey(messageId: string, emoji: string): string { return `${messageId}:${emoji}` }
function delay(milliseconds: number): Promise<void> { return new Promise((resolve) => window.setTimeout(resolve, milliseconds)) }
function removeDeletingState(messageId: string) {
  const next = new Set(deletingMessageIds.value)
  next.delete(messageId)
  deletingMessageIds.value = next
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

function onForwarded() {}
</script>

<style scoped>
.composer-wrap {
  width: min(calc(100% - 32px), 780px);
  margin: 0 auto 14px;
}
.message-item.is-deleting { overflow: hidden; animation: message-delete 260ms ease-in forwards; }
.message-bubble { position: relative; transition: transform 140ms ease, filter 140ms ease; }
.message-bubble.is-context-pressed { filter: brightness(0.96); transform: scale(0.975); }
.message-bubble__reference { display: block; margin-bottom: 4px; color: var(--color-text-secondary); font-size: var(--font-size-sm); }
.message-bubble__reply-preview { display: flex; flex-direction: column; gap: 2px; margin-bottom: 8px; padding: 6px 8px; border-left: 3px solid var(--color-accent); background: var(--color-bg-subtle); border-radius: var(--radius-sm); font-size: var(--font-size-sm); }
.message-bubble__reply-preview span { overflow: hidden; color: var(--color-text-secondary); text-overflow: ellipsis; white-space: nowrap; }
.message-bubble p.is-deleted { color: var(--color-text-secondary); font-style: italic; }
.message-bubble__reactions { display: flex; flex-wrap: wrap; gap: 4px; margin: 6px 0 0; }
.message-bubble__reactions button { position: relative; padding: 2px 6px; color: inherit; background: var(--color-bg-subtle); border: 1px solid var(--color-border); border-radius: 999px; cursor: pointer; }
.message-bubble__reactions button.is-active { border-color: var(--color-accent); background: color-mix(in srgb, var(--color-accent) 16%, transparent); }
.message-bubble__reactions button.is-appearing { animation: reaction-pop 280ms cubic-bezier(.2, .9, .25, 1.35); }
.message-reaction-picker { position: absolute; top: 0; display: flex; flex-direction: column; gap: 3px; padding: 4px; background: var(--color-bg-elevated); border: 1px solid var(--color-border); border-radius: 999px; box-shadow: 0 6px 16px rgb(15 23 42 / 18%); z-index: 2; }
.message-item:not(.is-own) .message-reaction-picker { right: -44px; }
.message-item.is-own .message-reaction-picker { left: -44px; }
.message-reaction-picker button { width: 30px; height: 30px; padding: 0; font-size: 1rem; background: transparent; border: 0; border-radius: 50%; cursor: pointer; }
.message-reaction-picker button:hover, .message-reaction-picker button.is-active { background: var(--color-bg-subtle); }
.message-reaction-picker__emoji-picker { position: absolute; top: 0; width: 320px; max-width: min(320px, 75vw); z-index: 3; }
.message-item:not(.is-own) .message-reaction-picker__emoji-picker { left: 42px; }
.message-item.is-own .message-reaction-picker__emoji-picker { right: 42px; }
@keyframes message-delete { from { max-height: 220px; opacity: 1; transform: scale(1); } to { max-height: 0; margin: 0; opacity: 0; transform: scale(.92); } }
@keyframes reaction-pop { 0% { transform: scale(.72); } 65% { transform: scale(1.18); } 100% { transform: scale(1); } }
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
