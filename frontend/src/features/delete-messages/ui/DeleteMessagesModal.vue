<template>
  <BaseModalSheet
    :open="open"
    :title="title"
    title-id="delete-messages-title"
    compact
    @close="emit('close')"
  >
    <p v-if="!canDeleteForEveryone" class="delete-messages-modal__text">
      {{ deleteMessagesCopy.onlyForYou }}
    </p>

    <label v-if="canDeleteForEveryone" class="delete-messages-modal__option">
      <input v-model="forEveryone" class="delete-messages-modal__native" type="checkbox" />
      <span class="delete-messages-modal__box" :class="{ 'is-checked': forEveryone }" aria-hidden="true">
        <AppIcon v-if="forEveryone" name="check" size="sm" />
      </span>
      <span>{{ deleteMessagesCopy.forEveryone }}</span>
    </label>

    <div class="delete-messages-modal__actions">
      <BaseButton variant="ghost" :disabled="pending" @click="emit('close')">
        {{ deleteMessagesCopy.cancel }}
      </BaseButton>
      <BaseButton variant="danger" :loading="pending" @click="confirm()">
        {{ deleteMessagesCopy.confirm }}
      </BaseButton>
    </div>
  </BaseModalSheet>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { deleteMessagesCopy, deleteMessagesTitle } from '../model/copy'
import AppIcon from '@/shared/ui/AppIcon.vue'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'

const props = withDefaults(
  defineProps<{
    open: boolean
    count: number
    canDeleteForEveryone: boolean
    pending?: boolean
  }>(),
  { pending: false },
)

const emit = defineEmits<{
  close: []
  confirm: [payload: { forEveryone: boolean }]
}>()

const forEveryone = ref(false)
const title = computed(() => deleteMessagesTitle(props.count))

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) forEveryone.value = false
  },
)

function confirm() {
  emit('confirm', {
    forEveryone: props.canDeleteForEveryone && forEveryone.value,
  })
}
</script>

<style scoped>
.delete-messages-modal__text {
  margin: 0 0 16px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-normal);
}
.delete-messages-modal__option {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0 0 18px;
  color: var(--color-text-primary);
  font-size: var(--font-size-sm);
  cursor: pointer;
  user-select: none;
}
.delete-messages-modal__native {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0 0 0 0);
  white-space: nowrap;
}
.delete-messages-modal__box {
  display: grid;
  flex: 0 0 auto;
  place-items: center;
  width: 20px;
  height: 20px;
  color: var(--color-text-on-accent);
  background: transparent;
  border: 1.5px solid var(--color-border);
  border-radius: 5px;
}
.delete-messages-modal__box.is-checked {
  background: var(--color-accent);
  border-color: var(--color-accent);
}
.delete-messages-modal__box :deep(.app-icon) {
  width: 12px;
  height: 12px;
}
.delete-messages-modal__actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
