<template>
  <BaseModalSheet
    :open="open"
    title="Добавить в папку"
    title-id="assign-folder-title"
    compact
    scrollable
    @close="emit('close')"
  >
    <div class="assign-folder">
      <p v-if="availableFolders.length === 0" class="assign-folder__empty">
        {{ emptyMessage }}
      </p>

      <button
        v-for="folder in availableFolders"
        :key="folder.id"
        type="button"
        class="assign-folder__item"
        :class="{ 'is-selected': selectedId === folder.id }"
        @click="selectedId = folder.id"
      >
        <FolderBadge :name="folder.name" size="md" />
        <span class="assign-folder__name">{{ folder.name }}</span>
        <AppIcon v-if="selectedId === folder.id" name="check" size="sm" class="assign-folder__check" />
      </button>

      <p v-if="error" class="assign-folder__error">{{ error }}</p>

      <div v-if="availableFolders.length > 0" class="assign-folder__footer">
        <BaseButton block :disabled="!selectedId" :loading="isPending" @click="submit()">
          Добавить
        </BaseButton>
      </div>
    </div>
  </BaseModalSheet>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'
import BaseButton from '@/shared/ui/BaseButton.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import FolderBadge from '@/entities/folder/ui/FolderBadge.vue'
import { useFolders, useAssignChatFolder } from '@/entities/folder/api/folder.queries'

const props = defineProps<{
  open: boolean
  conversationId: string
  assignedFolderIds?: string[]
}>()

const emit = defineEmits<{ close: [] }>()

const { folders } = useFolders()
const { mutateAsync: addItem, isPending } = useAssignChatFolder()

const selectedId = ref<string | null>(null)
const error = ref<string | null>(null)

const availableFolders = computed(() => {
  const assigned = new Set(props.assignedFolderIds ?? [])
  return folders.value.filter((folder) => !assigned.has(folder.id))
})

const emptyMessage = computed(() =>
  folders.value.length === 0
    ? 'Нет папок. Создайте папку в боковой панели.'
    : 'Чат уже добавлен во все папки.',
)

watch(
  () => props.open,
  (open) => {
    if (!open) return
    selectedId.value = null
    error.value = null
  },
)

async function submit() {
  if (!selectedId.value) return

  error.value = null
  try {
    await addItem({ folderId: selectedId.value, conversationIds: [props.conversationId] })
    emit('close')
  } catch {
    error.value = 'Не удалось добавить чат в папку'
  }
}
</script>

<style scoped>
.assign-folder {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.assign-folder__empty,
.assign-folder__error {
  margin: 0;
  font-size: var(--font-size-sm);
}

.assign-folder__empty {
  color: var(--color-text-secondary);
  padding: 8px 0;
}

.assign-folder__error {
  color: var(--color-danger);
  padding-top: 4px;
}

.assign-folder__item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 10px 12px;
  background: transparent;
  border: 0;
  border-radius: var(--radius-sm);
  cursor: pointer;
  text-align: left;
  color: var(--color-text-primary);
  font: inherit;
  font-size: var(--font-size-sm);
}

.assign-folder__item:hover {
  background: var(--color-bg-subtle);
}

.assign-folder__item.is-selected {
  background: var(--color-bg-subtle);
  color: var(--color-primary);
}

.assign-folder__name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.assign-folder__check {
  color: var(--color-primary);
}

.assign-folder__footer {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid var(--color-border);
}
</style>
