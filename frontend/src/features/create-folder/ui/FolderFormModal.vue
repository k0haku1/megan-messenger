<template>
  <BaseModalSheet
    :open="open"
    :title="mode === 'create' ? 'Новая папка' : 'Переименовать папку'"
    :title-id="mode === 'create' ? 'create-folder-title' : 'rename-folder-title'"
    compact
    @close="emit('close')"
  >
    <form class="folder-form" @submit.prevent="submit()">
      <BaseTextField
        v-model="name"
        label="Название папки"
        :placeholder="mode === 'create' ? 'Работа, Друзья' : undefined"
        autocomplete="off"
        :disabled="isPending"
        autofocus
      />
      <p v-if="error" class="folder-form__error">{{ error }}</p>
      <BaseButton type="submit" block :loading="isPending">
        {{ mode === 'create' ? 'Создать' : 'Сохранить' }}
      </BaseButton>
    </form>
  </BaseModalSheet>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'
import BaseButton from '@/shared/ui/BaseButton.vue'
import { useCreateFolder, useUpdateFolder } from '@/entities/folder/api/folder.queries'

const props = defineProps<{
  open: boolean
  mode: 'create' | 'rename'
  folderId?: string
  initialName?: string
}>()

const emit = defineEmits<{ close: [] }>()

const name = ref(props.initialName ?? '')
const error = ref<string | null>(null)

const { mutateAsync: createFolder, isPending: isCreating } = useCreateFolder()
const { mutateAsync: updateFolder, isPending: isUpdating } = useUpdateFolder()

const isPending = computed(() => isCreating.value || isUpdating.value)

watch(
  () => [props.open, props.initialName] as const,
  ([open, initialName]) => {
    if (!open) return
    name.value = initialName ?? ''
    error.value = null
  },
)

async function submit() {
  const trimmed = name.value.trim()
  if (!trimmed) return

  error.value = null
  try {
    if (props.mode === 'create') {
      await createFolder({ name: trimmed })
    } else if (props.folderId) {
      await updateFolder({ folderId: props.folderId, name: trimmed })
    }
    emit('close')
  } catch {
    error.value =
      props.mode === 'create' ? 'Не удалось создать папку' : 'Не удалось переименовать папку'
  }
}
</script>

<style scoped>
.folder-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.folder-form__error {
  color: var(--color-danger);
  font-size: var(--font-size-sm);
  margin: 0;
}
</style>
