<template>
  <aside class="project-doc-list">
    <header class="project-doc-list__header">
      <h3>Документы</h3>
      <BaseButton @click="emit('create')">Новый</BaseButton>
    </header>

    <div class="project-doc-list__scroll">
      <button
        v-for="doc in docs"
        :key="doc.id"
        class="project-doc-list__item"
        :class="{ 'is-active': doc.id === activeDocId }"
        type="button"
        @click="emit('select', doc.id)"
      >
        <strong>{{ doc.title }}</strong>
        <time>{{
          new Date(doc.updatedAt).toLocaleDateString('ru', { day: 'numeric', month: 'short' })
        }}</time>
      </button>

      <div v-if="docs.length === 0" class="project-doc-list__empty">Пока нет MD-документов</div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import type { ProjectDocSummary } from '@/entities/project/model/types'
import BaseButton from '@/shared/ui/BaseButton.vue'

defineProps<{
  docs: ProjectDocSummary[]
  activeDocId: string | null
}>()

const emit = defineEmits<{ select: [docId: string]; create: [] }>()
</script>

<style scoped>
.project-doc-list {
  display: flex;
  min-width: 240px;
  max-width: 280px;
  min-height: 0;
  flex-direction: column;
  overflow: hidden;
  border-right: 1px solid var(--color-border);
  background: var(--color-bg-elevated);
}
.project-doc-list__header {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 14px;
  border-bottom: 1px solid var(--color-border);
}
.project-doc-list__header h3 {
  margin: 0;
  font-size: var(--font-size-md);
}
.project-doc-list__scroll {
  min-height: 0;
  flex: 1;
  overflow-y: auto;
}
.project-doc-list__item {
  display: flex;
  width: 100%;
  flex-direction: column;
  gap: 4px;
  padding: 12px 14px;
  text-align: left;
  background: transparent;
  border: 0;
  border-bottom: 1px solid var(--color-border);
  cursor: pointer;
}
.project-doc-list__item.is-active {
  background: color-mix(in srgb, var(--color-accent) 10%, transparent);
}
.project-doc-list__item time {
  color: var(--color-text-tertiary);
  font-size: var(--font-size-xs);
}
.project-doc-list__empty {
  padding: 16px 14px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
</style>
