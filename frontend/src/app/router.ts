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
      path: '/u/:username',
      name: 'user-link',
      component: () => import('@/pages/messenger/ui/MessengerPage.vue'),
      meta: { requiresAuth: true, requiresOnboarding: true },
    },
    {
      path: '/',
      name: 'messenger',
      component: () => import('@/pages/messenger/ui/MessengerPage.vue'),
      meta: { requiresAuth: true, requiresOnboarding: true },
    },
    {
      path: '/projects',
      name: 'projects',
      component: () => import('@/pages/projects/ui/ProjectsPage.vue'),
      meta: { requiresAuth: true, requiresOnboarding: true },
    },
    {
      path: '/projects/:projectId',
      name: 'project-detail',
      component: () => import('@/pages/project-detail/ui/ProjectDetailPage.vue'),
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
