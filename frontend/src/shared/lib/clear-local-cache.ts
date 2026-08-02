import { db } from '@/shared/model/db'

export async function clearLocalCache(): Promise<void> {
  await db.messages.clear()
}
