<template>
  <AppShell>
    <section class="projects-page">
      <header class="projects-page__header">
        <div>
          <h1>Проекты</h1>
          <p>База знаний команды: MD-документы, решения и привязка чатов</p>
        </div>
        <BaseButton @click="createOpen = true">Создать проект</BaseButton>
      </header>

      <div v-if="isLoading" class="projects-page__state">Загрузка…</div>
      <div v-else-if="isError" class="projects-page__state is-error">Не удалось загрузить проекты</div>

      <div v-else class="projects-page__grid">
        <RouterLink
          v-for="project in projects"
          :key="project.id"
          class="project-card"
          :to="{ name: 'project-detail', params: { projectId: project.id } }"
        >
          <strong>{{ project.name }}</strong>
          <span v-if="project.description">{{ project.description }}</span>
          <small>{{ project.docCount ?? 0 }} док. · {{ project.memberCount ?? 0 }} участн.</small>
        </RouterLink>

        <div v-if="projects.length === 0" class="projects-page__empty">
          <AppIcon class="projects-page__empty-icon" name="book" />
          <strong>Пока нет проектов</strong>
          <span>Создайте пространство для документов и решений команды</span>
        </div>
      </div>
    </section>

    <CreateProjectModal :open="createOpen" @close="createOpen = false" />
  </AppShell>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useProjects } from '@/entities/project/api/project.queries'
import CreateProjectModal from '@/features/create-project/ui/CreateProjectModal.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import AppShell from '@/shared/ui/AppShell.vue'
import BaseButton from '@/shared/ui/BaseButton.vue'

const { projects, isLoading, isError } = useProjects()
const createOpen = ref(false)
</script>

<style scoped>
.projects-page {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
  overflow: auto;
  background: var(--color-bg-elevated);
}
.projects-page__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 18px 20px;
  border-bottom: 1px solid var(--color-border);
}
.projects-page__header h1 {
  margin: 0;
  font-size: var(--font-size-2xl);
}
.projects-page__header p {
  margin: 6px 0 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.projects-page__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
  padding: 16px 20px 24px;
}
.project-card {
  display: flex;
  min-height: 120px;
  flex-direction: column;
  gap: 6px;
  padding: 14px;
  text-decoration: none;
  color: inherit;
  background: var(--color-bg-subtle);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: background var(--transition-fast);
}
.project-card:hover {
  background: var(--color-bg-hover);
}
.project-card span {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-tight);
}
.project-card small {
  margin-top: auto;
  color: var(--color-text-tertiary);
  font-size: var(--font-size-xs);
}
.projects-page__state {
  padding: 24px 20px;
  color: var(--color-text-secondary);
}
.projects-page__state.is-error {
  color: var(--color-danger);
}
.projects-page__empty {
  grid-column: 1 / -1;
  display: flex;
  min-height: 280px;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  color: var(--color-text-secondary);
  text-align: center;
}
.projects-page__empty strong {
  margin-top: 10px;
  color: var(--color-text-primary);
}
.projects-page__empty span {
  margin-top: 4px;
  font-size: var(--font-size-sm);
}
.projects-page__empty-icon {
  width: 40px;
  height: 40px;
  color: var(--color-accent);
}
</style>
