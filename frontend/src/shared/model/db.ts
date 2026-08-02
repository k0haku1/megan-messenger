import Dexie, { type EntityTable } from 'dexie'
import type { Message } from '@/entities/message/model/types'

class MessengerDatabase extends Dexie {
  messages!: EntityTable<Message, 'id'>
  hiddenMessages!: EntityTable<{ id: string }, 'id'>

  constructor() {
    super('megan-messenger')
    this.version(1).stores({
      conversations: 'id, type, createdAt',
      messages: 'id, conversationId, [conversationId+createdAt], createdAt',
      syncStates: 'key, status, updatedAt',
    })
    this.version(2).stores({
      messages: 'id, conversationId, [conversationId+createdAt], createdAt',
    })
    this.version(3).stores({
      hiddenMessages: 'id',
    })
  }
}

export const db = new MessengerDatabase()
