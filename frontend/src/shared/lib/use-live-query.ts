import { liveQuery } from 'dexie'
import { onScopeDispose, ref, type Ref } from 'vue'

/** Reactive reads from IndexedDB. Components must use this, never Query results. */
export function useLiveQuery<T>(read: () => Promise<T>, initial: T): Ref<T> {
  const value = ref(initial) as Ref<T>
  const subscription = liveQuery(read).subscribe({ next: (next) => { value.value = next } })
  onScopeDispose(() => subscription.unsubscribe())
  return value
}
