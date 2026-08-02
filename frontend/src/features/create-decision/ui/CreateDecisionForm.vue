<template>
  <form class="decision-form" @submit.prevent="submit()">
    <BaseTextField
      v-model="summary"
      label="Решение"
      placeholder="Кратко: что решили"
      :disabled="pending"
    />
    <label class="decision-form__field">
      <span>Контекст</span>
      <textarea
        v-model="context"
        rows="4"
        placeholder="Почему, альтернативы, ограничения"
        :disabled="pending"
      />
    </label>
    <p v-if="error" class="decision-form__error">{{ error }}</p>
    <BaseButton type="submit" block :loading="pending">Зафиксировать</BaseButton>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { projectApi } from '@/entities/project/api/project.api'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'

const props = defineProps<{ projectId: string }>()
const emit = defineEmits<{ created: [] }>()

const summary = ref('')
const context = ref('')
const error = ref('')
const pending = ref(false)

async function submit() {
  error.value = ''
  if (summary.value.trim().length < 3) {
    error.value = 'Минимум 3 символа'
    return
  }

  pending.value = true
  try {
    await projectApi.createDecision(props.projectId, {
      summary: summary.value.trim(),
      context: context.value.trim(),
    })
    summary.value = ''
    context.value = ''
    emit('created')
  } catch (err) {
    error.value = getApiErrorMessage(err, { fallback: 'Не удалось сохранить решение' })
  } finally {
    pending.value = false
  }
}
</script>

<style scoped>
.decision-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 18px;
  padding: 14px;
  background: var(--color-bg-subtle);
  border-radius: var(--radius-md);
}
.decision-form__field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.decision-form__field span {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.decision-form__field textarea {
  min-height: 96px;
  padding: 10px 12px;
  resize: vertical;
  font: inherit;
  color: inherit;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}
.decision-form__error {
  margin: 0;
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}
</style>
