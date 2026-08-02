import { db } from '@/shared/model/db'

export async function hideMessageForMe(messageId: string): Promise<void> {
  await db.hiddenMessages.put({ id: messageId })
}

export async function isMessageHidden(messageId: string): Promise<boolean> {
  const row = await db.hiddenMessages.get(messageId)
  return row != null
}
