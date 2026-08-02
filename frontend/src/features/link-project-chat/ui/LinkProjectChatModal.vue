<template>
  <BaseModalSheet
    :open="open"
    title="Привязать чат"
    title-id="link-project-chat-title"
    compact
    scrollable
    @close="emit('close')"
  >
    <p class="link-chat-modal__hint">Можно привязать только групповые чаты, где вы участник.</p>

    <div v-if="groupConversations.length === 0" class="link-chat-modal__empty">
      Нет доступных групповых чатов
    </div>

    <button
      v-for="conversation in groupConversations"
      :key="conversation.id"
      class="link-chat-modal__item"
      type="button"
      :disabled="pendingId === conversation.id || isLinked(conversation.id)"
      @click="link(conversation.id)"
    >
      <strong>{{ getConversationTitle(conversation) }}</strong>
      <span v-if="isLinked(conversation.id)">уже привязан</span>
      <span v-else>групповой чат</span>
    </button>

    <p v-if="error" class="link-chat-modal__error">{{ error }}</p>
  </BaseModalSheet>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useConversations } from '@/entities/conversation/api/conversation.queries'
import { getConversationTitle } from '@/entities/conversation/lib/display'
import { projectApi } from '@/entities/project/api/project.api'
import type { ProjectConversationLink } from '@/entities/project/model/types'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'

const props = defineProps<{
  open: boolean
  projectId: string
  linked: ProjectConversationLink[]
}>()

const emit = defineEmits<{ close: [] ; linked: [] }>()

const { conversations } = useConversations()
const error = ref('')
const pendingId = ref<string | null>(null)

const groupConversations = computed(() =>
  conversations.value.filter((conversation) => conversation.type === 'group'),
)

function isLinked(conversationId: string) {
  return props.linked.some((item) => item.conversationId === conversationId)
}

async function link(conversationId: string) {
  error.value = ''
  pendingId.value = conversationId
  try {
    await projectApi.linkConversation(props.projectId, conversationId)
    emit('linked')
    emit('close')
  } catch (err) {
    error.value = getApiErrorMessage(err, { fallback: 'Не удалось привязать чат' })
  } finally {
    pendingId.value = null
  }
}
</script>

<style scoped>
.link-chat-modal__hint {
  margin: 0 0 12px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.link-chat-modal__empty {
  padding: 12px 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.link-chat-modal__item {
  display: flex;
  width: 100%;
  flex-direction: column;
  gap: 2px;
  margin-bottom: 6px;
  padding: 10px 12px;
  text-align: left;
  background: var(--color-bg-subtle);
  border: 0;
  border-radius: var(--radius-sm);
  cursor: pointer;
}
.link-chat-modal__item:disabled {
  cursor: not-allowed;
  opacity: 0.65;
}
.link-chat-modal__item span {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.link-chat-modal__error {
  margin: 10px 0 0;
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}
</style>
