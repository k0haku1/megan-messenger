export type MessageReceiptStatus = 'sent' | 'read'

export function ownMessageReceiptStatus(
  createdAt: string,
  othersReadAt: string | null | undefined,
): MessageReceiptStatus {
  if (!othersReadAt) return 'sent'
  const created = Date.parse(createdAt)
  const readAt = Date.parse(othersReadAt)
  if (Number.isNaN(created) || Number.isNaN(readAt)) return 'sent'
  return created <= readAt ? 'read' : 'sent'
}
