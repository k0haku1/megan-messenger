<template>
  <BaseModalSheet
    :open="open"
    title="Вступить в группу"
    title-id="join-group-title"
    compact
    @close="emit('close')"
  >
    <form class="join-group-form" @submit.prevent="submit()">
      <BaseTextField
        v-model="inviteInput"
        label="Ссылка или код приглашения"
        placeholder="abc123def45"
        autocomplete="off"
        :disabled="pending"
      />
      <p class="join-group-form__hint">11 символов из ссылки-приглашения группы</p>
      <p v-if="error" class="join-group-form__error">{{ error }}</p>
      <BaseButton type="submit" block :loading="pending">Вступить</BaseButton>
    </form>
  </BaseModalSheet>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { joinGroupChat } from '@/features/group-chat/model/group-chat.actions'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'
import { parseGroupInviteSlug } from '@/shared/lib/group-invite'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'
import BaseTextField from '@/shared/ui/BaseTextField.vue'

defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: [] }>()

const inviteInput = ref('')
const pending = ref(false)
const error = ref('')

watch(inviteInput, () => {
  error.value = ''
})

async function submit() {
  const slug = parseGroupInviteSlug(inviteInput.value)
  if (!slug) {
    error.value = 'Некорректный код приглашения'
    return
  }

  pending.value = true
  try {
    await joinGroupChat(slug)
    inviteInput.value = ''
    emit('close')
  } catch (err) {
    error.value = getApiErrorMessage(err, { fallback: 'Не удалось вступить в группу' })
  } finally {
    pending.value = false
  }
}
</script>

<style scoped>
.join-group-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.join-group-form__hint {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.join-group-form__error {
  margin: 0;
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}
</style>
