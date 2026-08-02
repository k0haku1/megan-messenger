<template>
  <BaseModalSheet
    :open="open"
    title="О группе"
    title-id="group-info-title"
    compact
    @close="emit('close')"
  >
    <section class="group-info">
      <div class="group-info__identity">
        <BaseAvatar :name="title" size="lg" :color="2" />
        <div>
          <strong>{{ title }}</strong>
          <span>Групповой чат</span>
        </div>
      </div>

      <div v-if="slug" class="group-info__invite">
        <span class="group-info__label">Код приглашения</span>
        <div class="group-info__slug">
          <code>{{ slug }}</code>
          <BaseButton variant="ghost" @click="copySlug()">
            {{ copied ? 'Скопировано' : 'Копировать' }}
          </BaseButton>
        </div>
        <p>Отправьте код коллегам — они смогут вступить через «Новый чат → Вступить по ссылке».</p>
      </div>

      <div class="group-info__leave">
        <BaseButton variant="danger" block :loading="leaving" @click="leave()">
          Покинуть группу
        </BaseButton>
        <p v-if="leaveError" class="group-info__error">{{ leaveError }}</p>
      </div>
    </section>
  </BaseModalSheet>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { leaveGroupChat } from '@/features/group-chat/model/group-chat.actions'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'
import { copyGroupInviteSlug } from '@/shared/lib/group-invite'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'

const props = defineProps<{
  open: boolean
  conversationId: string
  title: string
  slug?: string
}>()

const emit = defineEmits<{ close: [] }>()

const copied = ref(false)
const leaving = ref(false)
const leaveError = ref('')

watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen) {
      copied.value = false
      leaveError.value = ''
    }
  },
)

async function copySlug() {
  if (!props.slug) return
  await copyGroupInviteSlug(props.slug)
  copied.value = true
}

async function leave() {
  if (!window.confirm('Покинуть группу? Вы сможете вернуться по коду приглашения.')) return

  leaving.value = true
  leaveError.value = ''
  try {
    await leaveGroupChat(props.conversationId)
    emit('close')
  } catch (err) {
    leaveError.value = getApiErrorMessage(err, { fallback: 'Не удалось покинуть группу' })
  } finally {
    leaving.value = false
  }
}
</script>

<style scoped>
.group-info__identity {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  padding: 12px;
  background: var(--color-bg-subtle);
  border-radius: var(--radius-md);
}
.group-info__identity span {
  display: block;
  margin-top: 2px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.group-info__label {
  display: block;
  margin-bottom: 6px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}
.group-info__slug {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  padding: 10px 12px;
  background: var(--color-bg-subtle);
  border-radius: var(--radius-sm);
}
.group-info__slug code {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
}
.group-info__invite p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-normal);
}
.group-info__leave {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--color-border-subtle);
}
.group-info__error {
  margin: 8px 0 0;
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}
</style>
