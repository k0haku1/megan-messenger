import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { folderApi } from './folder.api'
import { patchConversationFolderMembership } from '@/entities/folder/lib/apply-folder-membership'
import { useSessionStore } from '@/entities/session/model/session.store'
import { queryClient } from '@/shared/api/query-client'
import { queryKeys } from '@/shared/api/query-keys'

export function useFolders() {
  const session = useSessionStore()
  const enabled = computed(() => session.isAuthenticated && !session.needsUsername)

  const query = useQuery({
    queryKey: queryKeys.folders,
    queryFn: folderApi.list,
    enabled,
  })

  const folders = computed(() => query.data.value?.folders ?? [])

  return { query, folders }
}

export function useCreateFolder() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: { name: string; icon?: string }) => folderApi.create(payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.folders }),
  })
}

export function useUpdateFolder() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ folderId, name, icon }: { folderId: string; name: string; icon?: string }) =>
      folderApi.update(folderId, { name, icon }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.folders }),
  })
}

export function useDeleteFolder() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (folderId: string) => folderApi.delete(folderId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.folders })
      void qc.invalidateQueries({ queryKey: queryKeys.conversations })
    },
  })
}

export function useAssignChatFolder() {
  return useMutation({
    mutationFn: ({ folderId, conversationIds }: { folderId: string; conversationIds: string[] }) =>
      folderApi.addItems(folderId, conversationIds),
    onSuccess: (_data, { folderId, conversationIds }) => {
      let patched = true
      for (const conversationId of conversationIds) {
        if (!patchConversationFolderMembership(conversationId, folderId, true)) {
          patched = false
        }
      }
      if (!patched) {
        void queryClient.invalidateQueries({ queryKey: queryKeys.conversations })
      }
    },
  })
}

export function useRemoveChatFromFolder() {
  return useMutation({
    mutationFn: ({ folderId, conversationId }: { folderId: string; conversationId: string }) =>
      folderApi.removeItem(folderId, conversationId),
    onSuccess: (_data, { folderId, conversationId }) => {
      if (!patchConversationFolderMembership(conversationId, folderId, false)) {
        void queryClient.invalidateQueries({ queryKey: queryKeys.conversations })
      }
    },
  })
}
