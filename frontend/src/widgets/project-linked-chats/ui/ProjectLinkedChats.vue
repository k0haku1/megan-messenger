<template>
  <section class="project-linked-chats">
    <header class="project-linked-chats__header">
      <h3>Привязанные чаты</h3>
      <BaseButton @click="emit('link')">Привязать</BaseButton>
    </header>

    <article v-for="chat in conversations" :key="chat.conversationId" class="project-linked-chats__item">
      <div>
        <strong>{{ chat.title || 'Групповой чат' }}</strong>
        <span>группа</span>
      </div>
      <BaseButton
        variant="ghost"
        :loading="pendingId === chat.conversationId"
        @click="unlink(chat.conversationId)"
      >
        Отвязать
      </BaseButton>
    </article>

    <div v-if="conversations.length === 0" class="project-linked-chats__empty">
      Нет привязанных чатов
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { projectApi } from '@/entities/project/api/project.api'
import type { ProjectConversationLink } from '@/entities/project/model/types'
import BaseButton from '@/shared/ui/BaseButton.vue'

const props = defineProps<{ projectId: string; conversations: ProjectConversationLink[] }>()
const emit = defineEmits<{ link: []; updated: [] }>()

const pendingId = ref<string | null>(null)

async function unlink(conversationId: string) {
  pendingId.value = conversationId
  try {
    await projectApi.unlinkConversation(props.projectId, conversationId)
    emit('updated')
  } finally {
    pendingId.value = null
  }
}
</script>

<style scoped>
.project-linked-chats__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 12px;
}
.project-linked-chats__header h3 {
  margin: 0;
  font-size: var(--font-size-lg);
}
.project-linked-chats__item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
  padding: 12px 14px;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}
.project-linked-chats__item span {
  display: block;
  margin-top: 2px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.project-linked-chats__empty {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
</style>
