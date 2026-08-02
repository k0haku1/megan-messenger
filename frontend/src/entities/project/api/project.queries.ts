import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { projectApi } from '@/entities/project/api/project.api'
import { queryKeys } from '@/shared/api/query-keys'

export function useProjects() {
  const query = useQuery({
    queryKey: queryKeys.projects,
    queryFn: async () => (await projectApi.list()).projects,
  })

  return {
    projects: computed(() => query.data.value ?? []),
    isLoading: query.isLoading,
    isError: query.isError,
    refetch: query.refetch,
  }
}

export function useProject(projectId: Ref<string>) {
  return useQuery({
    queryKey: computed(() => queryKeys.project(projectId.value)),
    queryFn: async () => projectApi.get(projectId.value),
    enabled: computed(() => projectId.value.length > 0),
  })
}

export function useProjectMembers(projectId: Ref<string>) {
  return useQuery({
    queryKey: computed(() => queryKeys.projectMembers(projectId.value)),
    queryFn: async () => (await projectApi.listMembers(projectId.value)).members,
    enabled: computed(() => projectId.value.length > 0),
  })
}

export function useProjectDocs(projectId: Ref<string>) {
  return useQuery({
    queryKey: computed(() => queryKeys.projectDocs(projectId.value)),
    queryFn: async () => (await projectApi.listDocs(projectId.value)).docs,
    enabled: computed(() => projectId.value.length > 0),
  })
}

export function useProjectDoc(projectId: Ref<string>, docId: Ref<string | null>) {
  return useQuery({
    queryKey: computed(() => queryKeys.projectDoc(projectId.value, docId.value ?? '')),
    queryFn: async () => (await projectApi.getDoc(projectId.value, docId.value!)).doc,
    enabled: computed(() => projectId.value.length > 0 && !!docId.value),
  })
}

export function useProjectDecisions(projectId: Ref<string>) {
  return useQuery({
    queryKey: computed(() => queryKeys.projectDecisions(projectId.value)),
    queryFn: async () => (await projectApi.listDecisions(projectId.value)).decisions,
    enabled: computed(() => projectId.value.length > 0),
  })
}

export function useProjectConversations(projectId: Ref<string>) {
  return useQuery({
    queryKey: computed(() => queryKeys.projectConversations(projectId.value)),
    queryFn: async () => (await projectApi.listConversations(projectId.value)).conversations,
    enabled: computed(() => projectId.value.length > 0),
  })
}

export function useInvalidateProject(projectId: string) {
  const queryClient = useQueryClient()

  return () => {
    void queryClient.invalidateQueries({ queryKey: queryKeys.projects })
    void queryClient.invalidateQueries({ queryKey: queryKeys.project(projectId) })
    void queryClient.invalidateQueries({ queryKey: queryKeys.projectDocs(projectId) })
    void queryClient.invalidateQueries({ queryKey: queryKeys.projectDecisions(projectId) })
    void queryClient.invalidateQueries({ queryKey: queryKeys.projectConversations(projectId) })
    void queryClient.invalidateQueries({ queryKey: queryKeys.projectMembers(projectId) })
  }
}

export function useCreateProject() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: projectApi.create,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: queryKeys.projects })
    },
  })
}
