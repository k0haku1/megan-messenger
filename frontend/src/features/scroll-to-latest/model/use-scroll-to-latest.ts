import { computed, nextTick, onBeforeUnmount, ref, unref, watch, type Ref } from 'vue'

const NEAR_BOTTOM_PX = 80

type ScrollMessage = {
  id: string
  createdAt: string
  senderId: string
}

type Options = {
  container: Ref<HTMLElement | null>
  conversationKey: Ref<string | null | undefined> | (() => string | null | undefined)
  latestMessageId: Ref<string | null | undefined> | (() => string | null | undefined)
  latestIsOwn?: Ref<boolean> | (() => boolean)
  unreadCount?: Ref<number> | (() => number)
  messages?: Ref<ScrollMessage[]> | (() => ScrollMessage[])
  currentUserId?: Ref<string | null | undefined> | (() => string | null | undefined)
  lastReadAt?: Ref<string | null | undefined> | (() => string | null | undefined)
  readStateReady?: Ref<boolean> | (() => boolean)
  /** Initial page written to Dexie. */
  messagesSynced?: Ref<boolean> | (() => boolean)
  /** True after older pages needed for first-unread are loaded. */
  openHistoryReady?: Ref<boolean> | (() => boolean)
  syncedLatestMessageId?: Ref<string | null | undefined> | (() => string | null | undefined)
  syncedMessageCount?: Ref<number> | (() => number)
  /** Called when unseen baseline should clear (near bottom / jump). */
  onCaughtUp?: () => void
}

function resolve<T>(source: Ref<T> | (() => T)): T {
  return typeof source === 'function' ? source() : unref(source)
}

function isNearBottom(el: HTMLElement) {
  return el.scrollHeight - el.scrollTop - el.clientHeight <= NEAR_BOTTOM_PX
}

function findFirstUnreadId(
  messages: ScrollMessage[],
  currentUserId: string | null | undefined,
  lastReadAt: string | null | undefined,
): string | null {
  const watermark = lastReadAt ? new Date(lastReadAt).getTime() : Number.NEGATIVE_INFINITY
  if (Number.isNaN(watermark)) return null

  for (const message of messages) {
    if (currentUserId && message.senderId === currentUserId) continue
    const created = new Date(message.createdAt).getTime()
    if (Number.isNaN(created)) continue
    if (created > watermark) return message.id
  }
  return null
}

export function useScrollToLatest(options: Options) {
  const stickToBottom = ref(false)
  const showJumpToLatest = ref(false)
  const unseenSinceLeave = ref(0)
  const pendingOpenRestore = ref(true)
  const programmaticScroll = ref(false)
  const openPositionReady = computed(() => !pendingOpenRestore.value)

  let revealTimer: ReturnType<typeof setTimeout> | null = null

  const unreadBaseline = computed(() => Math.max(0, resolve(options.unreadCount ?? (() => 0))))
  const unseenCount = computed(() => Math.max(unreadBaseline.value, unseenSinceLeave.value))

  function clearRevealTimer() {
    if (!revealTimer) return
    clearTimeout(revealTimer)
    revealTimer = null
  }

  function finishOpenRestore() {
    pendingOpenRestore.value = false
    clearRevealTimer()
  }

  function markCaughtUp() {
    unseenSinceLeave.value = 0
    options.onCaughtUp?.()
  }

  function updateFlags() {
    if (programmaticScroll.value || pendingOpenRestore.value) return
    const el = options.container.value
    if (!el) return
    const near = isNearBottom(el)
    stickToBottom.value = near
    showJumpToLatest.value = !near
    if (near) markCaughtUp()
  }

  function onContainerScroll() {
    updateFlags()
  }

  function scrollToLatest(behavior: ScrollBehavior = 'smooth') {
    const el = options.container.value
    if (!el) return
    programmaticScroll.value = true
    el.scrollTo({ top: el.scrollHeight, behavior })
    stickToBottom.value = true
    showJumpToLatest.value = false
    markCaughtUp()
    window.setTimeout(() => {
      programmaticScroll.value = false
      updateFlags()
    }, behavior === 'smooth' ? 350 : 0)
  }

  function scrollToMessage(messageId: string, behavior: ScrollBehavior = 'auto') {
    const el = options.container.value
    if (!el) return false
    const target = el.querySelector(`[data-message-id="${CSS.escape(messageId)}"]`) as HTMLElement | null
    if (!target) return false
    programmaticScroll.value = true
    const top = el.scrollTop + (target.getBoundingClientRect().top - el.getBoundingClientRect().top) - 12
    el.scrollTop = Math.max(0, top)
    if (behavior === 'smooth') {
      el.scrollTo({ top: Math.max(0, top), behavior })
    }
    window.setTimeout(() => {
      programmaticScroll.value = false
      updateFlags()
    }, behavior === 'smooth' ? 350 : 0)
    return true
  }

  function localListCaughtUp(): boolean {
    if (!options.messagesSynced || !resolve(options.messagesSynced)) return false
    const expectedCount = resolve(options.syncedMessageCount ?? (() => 0))
    const expectedLatest = resolve(options.syncedLatestMessageId ?? (() => null))
    const messages = resolve(options.messages ?? (() => []))

    if (!expectedLatest && expectedCount === 0) return true
    if (expectedCount > 0 && messages.length === 0) return false
    if (expectedLatest && !messages.some((message) => message.id === expectedLatest)) return false
    return true
  }

  function canRestore(): boolean {
    if (!pendingOpenRestore.value) return false
    if (!options.container.value) return false
    if (options.readStateReady && !resolve(options.readStateReady)) return false
    if (options.messagesSynced && !localListCaughtUp()) return false
    if (options.openHistoryReady && !resolve(options.openHistoryReady)) return false
    return true
  }

  function applyOpenPosition() {
    const container = options.container.value
    if (!container || !pendingOpenRestore.value) return false

    const messages = resolve(options.messages ?? (() => []))
    const firstUnreadId = findFirstUnreadId(
      messages,
      resolve(options.currentUserId ?? (() => null)),
      resolve(options.lastReadAt ?? (() => null)),
    )

    if (firstUnreadId) {
      if (!scrollToMessage(firstUnreadId, 'auto')) return false
      stickToBottom.value = isNearBottom(container)
      showJumpToLatest.value = !stickToBottom.value
      if (stickToBottom.value) markCaughtUp()
    } else {
      programmaticScroll.value = true
      container.scrollTop = container.scrollHeight
      stickToBottom.value = true
      showJumpToLatest.value = false
      markCaughtUp()
      programmaticScroll.value = false
    }

    finishOpenRestore()
    return true
  }

  function restoreOpenPosition() {
    if (!canRestore()) return

    // Position before paint when possible; retry once if DOM nodes lag Dexie.
    void nextTick(() => {
      if (applyOpenPosition()) return
      requestAnimationFrame(() => {
        void nextTick(() => {
          applyOpenPosition()
        })
      })
    })
  }

  watch(
    () => resolve(options.conversationKey),
    () => {
      stickToBottom.value = false
      showJumpToLatest.value = false
      unseenSinceLeave.value = 0
      pendingOpenRestore.value = true
      programmaticScroll.value = false
      clearRevealTimer()
      // Failsafe: never leave the thread invisible if restore stalls.
      revealTimer = setTimeout(() => {
        if (!pendingOpenRestore.value) return
        const el = options.container.value
        if (el) el.scrollTop = el.scrollHeight
        stickToBottom.value = true
        showJumpToLatest.value = false
        finishOpenRestore()
      }, 2500)
    },
  )

  watch(
    () => [
      resolve(options.conversationKey),
      resolve(options.messages ?? (() => [])).length,
      resolve(options.messages ?? (() => [])).at(-1)?.id ?? null,
      resolve(options.lastReadAt ?? (() => null)),
      options.readStateReady ? resolve(options.readStateReady) : true,
      options.messagesSynced ? resolve(options.messagesSynced) : true,
      options.openHistoryReady ? resolve(options.openHistoryReady) : true,
      resolve(options.syncedLatestMessageId ?? (() => null)),
      resolve(options.syncedMessageCount ?? (() => 0)),
      options.container.value,
    ],
    () => {
      if (!canRestore()) return
      restoreOpenPosition()
    },
  )

  watch(
    () => resolve(options.latestMessageId),
    (id, prev) => {
      if (!id || id === prev) return
      if (pendingOpenRestore.value) return

      const own = resolve(options.latestIsOwn ?? (() => false))
      if (own || stickToBottom.value) {
        void nextTick(() => scrollToLatest('auto'))
        return
      }
      unseenSinceLeave.value += 1
      showJumpToLatest.value = true
    },
  )

  let resizeObserver: ResizeObserver | null = null
  watch(
    () => options.container.value,
    (el, _, onCleanup) => {
      resizeObserver?.disconnect()
      resizeObserver = null
      if (!el) return

      const onScroll = () => updateFlags()
      el.addEventListener('scroll', onScroll, { passive: true })
      resizeObserver = new ResizeObserver(() => {
        if (pendingOpenRestore.value || programmaticScroll.value) return
        if (stickToBottom.value) scrollToLatest('auto')
        else updateFlags()
      })
      resizeObserver.observe(el)

      onCleanup(() => {
        el.removeEventListener('scroll', onScroll)
        resizeObserver?.disconnect()
        resizeObserver = null
      })
    },
    { immediate: true },
  )

  onBeforeUnmount(() => {
    clearRevealTimer()
    resizeObserver?.disconnect()
  })

  return {
    stickToBottom,
    showJumpToLatest,
    unseenCount,
    openPositionReady,
    onContainerScroll,
    scrollToLatest,
  }
}
