<template>
  <span
    class="message-status-ticks"
    :class="{ 'is-read': status === 'read' }"
    :aria-label="status === 'read' ? 'Прочитано' : 'Отправлено'"
  >
    <AppIcon :name="status === 'read' ? 'checks' : 'check'" size="sm" />
  </span>
</template>

<script setup lang="ts">
import type { MessageReceiptStatus } from '@/entities/message/lib/receipt-status'
import AppIcon from '@/shared/ui/AppIcon.vue'

defineProps<{
  status: MessageReceiptStatus
}>()
</script>

<style scoped>
.message-status-ticks {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  color: var(--color-text-tertiary);
}
.message-status-ticks.is-read {
  color: var(--color-accent);
}
.message-status-ticks :deep(.app-icon) {
  width: 14px;
  height: 14px;
}
</style>
