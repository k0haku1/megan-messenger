<template>
  <section class="project-decision-list">
    <article v-for="decision in decisions" :key="decision.id" class="project-decision-list__item">
      <header>
        <strong>{{ decision.summary }}</strong>
        <select
          :value="decision.status"
          :disabled="pendingId === decision.id"
          @change="updateStatus(decision.id, ($event.target as HTMLSelectElement).value)"
        >
          <option value="proposed">Предложено</option>
          <option value="accepted">Принято</option>
          <option value="deprecated">Устарело</option>
        </select>
      </header>
      <p v-if="decision.context">{{ decision.context }}</p>
      <footer class="project-decision-list__meta">
        <time>{{
          new Date(decision.createdAt).toLocaleString('ru', {
            day: 'numeric',
            month: 'short',
            hour: '2-digit',
            minute: '2-digit',
          })
        }}</time>
        <span v-if="decision.messageId" class="project-decision-list__source">Из чата</span>
      </footer>
    </article>

    <div v-if="decisions.length === 0" class="project-decision-list__empty">
      Решений пока нет
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { projectApi } from '@/entities/project/api/project.api'
import type { ProjectDecision } from '@/entities/project/model/types'

const props = defineProps<{ projectId: string; decisions: ProjectDecision[] }>()
const emit = defineEmits<{ updated: [] }>()

const pendingId = ref<string | null>(null)

async function updateStatus(decisionId: string, status: string) {
  pendingId.value = decisionId
  try {
    await projectApi.updateDecisionStatus(props.projectId, decisionId, status)
    emit('updated')
  } finally {
    pendingId.value = null
  }
}
</script>

<style scoped>
.project-decision-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.project-decision-list__item {
  padding: 14px;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}
.project-decision-list__item header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 8px;
}
.project-decision-list__item p {
  margin: 0 0 8px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-normal);
  white-space: pre-wrap;
}
.project-decision-list__meta {
  display: flex;
  align-items: center;
  gap: 8px;
}
.project-decision-list__item time {
  color: var(--color-text-tertiary);
  font-size: var(--font-size-xs);
}
.project-decision-list__source {
  padding: 2px 8px;
  color: var(--color-accent);
  font-size: var(--font-size-xs);
  background: rgb(47 111 237 / 10%);
  border-radius: var(--radius-sm);
}
.project-decision-list__empty {
  padding: 20px 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
</style>
