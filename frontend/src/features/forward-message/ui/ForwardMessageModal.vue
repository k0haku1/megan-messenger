<template>
  <BaseModalSheet
    :open="open"
    title="Переслать сообщение"
    title-id="forward-message-title"
    compact
    scrollable
    @close="emit('close')"
  >
    <blockquote v-if="message" class="forward-message-modal__quote">
      <strong>{{ message.sender.username }}</strong>
      <p>{{ message.content }}</p>
    </blockquote>

    <p class="forward-message-modal__hint">Выберите чат</p>

    <div v-if="targets.length === 0" class="forward-message-modal__empty">
      Нет других чатов для пересылки
    </div>

    <button
      v-for="conversation in targets"
      :key="conversation.id"
      class="forward-message-modal__item"
      type="button"
      :disabled="pendingId === conversation.id"
      @click="forwardTo(conversation.id)"
    >
      <strong>{{ getConversationTitle(conversation) }}</strong>
      <span>{{ conversation.type === 'group' ? 'групповой чат' : 'личный чат' }}</span>
    </button>

    <p v-if="error" class="forward-message-modal__error">{{ error }}</p>
  </BaseModalSheet>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useConversations } from '@/entities/conversation/api/conversation.queries'
import { getConversationTitle } from '@/entities/conversation/lib/display'
import type { Message } from '@/entities/message/model/types'
import { messageApi } from '@/entities/message/api/message.api'
import { upsertMessage } from '@/entities/message/api/message.repository'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'

const props = defineProps<{
  open: boolean
  message: Message | null
  currentConversationId: string | null
}>()

const emit = defineEmits<{ close: []; forwarded: [conversationId: string] }>()

const { conversations } = useConversations()
const error = ref('')
const pendingId = ref<string | null>(null)

const targets = computed(() =>
  conversations.value.filter((conversation) => conversation.id !== props.currentConversationId),
)

watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen) {
      error.value = ''
      pendingId.value = null
    }
  },
)

async function forwardTo(conversationId: string) {
  if (!props.message) return

  error.value = ''
  pendingId.value = conversationId
  try {
    if (!props.currentConversationId) return
    const { message } = await messageApi.forward(props.currentConversationId, props.message.id, conversationId)
    await upsertMessage(message)
    emit('forwarded', conversationId)
    emit('close')
  } catch (err) {
    error.value = getApiErrorMessage(err, { fallback: 'Не удалось переслать сообщение' })
  } finally {
    pendingId.value = null
  }
}
</script>

<style scoped>
.forward-message-modal__quote {
  margin: 0 0 12px;
  padding: 10px 12px;
  background: var(--color-bg-subtle);
  border-left: 3px solid var(--color-border);
  border-radius: var(--radius-sm);
}
.forward-message-modal__quote strong {
  display: block;
  margin-bottom: 4px;
  font-size: var(--font-size-sm);
}
.forward-message-modal__quote p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  white-space: pre-wrap;
}
.forward-message-modal__hint {
  margin: 0 0 10px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.forward-message-modal__empty {
  padding: 8px 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.forward-message-modal__item {
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
.forward-message-modal__item:disabled {
  cursor: not-allowed;
  opacity: 0.65;
}
.forward-message-modal__item span {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.forward-message-modal__error {
  margin: 10px 0 0;
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}
</style>
