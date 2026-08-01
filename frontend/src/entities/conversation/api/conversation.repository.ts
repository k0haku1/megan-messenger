import { db } from '@/shared/model/db'
import type { Conversation } from '@/entities/conversation/model/types'

export async function replaceConversations(conversations: Conversation[]): Promise<void> {
  await db.transaction('rw', db.conversations, async () => {
    await db.conversations.clear()
    await db.conversations.bulkPut(conversations)
  })
}
