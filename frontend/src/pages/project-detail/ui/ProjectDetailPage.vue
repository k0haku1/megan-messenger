<template>
  <AppShell>
    <section class="project-detail-page">
      <header class="project-detail-page__header">
        <div class="project-detail-page__title">
          <RouterLink class="project-detail-page__back" :to="{ name: 'projects' }">
            <AppIcon name="chevron-left" size="sm" />
          </RouterLink>
          <div>
            <h1>{{ projectQuery.data.value?.name ?? 'Проект' }}</h1>
            <p v-if="projectQuery.data.value?.description">{{ projectQuery.data.value.description }}</p>
          </div>
        </div>

        <nav class="project-detail-page__tabs" aria-label="Разделы проекта">
          <button
            v-for="item in tabs"
            :key="item.id"
            type="button"
            :class="{ 'is-active': activeTab === item.id }"
            @click="activeTab = item.id"
          >
            {{ item.label }}
          </button>
        </nav>
      </header>

      <div v-if="projectQuery.isLoading.value" class="project-detail-page__state">Загрузка…</div>

      <template v-else>
        <div v-if="activeTab === 'docs'" class="project-detail-page__docs">
          <ProjectDocList
            :docs="docsQuery.data.value ?? []"
            :active-doc-id="activeDocId"
            @select="selectDoc"
            @create="createDoc"
          />
          <ProjectDocEditor
            :doc-id="activeDocId"
            :title="draftTitle"
            :body-md="draftBody"
            :saving="savingDoc"
            :disabled="!activeDocId"
            @update:title="draftTitle = $event"
            @update:body-md="draftBody = $event"
            @save="saveDoc"
            @delete="deleteDoc"
          />
        </div>

        <div v-else-if="activeTab === 'decisions'" class="project-detail-page__panel">
          <CreateDecisionForm :project-id="projectId" @created="refreshDecisions" />
          <ProjectDecisionList
            :project-id="projectId"
            :decisions="decisionsQuery.data.value ?? []"
            @updated="refreshDecisions"
          />
        </div>

        <div v-else class="project-detail-page__panel">
          <ProjectLinkedChats
            :project-id="projectId"
            :conversations="conversationsQuery.data.value ?? []"
            @link="linkModalOpen = true"
            @updated="refreshConversations"
          />
          <AddProjectMemberPanel
            :project-id="projectId"
            :members="membersQuery.data.value ?? []"
            @added="refreshMembers"
          />
        </div>
      </template>
    </section>

    <LinkProjectChatModal
      :open="linkModalOpen"
      :project-id="projectId"
      :linked="conversationsQuery.data.value ?? []"
      @close="linkModalOpen = false"
      @linked="refreshConversations"
    />
  </AppShell>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useQueryClient } from '@tanstack/vue-query'
import {
  useProject,
  useProjectConversations,
  useProjectDecisions,
  useProjectDoc,
  useProjectDocs,
  useProjectMembers,
} from '@/entities/project/api/project.queries'
import { projectApi } from '@/entities/project/api/project.api'
import type { ProjectTab } from '@/entities/project/model/types'
import AddProjectMemberPanel from '@/features/add-project-member/ui/AddProjectMemberPanel.vue'
import CreateDecisionForm from '@/features/create-decision/ui/CreateDecisionForm.vue'
import LinkProjectChatModal from '@/features/link-project-chat/ui/LinkProjectChatModal.vue'
import ProjectDocEditor from '@/widgets/project-doc-editor/ui/ProjectDocEditor.vue'
import ProjectDocList from '@/widgets/project-doc-list/ui/ProjectDocList.vue'
import ProjectDecisionList from '@/widgets/project-decision-list/ui/ProjectDecisionList.vue'
import ProjectLinkedChats from '@/widgets/project-linked-chats/ui/ProjectLinkedChats.vue'
import { queryKeys } from '@/shared/api/query-keys'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'
import AppIcon from '@/shared/ui/AppIcon.vue'
import AppShell from '@/shared/ui/AppShell.vue'

const route = useRoute()
const queryClient = useQueryClient()

const projectId = computed(() => String(route.params.projectId ?? ''))
const activeTab = ref<ProjectTab>('docs')
const activeDocId = ref<string | null>(null)
const draftTitle = ref('')
const draftBody = ref('')
const savingDoc = ref(false)
const linkModalOpen = ref(false)

const tabs = [
  { id: 'docs' as const, label: 'Документы' },
  { id: 'decisions' as const, label: 'Решения' },
  { id: 'chats' as const, label: 'Чаты и участники' },
]

const projectQuery = useProject(projectId)
const docsQuery = useProjectDocs(projectId)
const docQuery = useProjectDoc(projectId, activeDocId)
const decisionsQuery = useProjectDecisions(projectId)
const conversationsQuery = useProjectConversations(projectId)
const membersQuery = useProjectMembers(projectId)

watch(
  () => docsQuery.data.value,
  (docs) => {
    if (!docs || docs.length === 0) {
      activeDocId.value = null
      draftTitle.value = ''
      draftBody.value = ''
      return
    }
    if (!activeDocId.value || !docs.some((doc) => doc.id === activeDocId.value)) {
      activeDocId.value = docs[0].id
    }
  },
  { immediate: true },
)

watch(
  () => docQuery.data.value,
  (doc) => {
    if (!doc) return
    draftTitle.value = doc.title
    draftBody.value = doc.bodyMd
  },
)

function selectDoc(docId: string) {
  activeDocId.value = docId
}

async function createDoc() {
  try {
    const { doc } = await projectApi.createDoc(projectId.value, {
      title: 'Новый документ',
      bodyMd: '# Новый документ\n\n',
    })
    await queryClient.invalidateQueries({ queryKey: queryKeys.projectDocs(projectId.value) })
    activeDocId.value = doc.id
  } catch (err) {
    window.alert(getApiErrorMessage(err, { fallback: 'Не удалось создать документ' }))
  }
}

async function saveDoc() {
  if (!activeDocId.value) return
  savingDoc.value = true
  try {
    await projectApi.updateDoc(projectId.value, activeDocId.value, {
      title: draftTitle.value.trim(),
      bodyMd: draftBody.value,
    })
    await queryClient.invalidateQueries({ queryKey: queryKeys.projectDocs(projectId.value) })
    await queryClient.invalidateQueries({
      queryKey: queryKeys.projectDoc(projectId.value, activeDocId.value),
    })
  } catch (err) {
    window.alert(getApiErrorMessage(err, { fallback: 'Не удалось сохранить документ' }))
  } finally {
    savingDoc.value = false
  }
}

async function deleteDoc() {
  if (!activeDocId.value) return
  if (!window.confirm('Удалить документ?')) return
  try {
    await projectApi.deleteDoc(projectId.value, activeDocId.value)
    activeDocId.value = null
    await queryClient.invalidateQueries({ queryKey: queryKeys.projectDocs(projectId.value) })
  } catch (err) {
    window.alert(getApiErrorMessage(err, { fallback: 'Не удалось удалить документ' }))
  }
}

function refreshDecisions() {
  void queryClient.invalidateQueries({ queryKey: queryKeys.projectDecisions(projectId.value) })
}

function refreshConversations() {
  void queryClient.invalidateQueries({ queryKey: queryKeys.projectConversations(projectId.value) })
}

function refreshMembers() {
  void queryClient.invalidateQueries({ queryKey: queryKeys.projectMembers(projectId.value) })
}
</script>

<style scoped>
.project-detail-page {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
  overflow: hidden;
  background: var(--color-bg-elevated);
}
.project-detail-page__header {
  flex-shrink: 0;
  padding: 14px 16px 0;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-bg-elevated);
}
.project-detail-page__title {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}
.project-detail-page__back {
  display: inline-grid;
  place-items: center;
  width: 34px;
  height: 34px;
  margin-top: 2px;
  color: var(--color-text-secondary);
  border-radius: var(--radius-sm);
}
.project-detail-page__back:hover {
  background: var(--color-bg-hover);
}
.project-detail-page__title h1 {
  margin: 0;
  font-size: var(--font-size-2xl);
}
.project-detail-page__title p {
  margin: 4px 0 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.project-detail-page__tabs {
  display: flex;
  gap: 6px;
  margin-top: 14px;
  padding-bottom: 0;
}
.project-detail-page__tabs button {
  padding: 10px 12px;
  color: var(--color-text-secondary);
  background: transparent;
  border: 0;
  border-bottom: 2px solid transparent;
  cursor: pointer;
}
.project-detail-page__tabs button.is-active {
  color: var(--color-text-primary);
  border-bottom-color: var(--color-accent);
}
.project-detail-page__docs {
  display: flex;
  min-height: 0;
  flex: 1;
}
.project-detail-page__panel {
  overflow: auto;
  flex: 1;
  padding: 16px 20px 24px;
}
.project-detail-page__state {
  padding: 24px 20px;
  color: var(--color-text-secondary);
}
</style>
