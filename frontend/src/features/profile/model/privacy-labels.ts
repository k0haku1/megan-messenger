import type { DMPolicy } from '@/entities/user/model/types'

export const PRIVACY_FIND_OPTIONS = [
  { value: true, label: 'Все', description: 'Вас можно найти по @username в поиске' },
  { value: false, label: 'Никто', description: '@username не показывается в поиске' },
] as const

export const PRIVACY_DM_OPTIONS = [
  { value: 'everyone' as DMPolicy, label: 'Все', description: 'Любой может начать диалог' },
  { value: 'nobody' as DMPolicy, label: 'Никто', description: 'Новые диалоги недоступны' },
] as const

export function privacyFindLabel(searchable: boolean | undefined): string {
  return searchable === false ? 'Никто' : 'Все'
}

export function privacyDmLabel(policy: DMPolicy | undefined): string {
  return policy === 'nobody' ? 'Никто' : 'Все'
}
