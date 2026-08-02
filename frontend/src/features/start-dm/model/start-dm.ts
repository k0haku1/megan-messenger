import { findDirectConversationByPeerId } from '@/entities/conversation/lib/find-dm'
import { useConversationSelectionStore } from '@/features/conversation-selection/model/conversation-selection.store'

export async function startDirectMessage(payload: {
  userId: string
  username: string
  avatarUrl?: string
}): Promise<void> {
  const navigation = useConversationSelectionStore()
  const existing = await findDirectConversationByPeerId(payload.userId)

  if (existing) {
    navigation.select(existing.id)
    return
  }

  navigation.openPendingPeer({
    id: payload.userId,
    username: payload.username,
    avatarUrl: payload.avatarUrl,
  })
}
