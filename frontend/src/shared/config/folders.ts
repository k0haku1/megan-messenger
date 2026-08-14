export const ALL_FOLDER_ID = 'all'

export interface SystemFolderConfig {
  id: typeof ALL_FOLDER_ID
  label: string
}

export const ALL_FOLDER: SystemFolderConfig = { id: ALL_FOLDER_ID, label: 'Все чаты' }

export type ActiveFolderId = typeof ALL_FOLDER_ID | string
