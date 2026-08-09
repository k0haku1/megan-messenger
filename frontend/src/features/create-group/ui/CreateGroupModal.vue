<template>
  <BaseModalSheet
    :open="open"
    :title="step === 'form' ? 'Новая группа' : 'Группа создана'"
    title-id="create-group-title"
    compact
    scrollable
    @close="close()"
  >
    <form v-if="step === 'form'" class="create-group-form" @submit.prevent="submit()">
      <BaseTextField
        v-model="title"
        label="Название группы"
        placeholder="Команда Megan"
        autocomplete="off"
        :disabled="pending"
      />

      <div class="create-group-form__members">
        <label class="create-group-form__search">
          <span>Участники</span>
          <input
            v-model="memberQuery"
            type="search"
            placeholder="Поиск @username"
            autocomplete="off"
            :disabled="pending"
          />
        </label>

        <div v-if="selectedMembers.length > 0" class="create-group-form__chips">
          <button
            v-for="member in selectedMembers"
            :key="member.id"
            class="create-group-form__chip"
            type="button"
            :disabled="pending"
            @click="removeMember(member.id)"
          >
            @{{ member.username }}
            <AppIcon name="close" size="sm" />
          </button>
        </div>

        <div v-if="memberQuery.trim().length >= 2" class="create-group-form__results">
          <p v-if="isSearching" class="create-group-form__hint">Ищем…</p>
          <button
            v-for="(user, index) in searchResults"
            :key="user.id"
            class="create-group-form__result"
            type="button"
            :disabled="pending || isSelected(user.id)"
            @click="addMember(user)"
          >
            <BaseAvatar :name="user.username" :src="user.avatarUrl" :color="index" size="sm" />
            <span>@{{ user.username }}</span>
          </button>
          <p
            v-if="!isSearching && searchResults.length === 0"
            class="create-group-form__hint"
          >
            Пользователи не найдены
          </p>
        </div>
      </div>

      <p v-if="error" class="create-group-form__error">{{ error }}</p>
      <BaseButton type="submit" block :loading="pending">Создать группу</BaseButton>
    </form>

    <div v-else class="create-group-success">
      <p>Группа «{{ createdTitle }}» готова. Отправьте ссылку-приглашение участникам:</p>
      <div class="create-group-success__slug">
        <code>{{ createdSlug }}</code>
        <BaseButton variant="ghost" @click="copySlug()">
          {{ copied ? 'Скопировано' : 'Копировать' }}
        </BaseButton>
      </div>
      <BaseButton block @click="finish()">Открыть чат</BaseButton>
    </div>
  </BaseModalSheet>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import { userApi } from '@/entities/user/api/user.api'
import type { UserSearchResult } from '@/entities/user/model/types'
import { normalizeUsernameQuery } from '@/entities/user/lib/username'
import { createGroupChat } from '@/features/group-chat/model/group-chat.actions'
import { useSessionStore } from '@/entities/session/model/session.store'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'
import { copyGroupInviteSlug } from '@/shared/lib/group-invite'
import AppIcon from '@/shared/ui/AppIcon.vue'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: [] }>()

const session = useSessionStore()
const step = ref<'form' | 'success'>('form')
const title = ref('')
const memberQuery = ref('')
const searchResults = ref<UserSearchResult[]>([])
const selectedMembers = ref<UserSearchResult[]>([])
const isSearching = ref(false)
const pending = ref(false)
const error = ref('')
const createdSlug = ref('')
const createdTitle = ref('')
const copied = ref(false)

const runSearch = useDebounceFn(async (rawQuery: string) => {
  const normalized = normalizeUsernameQuery(rawQuery)
  if (normalized.length < 2) {
    searchResults.value = []
    isSearching.value = false
    return
  }

  isSearching.value = true
  try {
    const response = await userApi.search(normalized)
    searchResults.value = response.users.filter((user) => user.id !== session.user?.id)
  } catch {
    searchResults.value = []
  } finally {
    isSearching.value = false
  }
}, 300)

watch(memberQuery, (value) => {
  void runSearch(value)
})

watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen) reset()
  },
)

function isSelected(userId: string) {
  return selectedMembers.value.some((member) => member.id === userId)
}

function addMember(user: UserSearchResult) {
  if (isSelected(user.id)) return
  selectedMembers.value = [...selectedMembers.value, user]
  memberQuery.value = ''
  searchResults.value = []
}

function removeMember(userId: string) {
  selectedMembers.value = selectedMembers.value.filter((member) => member.id !== userId)
}

async function submit() {
  error.value = ''
  const trimmedTitle = title.value.trim()
  if (trimmedTitle.length < 3) {
    error.value = 'Название — минимум 3 символа'
    return
  }

  pending.value = true
  try {
    const created = await createGroupChat(
      trimmedTitle,
      selectedMembers.value.map((member) => member.id),
    )
    createdSlug.value = created.slug
    createdTitle.value = trimmedTitle
    step.value = 'success'
  } catch (err) {
    error.value = getApiErrorMessage(err, { fallback: 'Не удалось создать группу' })
  } finally {
    pending.value = false
  }
}

async function copySlug() {
  if (!createdSlug.value) return
  await copyGroupInviteSlug(createdSlug.value)
  copied.value = true
}

function finish() {
  emit('close')
}

function close() {
  emit('close')
}

function reset() {
  step.value = 'form'
  title.value = ''
  memberQuery.value = ''
  searchResults.value = []
  selectedMembers.value = []
  error.value = ''
  createdSlug.value = ''
  createdTitle.value = ''
  copied.value = false
  pending.value = false
  isSearching.value = false
}
</script>

<style scoped>
.create-group-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.create-group-form__members {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.create-group-form__search {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.create-group-form__search span {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.create-group-form__search input {
  height: 44px;
  padding: 0 14px;
  font: inherit;
  color: inherit;
  background: var(--color-bg-subtle);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}
.create-group-form__chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.create-group-form__chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px 6px 10px;
  color: var(--color-text-primary);
  background: var(--color-bg-subtle);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-round);
  cursor: pointer;
}
.create-group-form__results {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.create-group-form__result {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  text-align: left;
  background: transparent;
  border: 0;
  border-radius: var(--radius-sm);
  cursor: pointer;
}
.create-group-form__result:hover:not(:disabled) {
  background: var(--color-bg-hover);
}
.create-group-form__result:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}
.create-group-form__hint {
  margin: 0;
  padding: 4px 2px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.create-group-form__error {
  margin: 0;
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}
.create-group-success p {
  margin: 0 0 12px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-normal);
}
.create-group-success__slug {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;
  padding: 10px 12px;
  background: var(--color-bg-subtle);
  border-radius: var(--radius-sm);
}
.create-group-success__slug code {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: var(--font-size-md);
}
</style>
