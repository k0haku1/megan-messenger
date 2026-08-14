export const queryKeys = {
  folders: ['folders'] as const,
  conversations: ['conversations'] as const,
  projects: ['projects'] as const,
  project: (id: string) => ['projects', id] as const,
  projectMembers: (id: string) => ['projects', id, 'members'] as const,
  projectDocs: (id: string) => ['projects', id, 'docs'] as const,
  projectDoc: (projectId: string, docId: string) => ['projects', projectId, 'docs', docId] as const,
  projectDecisions: (id: string) => ['projects', id, 'decisions'] as const,
  projectConversations: (id: string) => ['projects', id, 'conversations'] as const,
  conversationProjects: (conversationId: string) => ['conversations', conversationId, 'projects'] as const,
}
