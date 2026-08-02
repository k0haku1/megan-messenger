<template>
  <section class="project-member-panel">
    <h3>Участники</h3>
    <ul class="project-member-panel__list">
      <li v-for="member in members" :key="member.userId">
        <strong>@{{ member.username }}</strong>
        <span>{{ member.role === 'owner' ? 'владелец' : 'участник' }}</span>
      </li>
    </ul>

    <form class="project-member-panel__form" @submit.prevent="submit()">
      <BaseTextField
        v-model="username"
        label="Добавить по @username"
        placeholder="username"
        autocomplete="off"
        :disabled="pending"
      />
      <p v-if="error" class="project-member-panel__error">{{ error }}</p>
      <BaseButton type="submit" block :loading="pending">Добавить</BaseButton>
    </form>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { projectApi } from '@/entities/project/api/project.api'
import type { ProjectMember } from '@/entities/project/model/types'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'

const props = defineProps<{
  projectId: string
  members: ProjectMember[]
}>()

const emit = defineEmits<{ added: [] }>()

const username = ref('')
const error = ref('')
const pending = ref(false)

async function submit() {
  error.value = ''
  const value = username.value.trim().replace(/^@+/, '')
  if (value.length < 5) {
    error.value = 'Укажите username'
    return
  }

  pending.value = true
  try {
    await projectApi.addMember(props.projectId, value)
    username.value = ''
    emit('added')
  } catch (err) {
    error.value = getApiErrorMessage(err, { fallback: 'Не удалось добавить участника' })
  } finally {
    pending.value = false
  }
}
</script>

<style scoped>
.project-member-panel h3 {
  margin: 0 0 10px;
  font-size: var(--font-size-lg);
}
.project-member-panel__list {
  margin: 0 0 14px;
  padding: 0;
  list-style: none;
}
.project-member-panel__list li {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 8px 0;
  border-bottom: 1px solid var(--color-border);
}
.project-member-panel__list span {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.project-member-panel__form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.project-member-panel__error {
  margin: 0;
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}
</style>
