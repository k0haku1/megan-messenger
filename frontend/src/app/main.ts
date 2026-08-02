import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin } from '@tanstack/vue-query'
import App from './App.vue'
import { router } from './router'
import { queryClient } from '@/shared/api/query-client'
import { useThemeStore } from '@/features/theme/model/theme.store'
import '@/app/styles/index.css'

const pinia = createPinia()
const app = createApp(App)

app.use(pinia).use(router).use(VueQueryPlugin, { queryClient })

useThemeStore(pinia).init()

app.mount('#app')
