import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin } from '@tanstack/vue-query'
import App from './App.vue'
import { queryClient } from '@/shared/api/query-client'
import '@/app/styles/index.css'

createApp(App).use(createPinia()).use(VueQueryPlugin, { queryClient }).mount('#app')
