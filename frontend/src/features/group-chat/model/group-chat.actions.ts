import { conversationApi } from '@/entities/conversation/api/conversation.api'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'
import { queryClient } from '@/shared/api/query-client'
import { queryKeys } from '@/shared/api/query-keys'

export async function openConversation(conversationId: string) {
  const navigation = useConversationSelectionStore()

  navigation.select(conversationId)
  navigation.clearPendingPeer()

  await queryClient.invalidateQueries({ queryKey: queryKeys.conversations })
}

export async function createGroupChat(title: string, memberIds: string[]) {
  const created = await conversationApi.createGroup({ title, memberIds })
  await openConversation(created.id)
  return created
}

export async function joinGroupChat(slug: string) {
  const joined = await conversationApi.joinBySlug(slug)
  await openConversation(joined.id)
  return joined
}

export async function leaveGroupChat(conversationId: string) {
  await conversationApi.leave(conversationId)

  const navigation = useConversationSelectionStore()
  navigation.closeConversation(conversationId)

  await queryClient.invalidateQueries({ queryKey: queryKeys.conversations })
}
