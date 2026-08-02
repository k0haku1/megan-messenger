import type { AppIconName } from '@/shared/ui/app-icon.types'

export type FolderId = 'all' | 'personal' | 'groups' | 'work'

export interface FolderConfig {
  id: FolderId
  label: string
  icon: AppIconName
}

export const CHAT_FOLDERS: FolderConfig[] = [
  { id: 'all', icon: 'folder-all', label: 'Все чаты' },
  { id: 'personal', icon: 'folder-personal', label: 'Личные' },
  { id: 'groups', icon: 'folder-groups', label: 'Группы' },
  { id: 'work', icon: 'folder-work', label: 'Работа' },
]

export const CHAT_FOLDER_TITLES: Record<FolderId, string> = Object.fromEntries(
  CHAT_FOLDERS.map((folder) => [folder.id, folder.label]),
) as Record<FolderId, string>
