import { queryClient } from '@/shared/api/query-client'
import { queryKeys } from '@/shared/api/query-keys'
import type { Conversation } from '@/entities/conversation/model/types'

export function findDirectConversationByPeerId(peerUserId: string): Conversation | undefined {
  const data = queryClient.getQueryData<{ conversations: Conversation[] }>(queryKeys.conversations)
  return data?.conversations.find(
    (conversation) => conversation.type === 'dm' && conversation.peer?.id === peerUserId,
  )
}
