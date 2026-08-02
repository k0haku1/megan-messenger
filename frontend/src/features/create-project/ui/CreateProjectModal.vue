<template>
  <BaseModalSheet
    :open="open"
    title="Новый проект"
    title-id="create-project-title"
    compact
    @close="emit('close')"
  >
    <form class="project-form" @submit.prevent="submit()">
      <BaseTextField v-model="name" label="Название" autocomplete="off" :disabled="pending" />
      <BaseTextField
        v-model="description"
        label="Описание"
        autocomplete="off"
        :disabled="pending"
      />
      <p v-if="error" class="project-form__error">{{ error }}</p>
      <BaseButton type="submit" block :loading="pending">Создать</BaseButton>
    </form>
  </BaseModalSheet>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCreateProject } from '@/entities/project/api/project.queries'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'

defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: [] }>()

const router = useRouter()
const createProject = useCreateProject()

const name = ref('')
const description = ref('')
const error = ref('')
const pending = ref(false)

async function submit() {
  error.value = ''
  if (name.value.trim().length < 2) {
    error.value = 'Минимум 2 символа'
    return
  }

  pending.value = true
  try {
    const project = await createProject.mutateAsync({
      name: name.value.trim(),
      description: description.value.trim(),
    })
    emit('close')
    name.value = ''
    description.value = ''
    await router.push({ name: 'project-detail', params: { projectId: project.id } })
  } catch (err) {
    error.value = getApiErrorMessage(err, { fallback: 'Не удалось создать проект' })
  } finally {
    pending.value = false
  }
}
</script>

<style scoped>
.project-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.project-form__error {
  margin: 0;
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}
</style>
