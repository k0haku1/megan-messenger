import { upsertMessage } from '@/entities/message/api/message.repository'
import { messageApi } from '@/entities/message/api/message.api'
import type { Message } from '@/entities/message/model/types'
import { queryClient } from '@/shared/api/query-client'
import { queryKeys } from '@/shared/api/query-keys'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'

export async function sendConversationMessage(conversationId: string, content: string): Promise<Message> {
  const { message } = await messageApi.send(conversationId, content)
  await upsertMessage(message)
  return message
}

export async function sendDirectMessage(payload: {
  userId?: string
  username?: string
  content: string
}): Promise<{ conversationId: string; message: Message }> {
  const { conversation, message } = await messageApi.sendDM(payload)
  await upsertMessage(message)
  await queryClient.invalidateQueries({ queryKey: queryKeys.conversations })
  return { conversationId: conversation.id, message }
}

export function getSendErrorMessage(error: unknown): string {
  return getApiErrorMessage(error, {
    fallback: 'Не удалось отправить сообщение',
    forbidden: 'Пользователь не принимает новые сообщения',
  })
}
