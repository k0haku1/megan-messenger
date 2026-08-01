import { request } from '@/shared/api/http'
import type { Conversation } from '@/entities/conversation/model/types'

export const conversationApi = {
  list: () => request<{ conversations: Conversation[] }>('/conversations'),
}
