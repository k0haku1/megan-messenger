import { nextTick, onScopeDispose, ref, watch, type Ref } from 'vue'

const NEAR_BOTTOM_PX = 96

/**
 * Scroll policy:
 * - do not jump on conversation open/switch (keep reading position)
 * - follow when the user sends a message
 * - for incoming messages while scrolled up — side arrow + unseen counter
 */
export function useScrollToLatest(options: {
  container: Ref<HTMLElement | null>
  conversationKey: Ref<string | null | undefined>
  latestMessageId: Ref<string | null | undefined>
  /** True when the latest message was sent by the current user. */
  latestIsOwn: Ref<boolean>
}) {
  const stickToBottom = ref(true)
  const showJumpToLatest = ref(false)
  const unseenCount = ref(0)

  function isNearBottom(el: HTMLElement): boolean {
    const distance = el.scrollHeight - el.scrollTop - el.clientHeight
    return distance <= NEAR_BOTTOM_PX
  }

  function syncJumpVisibility(el: HTMLElement): void {
    const near = isNearBottom(el)
    stickToBottom.value = near
    showJumpToLatest.value = !near && el.scrollHeight > el.clientHeight + 8
    if (near) {
      unseenCount.value = 0
    }
  }

  function scrollToLatest(behavior: ScrollBehavior = 'smooth'): void {
    const el = options.container.value
    if (!el) return
    el.scrollTo({ top: el.scrollHeight, behavior })
    stickToBottom.value = true
    showJumpToLatest.value = false
    unseenCount.value = 0
  }

  function onContainerScroll(): void {
    const el = options.container.value
    if (!el) return
    syncJumpVisibility(el)
  }

  async function jumpNow(behavior: ScrollBehavior = 'auto'): Promise<void> {
    await nextTick()
    requestAnimationFrame(() => {
      scrollToLatest(behavior)
    })
  }

  watch(
    () => options.conversationKey.value,
    () => {
      showJumpToLatest.value = false
      stickToBottom.value = true
      unseenCount.value = 0
      void nextTick(() => {
        const el = options.container.value
        if (el) syncJumpVisibility(el)
      })
    },
  )

  watch(
    () => options.latestMessageId.value,
    (id, prev) => {
      if (!id || id === prev) return

      if (options.latestIsOwn.value) {
        unseenCount.value = 0
        void jumpNow(prev ? 'smooth' : 'auto')
        return
      }

      const el = options.container.value
      if (!el) {
        stickToBottom.value = false
        showJumpToLatest.value = true
        unseenCount.value += 1
        return
      }

      if (isNearBottom(el)) {
        syncJumpVisibility(el)
        return
      }

      stickToBottom.value = false
      showJumpToLatest.value = true
      unseenCount.value += 1
    },
    { flush: 'post' },
  )

  onScopeDispose(() => {
    stickToBottom.value = true
    showJumpToLatest.value = false
    unseenCount.value = 0
  })

  return {
    stickToBottom,
    showJumpToLatest,
    unseenCount,
    onContainerScroll,
    scrollToLatest,
  }
}
