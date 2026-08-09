<template>
  <div v-if="attachments?.length" class="message-attachments">
    <template v-for="att in attachments" :key="att.id">
      <button
        v-if="att.kind === 'image' && (att.url || att.thumbUrl)"
        type="button"
        class="message-attachments__image"
        :aria-label="`Открыть ${att.originalName || 'фото'}`"
        @click="openPreview(att)"
      >
        <img :src="att.thumbUrl || att.url" :alt="att.originalName" loading="lazy" />
      </button>

      <div v-else-if="att.kind === 'video'" class="message-attachments__video">
        <button
          type="button"
          class="message-attachments__video-poster"
          :aria-label="`Открыть ${att.originalName || 'видео'}`"
          @click="openPreview(att)"
        >
          <img
            v-if="att.thumbUrl"
            :src="att.thumbUrl"
            :alt="att.originalName"
            loading="lazy"
          />
          <span v-else class="message-attachments__video-fallback">{{ att.originalName }}</span>
          <span class="message-attachments__play" aria-hidden="true">▶</span>
        </button>
      </div>

      <button
        v-else
        type="button"
        class="message-attachments__file"
        :aria-label="`Открыть ${att.originalName || 'файл'}`"
        @click="openPreview(att)"
      >
        <span class="message-attachments__file-icon" aria-hidden="true">📄</span>
        <span class="message-attachments__file-meta">
          <strong>{{ att.originalName || 'Файл' }}</strong>
          <small>{{ formatSize(att.sizeBytes) }}</small>
        </span>
      </button>
    </template>

    <MediaPreviewOverlay :item="previewItem" @close="previewItem = null" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { MessageAttachment } from '@/entities/message/model/types'
import MediaPreviewOverlay from '@/shared/ui/MediaPreviewOverlay.vue'

defineProps<{
  attachments?: MessageAttachment[]
}>()

const previewItem = ref<MessageAttachment | null>(null)

function openPreview(att: MessageAttachment) {
  previewItem.value = att
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}
</script>

<style scoped>
.message-attachments {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 4px;
}
.message-attachments__image {
  display: block;
  overflow: hidden;
  max-width: 280px;
  padding: 0;
  border: 0;
  border-radius: 8px;
  background: transparent;
  cursor: pointer;
}
.message-attachments__image img {
  display: block;
  width: 100%;
  max-height: 320px;
  object-fit: cover;
}
.message-attachments__video {
  overflow: hidden;
  border-radius: 8px;
  max-width: 280px;
  background: #000;
}
.message-attachments__video-poster {
  position: relative;
  display: block;
  width: 100%;
  min-height: 120px;
  padding: 0;
  border: 0;
  background: #111;
  cursor: pointer;
}
.message-attachments__video-poster img {
  display: block;
  width: 100%;
  max-height: 320px;
  object-fit: cover;
}
.message-attachments__video-fallback {
  display: block;
  padding: 40px 16px;
  color: #fff;
  font-size: 13px;
}
.message-attachments__play {
  position: absolute;
  inset: 50% auto auto 50%;
  display: grid;
  place-items: center;
  width: 44px;
  height: 44px;
  margin: -22px 0 0 -22px;
  color: #fff;
  background: rgb(0 0 0 / 55%);
  border-radius: 999px;
  font-size: 14px;
}
.message-attachments__file {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 8px 10px;
  color: inherit;
  text-align: left;
  background: var(--color-bg-subtle);
  border: 0;
  border-radius: 8px;
  cursor: pointer;
  font: inherit;
}
.message-attachments__file-meta {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 2px;
}
.message-attachments__file-meta strong {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 13px;
}
.message-attachments__file-meta small {
  color: var(--color-text-secondary);
  font-size: 11px;
}
</style>
