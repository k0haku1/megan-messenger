export function getConversationWsUrl(conversationId: string, token: string): string {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${window.location.host}/api/conversations/${conversationId}/ws?token=${encodeURIComponent(token)}`
}

export function getInboxWsUrl(token: string): string {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${window.location.host}/api/ws?token=${encodeURIComponent(token)}`
}
