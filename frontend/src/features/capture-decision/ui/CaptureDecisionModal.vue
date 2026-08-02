<template>
  <BaseModalSheet
    :open="open"
    title="Зафиксировать решение"
    title-id="capture-decision-title"
    compact
    @close="emit('close')"
  >
    <form class="capture-decision-form" @submit.prevent="submit()">
      <blockquote v-if="message" class="capture-decision-form__quote">
        <strong>{{ message.sender.username }}</strong>
        <p>{{ message.content }}</p>
      </blockquote>

      <label v-if="projects.length > 1" class="capture-decision-form__field">
        <span>Проект</span>
        <select v-model="projectId" :disabled="pending">
          <option v-for="project in projects" :key="project.id" :value="project.id">
            {{ project.name }}
          </option>
        </select>
      </label>

      <BaseTextField
        v-model="summary"
        label="Решение"
        placeholder="Кратко: что решили"
        :disabled="pending"
      />

      <label class="capture-decision-form__field">
        <span>Контекст</span>
        <textarea
          v-model="context"
          rows="3"
          placeholder="Почему, альтернативы, ограничения"
          :disabled="pending"
        />
      </label>

      <p v-if="error" class="capture-decision-form__error">{{ error }}</p>
      <BaseButton type="submit" block :loading="pending">Зафиксировать</BaseButton>
    </form>
  </BaseModalSheet>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Message } from '@/entities/message/model/types'
import { projectApi } from '@/entities/project/api/project.api'
import { queryClient } from '@/shared/api/query-client'
import { queryKeys } from '@/shared/api/query-keys'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'

interface ProjectOption {
  id: string
  name: string
}

const props = defineProps<{
  open: boolean
  message: Message | null
  projects: ProjectOption[]
}>()

const emit = defineEmits<{ close: []; saved: [projectId: string] }>()

const projectId = ref('')
const summary = ref('')
const context = ref('')
const error = ref('')
const pending = ref(false)

watch(
  () => [props.open, props.message, props.projects] as const,
  ([isOpen, message, projects]) => {
    if (!isOpen || !message) return
    projectId.value = projects[0]?.id ?? ''
    summary.value = message.content.trim().slice(0, 300)
    context.value = ''
    error.value = ''
  },
  { immediate: true },
)

async function submit() {
  if (!props.message || !projectId.value) return

  error.value = ''
  if (summary.value.trim().length < 3) {
    error.value = 'Минимум 3 символа'
    return
  }

  pending.value = true
  try {
    await projectApi.createDecision(projectId.value, {
      summary: summary.value.trim(),
      context: context.value.trim(),
      messageId: props.message.id,
    })
    await queryClient.invalidateQueries({ queryKey: queryKeys.projectDecisions(projectId.value) })
    emit('saved', projectId.value)
    emit('close')
  } catch (err) {
    error.value = getApiErrorMessage(err, { fallback: 'Не удалось зафиксировать решение' })
  } finally {
    pending.value = false
  }
}
</script>

<style scoped>
.capture-decision-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.capture-decision-form__quote {
  margin: 0;
  padding: 10px 12px;
  background: var(--color-bg-subtle);
  border-left: 3px solid var(--color-accent);
  border-radius: var(--radius-sm);
}
.capture-decision-form__quote strong {
  display: block;
  margin-bottom: 4px;
  font-size: var(--font-size-sm);
}
.capture-decision-form__quote p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-normal);
  white-space: pre-wrap;
}
.capture-decision-form__field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.capture-decision-form__field span {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.capture-decision-form__field select,
.capture-decision-form__field textarea {
  padding: 10px 12px;
  font: inherit;
  color: inherit;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}
.capture-decision-form__field textarea {
  min-height: 72px;
  resize: vertical;
}
.capture-decision-form__error {
  margin: 0;
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}
</style>
