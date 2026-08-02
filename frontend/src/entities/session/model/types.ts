export type SessionUser = {
  id: string
  phone: string
  username?: string
  avatarUrl?: string
  hasPassword: boolean
  onboardingComplete: boolean
  usernameSearchable?: boolean
  dmPolicy?: 'everyone' | 'nobody'
}

export type AuthResult = {
  accessToken?: string
  refreshToken?: string
  needUsername: boolean
  needPassword: boolean
  challengeToken?: string
}

export type CompleteUsernameResponse = AuthResult & {
  user: SessionUser
}
