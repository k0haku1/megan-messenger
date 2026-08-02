export function getUserDisplayName(username?: string, phone?: string): string {
  return username || phone || 'User'
}

export function getUserAvatarColor(seed: string): number {
  let hash = 0
  for (let index = 0; index < seed.length; index += 1) {
    hash = (hash + seed.charCodeAt(index) * (index + 1)) % 6
  }
  return hash
}
