export const APP_ICON_NAMES = [
  'chevron-left',
  'close',
  'check',
  'edit',
  'username',
  'privacy',
  'theme',
  'password',
] as const

export type AppIconName = (typeof APP_ICON_NAMES)[number]
