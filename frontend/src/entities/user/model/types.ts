export interface PublicUser {
  id: string
  username: string
  avatarUrl?: string
  canMessage: boolean
}

export interface UserSearchResult {
  id: string
  username: string
  avatarUrl?: string
}

export type DMPolicy = 'everyone' | 'nobody'

export interface UserPrivacySettings {
  usernameSearchable: boolean
  dmPolicy: DMPolicy
}
