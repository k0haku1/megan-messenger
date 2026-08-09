export const deleteMessagesCopy = {
  cancel: 'Отмена',
  confirm: 'Удалить',
  forEveryone: 'Удалить для всех',
  onlyForYou: 'Выбранные сообщения будут удалены только для вас',
} as const

export function deleteMessagesTitle(count: number): string {
  if (count <= 1) return 'Вы действительно хотите удалить сообщение?'
  return `Вы действительно хотите удалить ${count} ${pluralizeMessages(count)}?`
}

function pluralizeMessages(count: number): string {
  const mod10 = count % 10
  const mod100 = count % 100
  if (mod10 === 1 && mod100 !== 11) return 'сообщение'
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14)) return 'сообщения'
  return 'сообщений'
}
