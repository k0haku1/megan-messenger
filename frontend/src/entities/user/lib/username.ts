export const USERNAME_HINT = 'По нему вас могут найти остальные пользователи'

export function normalizeUsernameQuery(raw: string): string {
  return raw.trim().replace(/^@+/, '').toLowerCase()
}
