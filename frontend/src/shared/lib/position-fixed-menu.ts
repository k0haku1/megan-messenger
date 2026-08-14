export type MenuPoint = {
  x: number
  y: number
}

export function positionFixedMenu(menu: HTMLElement, point: MenuPoint, gap = 6): void {
  const padding = 8
  const rect = menu.getBoundingClientRect()

  let x = point.x
  let y = point.y + gap

  if (y + rect.height > window.innerHeight - padding) {
    y = point.y - rect.height - gap
  }
  if (x + rect.width > window.innerWidth - padding) {
    x = window.innerWidth - rect.width - padding
  }
  if (x < padding) x = padding
  if (y < padding) y = padding

  menu.style.left = `${x}px`
  menu.style.top = `${y}px`
}
