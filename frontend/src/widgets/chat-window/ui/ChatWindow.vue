<template>
  <section
    class="chat-window"
    :class="{ 'is-embedded': embedded, 'is-pane-active': active }"
    aria-label="Переписка"
  >
    <template v-if="header">
      <header
        class="chat-header"
        :class="{ 'is-draggable': rearrangeable }"
        :data-tile-drag-handle="rearrangeable ? '' : undefined"
        @pointerdown="onHeaderPointerDown"
      >
        <span v-if="rearrangeable" class="chat-header__drag-hint" aria-hidden="true" />
        <template v-if="selectionMode">
          <div class="chat-header__selection" data-no-pane-drag>
            <button
              type="button"
              class="chat-header__selection-action"
              :disabled="selectedCount === 0"
              @click="openForwardSelected()"
            >
              Переслать {{ selectedCount }}
            </button>
            <button
              type="button"
              class="chat-header__selection-action is-danger"
              :disabled="selectedCount === 0 || deletingSelected"
              @click="deleteSelected()"
            >
              Удалить {{ selectedCount }}
            </button>
          </div>
          <div class="chat-header__meta">
            <strong>{{ selectedCount }} выбрано</strong>
          </div>
          <BaseIconButton label="Отменить выделение" data-no-pane-drag @click="exitSelection()">
            <AppIcon name="close" />
          </BaseIconButton>
        </template>
        <template v-else>
          <BaseAvatar
            :name="header.avatarName"
            :src="header.avatarUrl"
            :color="header.type === 'group' ? 2 : 0"
          />
          <div class="chat-header__meta">
            <strong>{{ header.title }}</strong>
            <span>{{ header.type === 'group' ? 'групповой чат' : 'был(а) недавно' }}</span>
          </div>
          <div class="chat-header__actions" data-no-pane-drag>
            <BaseIconButton
              v-if="rearrangeable"
              label="Закрыть чат"
              @click.stop="emit('close')"
            >
              <AppIcon name="close" />
            </BaseIconButton>
            <BaseIconButton label="Поиск по сообщениям">
              <AppIcon name="search" />
            </BaseIconButton>
            <BaseIconButton
              label="Информация о чате"
              @click="chatInfoOpen = true"
            >
              <AppIcon name="menu" />
            </BaseIconButton>
          </div>
        </template>
      </header>

      <ChatInfoSheet
        v-if="activeConversationId && header"
        :open="chatInfoOpen"
        :conversation-id="activeConversationId"
        :title="header.title"
        :type="header.type"
        :slug="conversation?.slug"
        :avatar-url="header.avatarUrl"
        @close="chatInfoOpen = false"
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
        :messages="forwardMessages"
        :current-conversation-id="activeConversationId"
        @close="closeForward()"
        @forwarded="onForwarded"
      />

      <DeleteMessagesModal
        :open="deleteOpen"
        :count="deleteMessages.length"
        :can-delete-for-everyone="deleteCanForEveryone"
        :pending="deletingSelected"
        @close="closeDelete()"
        @confirm="onDeleteConfirm"
      />

      <MessageContextMenu
        :open="menuOpen"
        :anchor="menuAnchor"
        :items="menuItems"
        :reaction-details="menuReactionDetails"
        @close="closeMenu()"
        @select="handleMenuAction"
      />

      <div class="message-area-wrap">
        <div
          ref="messageAreaRef"
          class="message-area"
          :style="{ paddingBottom: `${messageAreaBottomPad}px` }"
          @scroll.passive="onContainerScroll"
        >
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
              :class="{
                'is-own': isOwnMessage(message, currentUserId),
                'is-deleting': deletingMessageIds.has(message.id),
                'has-avatar-col': isGroupChat && !isOwnMessage(message, currentUserId),
                'is-selecting': selectionMode,
                'is-selected': selectionMode && selectedIds.has(message.id),
              }"
            >
              <template v-if="isGroupChat && !isOwnMessage(message, currentUserId)">
                <button
                  v-if="shouldShowSenderAvatar(message, index, visibleMessages, currentUserId)"
                  type="button"
                  class="message-item__avatar-btn"
                  :aria-label="`Профиль @${message.sender.username}`"
                  @click="selectionMode ? toggleSelect(message.id) : openSenderProfile(message.sender.username)"
                >
                  <BaseAvatar
                    :name="message.sender.username"
                    :src="message.sender.avatarUrl"
                    size="sm"
                    :color="1"
                  />
                </button>
                <span v-else class="message-item__avatar-spacer" aria-hidden="true" />
              </template>

              <div class="message-item__body">
                <span
                  v-if="shouldShowSenderName(message, index, visibleMessages, currentUserId)"
                  class="message-item__sender"
                >
                  {{ message.sender.username }}
                </span>
                <article
                  class="message-bubble"
                  :class="{ 'is-context-pressed': contextPressedMessageId === message.id }"
                  @contextmenu.prevent="openMenu($event, message)"
                  @click="selectionMode && toggleSelect(message.id)"
                  @mouseenter="!selectionMode && scheduleReactionPicker(message.id)"
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
                  <MessageAttachments :attachments="message.attachments" />
                  <p v-if="message.content">{{ message.content }}</p>
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

              <button
                v-if="selectionMode"
                type="button"
                class="message-item__check"
                :class="{ 'is-checked': selectedIds.has(message.id) }"
                :aria-label="selectedIds.has(message.id) ? 'Снять выделение' : 'Выделить сообщение'"
                :aria-pressed="selectedIds.has(message.id)"
                @click="toggleSelect(message.id)"
              >
                <AppIcon v-if="selectedIds.has(message.id)" name="check" size="sm" />
              </button>
            </div>
          </div>
        </div>

        <button
          v-if="showJumpToLatest"
          type="button"
          class="jump-to-latest"
          :style="{ bottom: `${jumpBottomOffset}px` }"
          :aria-label="unseenCount > 0 ? `К последнему сообщению, новых: ${unseenCount}` : 'К последнему сообщению'"
          @click="scrollToLatest('smooth')"
        >
          <AppIcon name="chevron-right" class="jump-to-latest__icon" />
          <span v-if="unseenCount > 0" class="jump-to-latest__badge">
            {{ unseenCount > 99 ? '99+' : unseenCount }}
          </span>
        </button>

        <div v-show="!selectionMode" ref="composerWrapRef" class="composer-wrap">
        <div v-if="replyTo" class="composer-reply">
          <div class="composer-reply__meta">
            <strong>Ответ для {{ replyTo.sender.username }}</strong>
            <BaseIconButton label="Отменить ответ" @click="clearReply()">
              <AppIcon name="close" size="sm" />
            </BaseIconButton>
          </div>
          <p>{{ replyTo.content || 'Вложение' }}</p>
        </div>

        <TransitionGroup
          name="composer-attach"
          tag="div"
          class="composer-attachments"
        >
          <div
            v-for="item in pendingAttachments"
            :key="item.localId"
            class="composer-attachments__item"
            :class="{ 'is-error': Boolean(item.error), 'is-uploading': item.uploading }"
          >
            <img v-if="item.previewUrl" :src="item.previewUrl" alt="" />
            <span v-else class="composer-attachments__file">{{ item.file.name }}</span>
            <span v-if="item.uploading" class="composer-attachments__spinner" aria-hidden="true" />
            <button type="button" :aria-label="`Убрать ${item.file.name}`" @click="removeAttachment(item.localId)">
              <AppIcon name="close" size="sm" />
            </button>
          </div>
        </TransitionGroup>

        <footer class="composer">
          <BaseIconButton
            label="Прикрепить файл"
            :disabled="!activeConversationId || isSending"
            @click="fileInputRef?.click()"
          >
            <AppIcon name="attach" />
          </BaseIconButton>
          <input
            ref="fileInputRef"
            type="file"
            class="composer-file-input"
            multiple
            accept="image/*,video/*,.pdf,.doc,.docx,.zip,.txt"
            @change="onAttachFiles"
          />
          <textarea
            ref="composerRef"
            v-model="draft"
            rows="1"
            :placeholder="replyTo ? 'Ответ...' : 'Сообщение'"
            aria-label="Сообщение"
            @focus="onComposerFocus"
            @keydown.enter.exact.prevent="sendAndKeepFocus()"
          />
          <BaseIconButton label="Эмодзи">
            <AppIcon name="emoji" />
          </BaseIconButton>
          <BaseIconButton label="Отправить" variant="primary" :disabled="!canSend" @click="sendAndKeepFocus()">
            <AppIcon name="send" />
          </BaseIconButton>
        </footer>

        <p v-if="sendError" class="composer-error">{{ sendError }}</p>
        <p v-if="actionNotice" class="composer-notice">{{ actionNotice }}</p>
      </div>
      <p
        v-if="selectionMode && actionNotice"
        class="composer-notice composer-notice--selection"
      >{{ actionNotice }}</p>
      </div>
    </template>

    <div v-else-if="!embedded" class="chat-placeholder">
      <div class="chat-placeholder__mark">M</div>
      <strong>Выберите чат</strong>
      <span>Откройте диалог слева, чтобы начать общение</span>
    </div>
  </section>
</template>

<script setup lang="ts">
import 'emoji-picker-element'
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useResizeObserver } from '@vueuse/core'
import CaptureDecisionModal from '@/features/capture-decision/ui/CaptureDecisionModal.vue'
import DeleteMessagesModal from '@/features/delete-messages/ui/DeleteMessagesModal.vue'
import ForwardMessageModal from '@/features/forward-message/ui/ForwardMessageModal.vue'
import ChatInfoSheet from '@/features/chat-info/ui/ChatInfoSheet.vue'
import MessageContextMenu, {
  type MessageContextMenuAnchor,
  type MessageContextMenuItem,
  type MessageContextMenuReactionDetail,
} from '@/features/message-context-menu/ui/MessageContextMenu.vue'
import { storeToRefs } from 'pinia'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import {
  getConversationAvatarName,
  getConversationAvatarUrl,
  getConversationTitle,
} from '@/entities/conversation/lib/display'
import { useConversationProjects } from '@/entities/conversation/lib/use-conversation-projects'
import { useConversation } from '@/entities/conversation/lib/use-conversation'
import { hideMessageForMe } from '@/entities/message/api/hidden-message.repository'
import { messageApi } from '@/entities/message/api/message.api'
import { removeMessage, upsertMessage } from '@/entities/message/api/message.repository'
import type { Message } from '@/entities/message/model/types'
import MessageAttachments from '@/entities/message/ui/MessageAttachments.vue'
import {
  formatMessageTime,
  isOwnMessage,
  shouldShowSenderAvatar,
  shouldShowSenderName,
} from '@/entities/message/lib/display'
import { useConversationMessages } from '@/entities/message/lib/use-conversation-messages'
import { useMessageComposer } from '@/features/compose-message/model/use-message-composer'
import { useMessagesSync } from '@/features/message-sync/model/use-messages-sync'
import { useConversationWs } from '@/features/conversation-ws/model/use-conversation-ws'
import { useScrollToLatest } from '@/features/scroll-to-latest/model/use-scroll-to-latest'
import { useSessionStore } from '@/entities/session/model/session.store'
import { useUserProfileStore } from '@/features/user-profile/model/user-profile.store'
import { queryClient } from '@/shared/api/query-client'
import { queryKeys } from '@/shared/api/query-keys'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'

const props = withDefaults(
  defineProps<{
    conversationId?: string | null
    active?: boolean
    embedded?: boolean
    rearrangeable?: boolean
  }>(),
  {
    conversationId: undefined,
    active: true,
    embedded: false,
    rearrangeable: false,
  },
)

const emit = defineEmits<{
  close: []
  headerDrag: [event: PointerEvent]
}>()

function onHeaderPointerDown(event: PointerEvent): void {
  if (!props.rearrangeable) return
  emit('headerDrag', event)
}

const navigation = useConversationSelectionStore()
const session = useSessionStore()
const userProfile = useUserProfileStore()
const { selectedConversationId, pendingPeer } = storeToRefs(navigation)
const activeConversationId = computed(() => props.conversationId ?? selectedConversationId.value)
const composerPendingPeer = computed(() =>
  props.conversationId ? null : pendingPeer.value,
)
const composerRef = ref<HTMLTextAreaElement | null>(null)
const fileInputRef = ref<HTMLInputElement | null>(null)
const composerWrapRef = ref<HTMLElement | null>(null)
const composerHeight = ref(72)
const messageAreaBottomPad = computed(() =>
  selectionMode.value ? 24 : Math.max(composerHeight.value + 20, 88),
)
const jumpBottomOffset = computed(() =>
  selectionMode.value ? 24 : Math.max(composerHeight.value + 22, 40),
)

useResizeObserver(composerWrapRef, (entries) => {
  const height = entries[0]?.contentRect.height
  if (typeof height === 'number' && height > 0) {
    composerHeight.value = Math.ceil(height)
  }
})

const chatInfoOpen = ref(false)
const captureOpen = ref(false)
const captureMessage = ref<Message | null>(null)
const forwardOpen = ref(false)
const forwardMessages = ref<Message[]>([])
const deleteOpen = ref(false)
const deleteMessages = ref<Message[]>([])
const actionNotice = ref('')
const replyTo = ref<Message | null>(null)
const selectionMode = ref(false)
const selectedIds = ref(new Set<string>())
const deletingSelected = ref(false)

const selectedCount = computed(() => selectedIds.value.size)
const selectedMessages = computed(() =>
  visibleMessages.value.filter((message) => selectedIds.value.has(message.id)),
)
const deleteCanForEveryone = computed(() => {
  const userId = currentUserId.value
  if (!userId || deleteMessages.value.length === 0) return false
  return deleteMessages.value.every((message) => isOwnMessage(message, userId) && !message.deletedAt)
})

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

const conversation = useConversation(activeConversationId)
const visibleMessages = useConversationMessages(activeConversationId)
const currentUserId = computed(() => session.user?.id)
const messagesById = computed(() => new Map(visibleMessages.value.map((message) => [message.id, message])))

const messageAreaRef = ref<HTMLElement | null>(null)
const latestMessageId = computed(() => visibleMessages.value.at(-1)?.id ?? null)
const latestIsOwn = computed(() => {
  const latest = visibleMessages.value.at(-1)
  if (!latest) return false
  return isOwnMessage(latest, currentUserId.value)
})

const { showJumpToLatest, unseenCount, onContainerScroll, scrollToLatest } = useScrollToLatest({
  container: messageAreaRef,
  conversationKey: activeConversationId,
  latestMessageId,
  latestIsOwn,
})

const isGroupChat = computed(() => conversation.value?.type === 'group')
const linkedProjectsQuery = useConversationProjects(activeConversationId, isGroupChat)
const linkedProjects = computed(() => linkedProjectsQuery.data.value ?? [])
const canCaptureDecision = computed(() => isGroupChat.value && linkedProjects.value.length > 0)

useMessagesSync(activeConversationId)
useConversationWs(activeConversationId)

const header = computed(() => {
  if (conversation.value) {
    return {
      type: conversation.value.type,
      title: getConversationTitle(conversation.value),
      avatarName: getConversationAvatarName(conversation.value),
      avatarUrl: getConversationAvatarUrl(conversation.value),
      isPending: false,
    }
  }

  if (!props.conversationId && pendingPeer.value) {
    return {
      type: 'dm' as const,
      title: pendingPeer.value.username,
      avatarName: pendingPeer.value.username,
      avatarUrl: pendingPeer.value.avatarUrl,
      isPending: true,
    }
  }

  return null
})

const menuItems = computed<MessageContextMenuItem[]>(() => {
  const items: MessageContextMenuItem[] = [
    { id: 'reply', label: 'Ответить' },
    { id: 'forward', label: 'Переслать' },
    { id: 'select', label: 'Выделить' },
  ]

  if (canCaptureDecision.value) {
    items.push({ id: 'capture', label: 'Зафиксировать решение' })
  }

  items.push({ id: 'delete', label: 'Удалить', danger: true })

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

const {
  draft,
  isSending,
  sendError,
  canSend,
  send,
  clearReply,
  pendingAttachments,
  addFiles,
  removeAttachment,
} = useMessageComposer({
  selectedConversationId: activeConversationId,
  pendingPeer: composerPendingPeer,
  replyTo,
  onConversationOpened: (conversationId) => navigation.select(conversationId),
  onSent: () => {
    void scrollToLatest('smooth')
    focusComposer()
  },
})

async function sendAndKeepFocus() {
  await send()
  focusComposer()
}

function focusComposer() {
  void nextTick(() => {
    composerRef.value?.focus({ preventScroll: true })
  })
}

function onAttachFiles(event: Event) {
  const input = event.target as HTMLInputElement
  if (!input.files?.length) return
  void addFiles(input.files)
  input.value = ''
  focusComposer()
}

function onComposerFocus(): void {
  if (!activeConversationId.value) return
  navigation.focusConversation(activeConversationId.value)
}

function openSenderProfile(username: string) {
  void userProfile.open(username)
}

watch(
  () => props.active,
  (isActive) => {
    if (!isActive) return
    focusComposer()
  },
  { immediate: true },
)

watch(activeConversationId, (id) => {
  exitSelection()
  if (!id || !props.active) return
  focusComposer()
})

function showNotice(text: string) {
  actionNotice.value = text
  window.setTimeout(() => {
    actionNotice.value = ''
  }, 4000)
}

function enterSelection(messageId: string) {
  selectionMode.value = true
  selectedIds.value = new Set([messageId])
  reactionPickerMessageId.value = null
  emojiPickerMessageId.value = null
}

function exitSelection() {
  selectionMode.value = false
  selectedIds.value = new Set()
  deletingSelected.value = false
}

function toggleSelect(messageId: string) {
  const next = new Set(selectedIds.value)
  if (next.has(messageId)) next.delete(messageId)
  else next.add(messageId)
  selectedIds.value = next
  if (next.size === 0) exitSelection()
}

function openForwardSelected() {
  const messages = selectedMessages.value
  if (!messages.length) return
  forwardMessages.value = messages
  forwardOpen.value = true
}

function openDeleteConfirm(messages: Message[]) {
  if (!messages.length) return
  deleteMessages.value = messages
  deleteOpen.value = true
}

function closeDelete() {
  if (deletingSelected.value) return
  deleteOpen.value = false
  deleteMessages.value = []
}

function deleteSelected() {
  openDeleteConfirm(selectedMessages.value)
}

async function onDeleteConfirm({ forEveryone }: { forEveryone: boolean }) {
  const conversationId = activeConversationId.value
  const messages = deleteMessages.value
  if (!conversationId || !messages.length || deletingSelected.value) return

  deletingSelected.value = true
  try {
    deletingMessageIds.value = new Set([
      ...deletingMessageIds.value,
      ...messages.map((message) => message.id),
    ])
    await Promise.all([
      Promise.all(
        messages.map((message) =>
          forEveryone
            ? messageApi.deleteForEveryone(conversationId, message.id)
            : messageApi.hideForMe(conversationId, message.id),
        ),
      ),
      delay(260),
    ])
    if (forEveryone) {
      await Promise.all(messages.map((message) => removeMessage(message.id)))
    } else {
      await Promise.all(messages.map((message) => hideMessageForMe(message.id)))
    }
    await queryClient.invalidateQueries({ queryKey: queryKeys.conversations })
    deleteOpen.value = false
    deleteMessages.value = []
    exitSelection()
  } catch {
    for (const message of messages) removeDeletingState(message.id)
    showNotice(forEveryone ? 'Не удалось удалить сообщения для всех' : 'Не удалось удалить сообщения')
  } finally {
    deletingSelected.value = false
  }
}

function onSelectionKeyDown(event: KeyboardEvent) {
  if (event.key === 'Escape' && selectionMode.value) {
    exitSelection()
  }
}

onMounted(() => {
  document.addEventListener('keydown', onSelectionKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', onSelectionKeyDown)
})

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
      exitSelection()
      replyTo.value = message
      break
    case 'forward':
      forwardMessages.value = [message]
      forwardOpen.value = true
      break
    case 'select':
      enterSelection(message.id)
      break
    case 'capture':
      captureMessage.value = message
      captureOpen.value = true
      break
    case 'delete':
      if (isOwnMessage(message, currentUserId.value)) {
        openDeleteConfirm([message])
      } else {
        await deleteForMe(message)
      }
      break
  }
}

function repliedMessage(message: Message): Message | undefined {
  return message.replyToId ? messagesById.value.get(message.replyToId) : undefined
}

function replyPreview(message: Message): string {
  if (message.deletedAt) return 'Сообщение удалено'
  if (message.content) return message.content
  const first = message.attachments?.[0]
  if (!first) return 'Сообщение'
  if (first.kind === 'image') return 'Фото'
  if (first.kind === 'video') return 'Видео'
  return first.originalName || 'Файл'
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
  const conversationId = activeConversationId.value
  if (!conversationId) return
  try {
    deletingMessageIds.value = new Set([...deletingMessageIds.value, message.id])
    await Promise.all([messageApi.hideForMe(conversationId, message.id), delay(260)])
    await hideMessageForMe(message.id)
    await queryClient.invalidateQueries({ queryKey: queryKeys.conversations })
  } catch {
    removeDeletingState(message.id)
    showNotice('Не удалось удалить сообщение')
  }
}

async function toggleReaction(message: Message, emoji: string, active: boolean) {
  const conversationId = activeConversationId.value
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
  forwardMessages.value = []
}

function onDecisionSaved(projectId: string) {
  const project = linkedProjects.value.find((item) => item.id === projectId)
  showNotice(project ? `Решение сохранено в «${project.name}»` : 'Решение сохранено')
}

function onForwarded() {
  exitSelection()
}
</script>

<style scoped>
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
.composer-wrap {
  position: absolute;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 3;
  width: min(calc(100% - 32px), 780px);
  margin: 0 auto 14px;
  padding-top: 8px;
  background: transparent;
  pointer-events: none;
}
.composer-wrap > * {
  pointer-events: auto;
}
.composer-wrap :deep(.composer) {
  width: 100%;
  margin: 0;
}
.composer-reply {
  margin-bottom: 8px;
  padding: 8px 10px;
  background: color-mix(in srgb, var(--color-bg-elevated) 82%, transparent);
  border-left: 3px solid var(--color-accent);
  border-radius: var(--radius-sm);
  backdrop-filter: blur(8px);
  box-shadow: 0 1px 6px rgb(15 23 42 / 8%);
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
.composer-attachments {
  display: flex;
  gap: 8px;
  margin: 0 0 8px;
  overflow-x: auto;
  padding: 2px 0;
}
.composer-attachments:empty {
  display: none;
  margin: 0;
  padding: 0;
}
.composer-attachments__item {
  position: relative;
  flex: 0 0 auto;
  width: 64px;
  height: 64px;
  overflow: hidden;
  background: var(--color-bg-subtle);
  border-radius: 8px;
  box-shadow: 0 1px 4px rgb(15 23 42 / 8%);
  transform-origin: center bottom;
}
.composer-attachments__item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.composer-attachments__file {
  display: grid;
  place-items: center;
  height: 100%;
  padding: 6px;
  font-size: 10px;
  text-align: center;
  word-break: break-word;
}
.composer-attachments__item button {
  position: absolute;
  top: 2px;
  right: 2px;
  z-index: 2;
  display: grid;
  place-items: center;
  width: 20px;
  height: 20px;
  padding: 0;
  color: #fff;
  background: rgb(0 0 0 / 55%);
  border: 0;
  border-radius: 999px;
  cursor: pointer;
  transition: transform 140ms ease, background 140ms ease;
}
.composer-attachments__item button :deep(.app-icon) {
  width: 12px;
  height: 12px;
}
.composer-attachments__item button:hover {
  background: rgb(0 0 0 / 72%);
  transform: scale(1.06);
}
.composer-attachments__item.is-uploading {
  opacity: 0.78;
}
.composer-attachments__item.is-uploading::after {
  position: absolute;
  inset: 0;
  z-index: 1;
  pointer-events: none;
  content: "";
  background: linear-gradient(
    110deg,
    transparent 30%,
    rgb(255 255 255 / 28%) 48%,
    transparent 66%
  );
  background-size: 200% 100%;
  animation: composer-attach-shimmer 1.1s ease-in-out infinite;
}
.composer-attachments__spinner {
  position: absolute;
  inset: auto 6px 6px auto;
  z-index: 2;
  width: 14px;
  height: 14px;
  border: 2px solid rgb(255 255 255 / 35%);
  border-top-color: #fff;
  border-radius: 999px;
  animation: composer-attach-spin 0.7s linear infinite;
}
.composer-attachments__item.is-error {
  outline: 1px solid var(--color-danger);
}
.composer-attach-enter-active {
  transition:
    opacity 220ms ease,
    transform 280ms cubic-bezier(0.22, 1, 0.36, 1);
}
.composer-attach-leave-active {
  transition:
    opacity 160ms ease,
    transform 180ms ease;
  position: absolute;
}
.composer-attach-enter-from {
  opacity: 0;
  transform: translateY(10px) scale(0.86);
}
.composer-attach-leave-to {
  opacity: 0;
  transform: translateY(4px) scale(0.9);
}
.composer-attach-move {
  transition: transform 220ms cubic-bezier(0.22, 1, 0.36, 1);
}
@keyframes composer-attach-shimmer {
  from { background-position: 120% 0; }
  to { background-position: -40% 0; }
}
@keyframes composer-attach-spin {
  to { transform: rotate(360deg); }
}
.composer-file-input {
  display: none;
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
.composer-notice--selection {
  position: absolute;
  right: 16px;
  bottom: 16px;
  left: 16px;
  z-index: 3;
  margin: 0;
  text-align: center;
}
</style>
