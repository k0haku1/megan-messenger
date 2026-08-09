import { db } from '@/shared/model/db'
import type { Message } from '@/entities/message/model/types'

export async function upsertMessage(message: Message): Promise<void> {
  await db.messages.put(message)
}

export async function removeMessage(messageId: string): Promise<void> {
  await db.messages.delete(messageId)
}

export async function replaceMessagesForConversation(
  conversationId: string,
  messages: Message[],
): Promise<void> {
  await db.transaction('rw', db.messages, async () => {
    await db.messages.where('conversationId').equals(conversationId).delete()
    if (messages.length > 0) {
      await db.messages.bulkPut(messages)
    }
  })
}
