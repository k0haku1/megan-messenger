import { createRouter, createWebHistory } from 'vue-router'
import { useSessionStore } from '@/entities/session/model/session.store'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/auth',
      name: 'auth',
      component: () => import('@/pages/auth/ui/AuthPage.vue'),
      meta: { guest: true },
    },
    {
      path: '/',
      name: 'messenger',
      component: () => import('@/pages/messenger/ui/MessengerPage.vue'),
      meta: { requiresAuth: true, requiresOnboarding: true },
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

router.beforeEach(async (to) => {
  const session = useSessionStore()
  if (!session.isReady) {
    await session.bootstrap()
  }

  const authenticated = session.isAuthenticated
  const onboarded = authenticated && !session.needsUsername

  if (to.meta.requiresAuth && !authenticated) {
    if (to.fullPath === '/' || to.path === '/') {
      return { name: 'auth' }
    }
    return { name: 'auth', query: { redirect: to.fullPath } }
  }

  if (to.meta.requiresOnboarding && session.needsUsername) {
    return { name: 'auth' }
  }

  if (to.meta.guest && onboarded && !session.needsPassword) {
    return { name: 'messenger' }
  }

  return true
})
