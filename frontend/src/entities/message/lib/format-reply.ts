import type { Message } from '@/entities/message/model/types'

export function formatReplyContent(replyTo: Message, content: string): string {
  const quoted = replyTo.content
    .split('\n')
    .map((line) => `> ${line}`)
    .join('\n')

  return `↩ ${replyTo.sender.username}:\n${quoted}\n\n${content}`
}
