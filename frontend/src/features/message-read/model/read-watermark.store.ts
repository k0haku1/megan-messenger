import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useReadWatermarkStore = defineStore('read-watermark', () => {
  const othersReadAtByConversation = ref<Record<string, string | null>>({})

  function getOthersReadAt(conversationId: string): string | null {
    return othersReadAtByConversation.value[conversationId] ?? null
  }

  function setOthersReadAt(conversationId: string, othersReadAt: string | null | undefined): void {
    othersReadAtByConversation.value = {
      ...othersReadAtByConversation.value,
      [conversationId]: othersReadAt ?? null,
    }
  }

  function clear(conversationId: string): void {
    const next = { ...othersReadAtByConversation.value }
    delete next[conversationId]
    othersReadAtByConversation.value = next
  }

  return {
    othersReadAtByConversation,
    getOthersReadAt,
    setOthersReadAt,
    clear,
  }
})
