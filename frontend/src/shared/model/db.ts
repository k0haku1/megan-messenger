import Dexie, { type EntityTable } from 'dexie'
import type { Conversation, Message } from '@/entities/conversation/model/types'

export interface SyncState { key: string; status: 'idle' | 'loading' | 'success' | 'error'; updatedAt: number; error?: string }

class MessengerDatabase extends Dexie {
  conversations!: EntityTable<Conversation, 'id'>
  messages!: EntityTable<Message, 'id'>
  syncStates!: EntityTable<SyncState, 'key'>

  constructor() {
    super('megan-messenger')
    this.version(1).stores({
      conversations: 'id, type, createdAt',
      messages: 'id, conversationId, [conversationId+createdAt], createdAt',
      syncStates: 'key, status, updatedAt',
    })
  }
}

export const db = new MessengerDatabase()
