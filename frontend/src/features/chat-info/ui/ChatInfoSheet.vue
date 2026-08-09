<template>
  <BaseModalSheet
    :open="open"
    :title="sheetTitle"
    title-id="chat-info-title"
    scrollable
    compact
    :hide-title="view === 'main'"
    @close="onClose"
  >
    <template #header-start>
      <BaseIconButton
        v-if="view !== 'main'"
        label="Назад"
        variant="ghost"
        @click="goBack()"
      >
        <AppIcon name="chevron-left" />
      </BaseIconButton>
    </template>

    <section v-if="view === 'main'" class="chat-info">
      <div class="chat-info__hero">
        <BaseAvatar
          class="chat-info__hero-avatar"
          :name="title"
          :src="avatarUrl"
          size="xl"
          :color="type === 'group' ? 2 : 0"
        />
        <h3 class="chat-info__hero-title">{{ title }}</h3>
        <p class="chat-info__hero-sub">
          {{ type === 'group' ? 'Групповой чат' : 'Личные сообщения' }}
        </p>
      </div>

      <div class="chat-info__group" role="list">
        <button
          v-if="type === 'group' && slug"
          type="button"
          class="chat-info__row"
          role="listitem"
          @click="view = 'invite'"
        >
          <span class="chat-info__row-icon is-accent">
            <AppIcon name="link" size="sm" />
          </span>
          <span class="chat-info__row-body">
            <strong>Пригласить в группу</strong>
            <small>Ссылка и код приглашения</small>
          </span>
          <AppIcon name="chevron-right" size="sm" class="chat-info__row-chevron" />
        </button>
        <button type="button" class="chat-info__row" role="listitem" @click="openMedia()">
          <span class="chat-info__row-icon is-accent">
            <AppIcon name="attach" size="sm" />
          </span>
          <span class="chat-info__row-body">
            <strong>Медиа и файлы</strong>
            <small>Фото, видео и документы</small>
          </span>
          <AppIcon name="chevron-right" size="sm" class="chat-info__row-chevron" />
        </button>
      </div>

      <div v-if="type === 'group'" class="chat-info__group">
        <button
          type="button"
          class="chat-info__row is-danger"
          :disabled="leaving"
          @click="leave()"
        >
          <span class="chat-info__row-icon is-danger">
            <AppIcon name="logout" size="sm" />
          </span>
          <span class="chat-info__row-body">
            <strong>{{ leaving ? 'Выходим…' : 'Покинуть группу' }}</strong>
          </span>
        </button>
        <p v-if="leaveError" class="chat-info__error">{{ leaveError }}</p>
      </div>
    </section>

    <section v-else-if="view === 'invite'" class="chat-info">
      <div class="chat-info__group chat-info__invite">
        <span class="chat-info__label">Код приглашения</span>
        <div class="chat-info__slug">
          <code>{{ slug }}</code>
          <BaseButton variant="ghost" @click="copySlug()">
            {{ copied ? 'Скопировано' : 'Копировать' }}
          </BaseButton>
        </div>
        <p class="chat-info__hint">
          Отправьте код коллегам — они смогут вступить через «Новый чат → Вступить по ссылке».
        </p>
      </div>
    </section>

    <section v-else class="chat-info chat-info--media">
      <div class="chat-info__filters" role="tablist">
        <button
          v-for="tab in mediaTabs"
          :key="tab.id"
          type="button"
          role="tab"
          :aria-selected="mediaKind === tab.id"
          :class="{ 'is-active': mediaKind === tab.id }"
          @click="setMediaKind(tab.id)"
        >
          {{ tab.label }}
        </button>
      </div>

      <div class="chat-info__media-stage">
        <div
          ref="mediaGridRef"
          class="chat-info__media-grid"
          :class="{ 'is-switching': mediaSwitching }"
          @scroll.passive="onMediaScroll"
        >
          <button
            v-for="item in mediaItems"
            :key="item.id"
            type="button"
            class="chat-info__media-cell"
            :aria-label="item.originalName"
            @click="openItem(item)"
          >
            <img
              v-if="(item.kind === 'image' || item.kind === 'video') && (item.thumbUrl || item.url)"
              :src="item.thumbUrl || item.url"
              :alt="item.originalName"
              loading="lazy"
            />
            <span v-else class="chat-info__file-cell">
              <strong>{{ item.originalName || 'Файл' }}</strong>
            </span>
            <span v-if="item.kind === 'video'" class="chat-info__media-badge">▶</span>
          </button>
        </div>

        <div
          v-if="mediaSwitching || (mediaLoading && mediaItems.length === 0)"
          class="chat-info__media-overlay"
          aria-hidden="true"
        />

        <p
          v-if="!mediaSwitching && !mediaLoading && mediaItems.length === 0"
          class="chat-info__media-empty"
        >
          Пока нет вложений
        </p>
        <p v-if="mediaError" class="chat-info__error">{{ mediaError }}</p>
      </div>

      <MediaPreviewOverlay :item="previewItem" @close="previewItem = null" />
    </section>
  </BaseModalSheet>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { leaveGroupChat } from '@/features/group-chat/model/group-chat.actions'
import { messageApi } from '@/entities/message/api/message.api'
import type { MessageAttachment } from '@/entities/message/model/types'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'
import { copyGroupInviteSlug } from '@/shared/lib/group-invite'
import AppIcon from '@/shared/ui/AppIcon.vue'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'
import MediaPreviewOverlay from '@/shared/ui/MediaPreviewOverlay.vue'

type ChatInfoView = 'main' | 'invite' | 'media'
type MediaKind = '' | 'image' | 'video' | 'file'

const props = defineProps<{
  open: boolean
  conversationId: string
  title: string
  type: 'dm' | 'group'
  slug?: string
  avatarUrl?: string
}>()

const emit = defineEmits<{ close: [] }>()

const view = ref<ChatInfoView>('main')
const copied = ref(false)
const leaving = ref(false)
const leaveError = ref('')

const mediaKind = ref<MediaKind>('')
const mediaItems = ref<MessageAttachment[]>([])
const mediaCursor = ref('')
const mediaLoading = ref(false)
const mediaSwitching = ref(false)
const mediaError = ref('')
const mediaGridRef = ref<HTMLElement | null>(null)
const previewItem = ref<MessageAttachment | null>(null)
let mediaRequestId = 0

const mediaTabs = [
  { id: '' as MediaKind, label: 'Все' },
  { id: 'image' as MediaKind, label: 'Фото' },
  { id: 'video' as MediaKind, label: 'Видео' },
  { id: 'file' as MediaKind, label: 'Файлы' },
]

const sheetTitle = computed(() => {
  if (view.value === 'invite') return 'Пригласить'
  if (view.value === 'media') return 'Медиа и файлы'
  return props.type === 'group' ? 'О группе' : 'О чате'
})

watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen) {
      view.value = 'main'
      copied.value = false
      leaveError.value = ''
      resetMedia()
    }
  },
)

function onClose() {
  previewItem.value = null
  emit('close')
}

function goBack() {
  previewItem.value = null
  view.value = 'main'
}

function resetMedia() {
  mediaRequestId += 1
  mediaItems.value = []
  mediaCursor.value = ''
  mediaError.value = ''
  mediaLoading.value = false
  mediaSwitching.value = false
  mediaKind.value = ''
  previewItem.value = null
}

async function openMedia() {
  view.value = 'media'
  resetMedia()
  await loadMedia(true)
}

async function setMediaKind(kind: MediaKind) {
  if (kind === mediaKind.value && !mediaLoading.value) return

  mediaKind.value = kind
  mediaCursor.value = ''
  mediaError.value = ''
  mediaSwitching.value = true
  if (mediaGridRef.value) mediaGridRef.value.scrollTop = 0
  await loadMedia(true)
}

async function loadMedia(replace: boolean) {
  if (!replace && mediaLoading.value) return
  if (!replace && !mediaCursor.value) return

  const requestId = ++mediaRequestId
  mediaLoading.value = true
  mediaError.value = ''

  try {
    const response = await messageApi.listMedia(props.conversationId, {
      limit: 36,
      cursor: replace ? undefined : mediaCursor.value || undefined,
      kind: mediaKind.value || undefined,
    })
    if (requestId !== mediaRequestId) return

    mediaItems.value = replace ? response.items : [...mediaItems.value, ...response.items]
    mediaCursor.value = response.nextCursor || ''
  } catch (err) {
    if (requestId !== mediaRequestId) return
    mediaError.value = getApiErrorMessage(err, { fallback: 'Не удалось загрузить вложения' })
    if (replace) mediaItems.value = []
  } finally {
    if (requestId === mediaRequestId) {
      mediaLoading.value = false
      mediaSwitching.value = false
    }
  }
}

function onMediaScroll() {
  const el = mediaGridRef.value
  if (!el || mediaLoading.value || mediaSwitching.value || !mediaCursor.value) return
  if (el.scrollTop + el.clientHeight >= el.scrollHeight - 80) {
    void loadMedia(false)
  }
}

function openItem(item: MessageAttachment) {
  previewItem.value = item
}

async function copySlug() {
  if (!props.slug) return
  await copyGroupInviteSlug(props.slug)
  copied.value = true
}

async function leave() {
  if (!window.confirm('Покинуть группу? Вы сможете вернуться по коду приглашения.')) return

  leaving.value = true
  leaveError.value = ''
  try {
    await leaveGroupChat(props.conversationId)
    emit('close')
  } catch (err) {
    leaveError.value = getApiErrorMessage(err, { fallback: 'Не удалось покинуть группу' })
  } finally {
    leaving.value = false
  }
}
</script>

<style scoped>
.chat-info {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.chat-info__hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 8px 8px 18px;
  text-align: center;
}
.chat-info__hero-avatar {
  box-shadow: 0 10px 28px rgb(15 23 42 / 14%);
}
.chat-info__hero-title {
  margin: 6px 0 0;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.03em;
  line-height: 1.2;
}
.chat-info__hero-sub {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.chat-info__group {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--color-bg-subtle);
  border-radius: 16px;
}
.chat-info__row {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 12px 14px;
  color: inherit;
  background: transparent;
  border: 0;
  border-bottom: 1px solid color-mix(in srgb, var(--color-border) 70%, transparent);
  cursor: pointer;
  font: inherit;
  text-align: left;
  transition: background 140ms ease;
}
.chat-info__group > .chat-info__row:last-child {
  border-bottom: 0;
}
.chat-info__row:hover {
  background: color-mix(in srgb, var(--color-bg-hover) 80%, transparent);
}
.chat-info__row:disabled {
  opacity: 0.65;
  cursor: wait;
}
.chat-info__row-icon {
  display: grid;
  place-items: center;
  width: 34px;
  height: 34px;
  flex: 0 0 auto;
  border-radius: 10px;
}
.chat-info__row-icon.is-accent {
  color: var(--color-accent);
  background: color-mix(in srgb, var(--color-accent) 14%, transparent);
}
.chat-info__row-icon.is-danger {
  color: var(--color-danger);
  background: color-mix(in srgb, var(--color-danger) 12%, transparent);
}
.chat-info__row-body {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
  gap: 2px;
}
.chat-info__row-body strong {
  font-size: 15px;
  font-weight: 600;
}
.chat-info__row-body small {
  color: var(--color-text-secondary);
  font-size: 12px;
}
.chat-info__row-chevron {
  color: var(--color-text-tertiary);
}
.chat-info__row.is-danger .chat-info__row-body strong {
  color: var(--color-danger);
}
.chat-info__label {
  display: block;
  margin-bottom: 8px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.chat-info__invite {
  padding: 14px;
}
.chat-info__slug {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  padding: 10px 12px;
  background: var(--color-bg-elevated);
  border-radius: var(--radius-sm);
}
.chat-info__slug code {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
}
.chat-info__hint {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-normal);
}
.chat-info__error {
  margin: 0;
  padding: 0 14px 12px;
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}
.chat-info__filters {
  display: flex;
  gap: 4px;
  padding: 4px;
  background: var(--color-bg-subtle);
  border-radius: 12px;
}
.chat-info__filters button {
  flex: 1;
  padding: 8px 10px;
  color: var(--color-text-secondary);
  background: transparent;
  border: 0;
  border-radius: 9px;
  cursor: pointer;
  font: inherit;
  font-size: 13px;
  font-weight: 600;
  transition:
    color 160ms ease,
    background 160ms ease;
}
.chat-info__filters button.is-active {
  color: var(--color-text-primary);
  background: var(--color-bg-elevated);
  box-shadow: 0 1px 3px rgb(15 23 42 / 8%);
}
.chat-info__media-stage {
  position: relative;
  min-height: min(52vh, 420px);
}
.chat-info__media-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  align-content: start;
  gap: 3px;
  height: min(52vh, 420px);
  overflow: auto;
  opacity: 1;
  transition: opacity 180ms ease;
  border-radius: 12px;
}
.chat-info__media-grid.is-switching {
  opacity: 0.35;
  pointer-events: none;
}
.chat-info__media-overlay {
  position: absolute;
  inset: 0;
  pointer-events: none;
  border-radius: 12px;
  background: linear-gradient(
    180deg,
    color-mix(in srgb, var(--color-bg-elevated) 8%, transparent),
    color-mix(in srgb, var(--color-bg-elevated) 18%, transparent)
  );
}
.chat-info__media-empty {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  pointer-events: none;
}
.chat-info__media-cell {
  position: relative;
  aspect-ratio: 1;
  padding: 0;
  overflow: hidden;
  background: var(--color-bg-subtle);
  border: 0;
  cursor: pointer;
}
.chat-info__media-cell img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.chat-info__file-cell {
  display: grid;
  place-items: center;
  height: 100%;
  padding: 8px;
  text-align: center;
  font-size: 11px;
}
.chat-info__media-badge {
  position: absolute;
  right: 6px;
  bottom: 6px;
  display: grid;
  place-items: center;
  width: 22px;
  height: 22px;
  color: #fff;
  background: rgb(0 0 0 / 55%);
  border-radius: 999px;
  font-size: 10px;
}
</style>
