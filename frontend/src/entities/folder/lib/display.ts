import { ALL_FOLDER } from '@/shared/config/folders'
import type { ChatFolder } from '@/entities/folder/model/types'

export function getFolderInitial(name: string): string {
  const trimmed = name.trim()
  return trimmed ? trimmed.charAt(0).toUpperCase() : '?'
}

export function getActiveFolderTitle(activeFolderId: string, folders: ChatFolder[]): string {
  if (activeFolderId === ALL_FOLDER.id) return ALL_FOLDER.label
  return folders.find((folder) => folder.id === activeFolderId)?.name ?? 'Папка'
}
