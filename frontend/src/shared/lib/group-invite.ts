export function parseGroupInviteSlug(input: string): string | null {
  const trimmed = input.trim().toLowerCase()
  if (/^[a-z0-9]{11}$/.test(trimmed)) {
    return trimmed
  }

  const match = trimmed.match(/([a-z0-9]{11})/)
  return match?.[1] ?? null
}

export async function copyGroupInviteSlug(slug: string): Promise<void> {
  await navigator.clipboard.writeText(slug)
}
