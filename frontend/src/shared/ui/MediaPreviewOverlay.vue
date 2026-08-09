<template>
  <Teleport to="body">
    <div
      v-if="item"
      class="media-preview"
      role="dialog"
      aria-modal="true"
      :aria-label="item.originalName || 'Превью'"
      @click.self="emit('close')"
    >
      <header class="media-preview__header">
        <div class="media-preview__meta">
          <strong>{{ item.originalName || 'Файл' }}</strong>
          <span v-if="item.kind !== 'file'">{{ kindLabel }}</span>
        </div>
        <BaseIconButton label="Закрыть" variant="ghost" class="media-preview__close" @click="emit('close')">
          <AppIcon name="close" />
        </BaseIconButton>
      </header>

      <div class="media-preview__body">
        <img
          v-if="item.kind === 'image' && previewUrl"
          :src="previewUrl"
          :alt="item.originalName"
          class="media-preview__image"
        />
        <video
          v-else-if="item.kind === 'video' && previewUrl"
          :src="previewUrl"
          class="media-preview__video"
          controls
          autoplay
          playsinline
        />
        <div v-else class="media-preview__file">
          <span class="media-preview__file-icon" aria-hidden="true">📄</span>
          <strong>{{ item.originalName || 'Файл' }}</strong>
          <small v-if="item.sizeBytes">{{ formatSize(item.sizeBytes) }}</small>
          <a
            v-if="previewUrl"
            class="media-preview__download"
            :href="previewUrl"
            :download="item.originalName"
            target="_blank"
            rel="noopener"
          >
            Открыть / скачать
          </a>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, onScopeDispose, watch } from 'vue'
import type { MessageAttachment } from '@/entities/message/model/types'
import AppIcon from '@/shared/ui/AppIcon.vue'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'

const props = defineProps<{
  item: MessageAttachment | null
}>()

const emit = defineEmits<{ close: [] }>()

const previewUrl = computed(() => props.item?.url || props.item?.thumbUrl || '')

const kindLabel = computed(() => {
  switch (props.item?.kind) {
    case 'image':
      return 'Фото'
    case 'video':
      return 'Видео'
    default:
      return 'Файл'
  }
})

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') emit('close')
}

watch(
  () => props.item,
  (item) => {
    if (item) {
      document.addEventListener('keydown', onKeydown)
      document.body.style.overflow = 'hidden'
    } else {
      document.removeEventListener('keydown', onKeydown)
      document.body.style.overflow = ''
    }
  },
)

onScopeDispose(() => {
  document.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
})
</script>

<style scoped>
.media-preview {
  position: fixed;
  inset: 0;
  z-index: 80;
  display: flex;
  flex-direction: column;
  background: rgb(8 12 18 / 92%);
  backdrop-filter: blur(8px);
}
.media-preview__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 14px;
  color: #fff;
}
.media-preview__meta {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 2px;
}
.media-preview__meta strong {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
}
.media-preview__meta span {
  color: rgb(255 255 255 / 65%);
  font-size: 12px;
}
.media-preview__close {
  color: #fff;
  background: rgb(255 255 255 / 12%);
}
.media-preview__close:hover {
  background: rgb(255 255 255 / 20%);
}
.media-preview__body {
  display: grid;
  flex: 1;
  place-items: center;
  min-height: 0;
  padding: 8px 16px 24px;
}
.media-preview__image,
.media-preview__video {
  max-width: min(100%, 1100px);
  max-height: 100%;
  object-fit: contain;
  border-radius: 8px;
}
.media-preview__video {
  width: min(100%, 960px);
  background: #000;
}
.media-preview__file {
  display: flex;
  max-width: 360px;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 28px 24px;
  color: #fff;
  text-align: center;
  background: rgb(255 255 255 / 8%);
  border-radius: 16px;
}
.media-preview__file-icon {
  font-size: 36px;
}
.media-preview__file small {
  color: rgb(255 255 255 / 65%);
}
.media-preview__download {
  margin-top: 8px;
  padding: 8px 14px;
  color: #fff;
  text-decoration: none;
  background: var(--color-accent, #3b82f6);
  border-radius: 999px;
  font-size: 13px;
  font-weight: 600;
}
</style>
