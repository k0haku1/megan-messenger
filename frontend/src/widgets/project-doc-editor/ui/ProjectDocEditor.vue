<template>
  <section class="project-doc-editor">
    <header class="project-doc-editor__header">
      <BaseTextField v-model="titleModel" label="Заголовок" :disabled="disabled" />
      <div class="project-doc-editor__actions">
        <BaseButton variant="ghost" :disabled="disabled || !docId" @click="emit('delete')">
          Удалить
        </BaseButton>
        <BaseButton :loading="saving" :disabled="disabled" @click="emit('save')">
          Сохранить
        </BaseButton>
      </div>
    </header>

    <div class="project-doc-editor__split">
      <label class="project-doc-editor__pane">
        <span>Markdown</span>
        <textarea
          :value="bodyMd"
          spellcheck="false"
          :disabled="disabled"
          @input="emit('update:bodyMd', ($event.target as HTMLTextAreaElement).value)"
        />
      </label>
      <div class="project-doc-editor__pane">
        <span>Просмотр</span>
        <article class="markdown-body" v-html="previewHtml" />
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { renderMarkdown } from '@/shared/lib/render-markdown'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'

const props = defineProps<{
  docId: string | null
  title: string
  bodyMd: string
  saving: boolean
  disabled: boolean
}>()

const emit = defineEmits<{
  'update:title': [value: string]
  'update:bodyMd': [value: string]
  save: []
  delete: []
}>()

const titleModel = computed({
  get: () => props.title,
  set: (value: string) => emit('update:title', value),
})

const previewHtml = computed(() => renderMarkdown(props.bodyMd || ''))
</script>

<style scoped>
.project-doc-editor {
  display: flex;
  min-width: 0;
  min-height: 0;
  flex: 1;
  flex-direction: column;
  overflow: hidden;
  background: var(--color-bg-chat);
}
.project-doc-editor__header {
  display: flex;
  flex-shrink: 0;
  align-items: flex-end;
  gap: 12px;
  padding: 14px 16px;
  background: var(--color-bg-elevated);
  border-bottom: 1px solid var(--color-border);
}
.project-doc-editor__header :deep(.text-field) {
  flex: 1;
}
.project-doc-editor__actions {
  display: flex;
  gap: 8px;
}
.project-doc-editor__split {
  display: grid;
  min-height: 0;
  flex: 1;
  overflow: hidden;
  grid-template-columns: 1fr 1fr;
}
.project-doc-editor__pane {
  display: flex;
  min-width: 0;
  min-height: 0;
  flex-direction: column;
  gap: 8px;
  padding: 14px 16px;
  overflow: hidden;
}
.project-doc-editor__pane > span {
  color: var(--color-text-secondary);
  font-size: var(--font-size-xs);
  font-weight: 600;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}
.project-doc-editor__pane textarea {
  min-height: 0;
  flex: 1;
  padding: 12px 14px;
  overflow-y: auto;
  resize: none;
  font: 14px/1.5 ui-monospace, SFMono-Regular, Menlo, monospace;
  color: inherit;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}
.project-doc-editor__pane article {
  min-height: 0;
  flex: 1;
  overflow: auto;
  padding: 12px 14px;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}
</style>
