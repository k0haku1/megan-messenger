<template>
  <div class="profile-view">
    <p class="profile-modal__hint">Выберите, как Megan Messenger выглядит на этом устройстве.</p>

    <section class="profile-view__section" aria-label="Тема оформления">
      <button
        v-for="option in options"
        :key="option.value"
        class="theme-option"
        :class="{ 'is-selected': themeStore.preference === option.value }"
        type="button"
        @click="themeStore.setPreference(option.value)"
      >
        <span class="theme-option__body">
          <strong>{{ option.label }}</strong>
          <small>{{ option.description }}</small>
        </span>
        <AppIcon
          v-if="themeStore.preference === option.value"
          class="theme-option__check"
          name="check"
          size="sm"
        />
      </button>
    </section>
  </div>
</template>

<script setup lang="ts">
import AppIcon from '@/shared/ui/AppIcon.vue'
import { useThemeStore } from '@/features/theme/model/theme.store'
import { THEME_LABELS } from '@/shared/theme/constants'
import type { ThemeMode } from '@/shared/theme/types'

const themeStore = useThemeStore()

const options: Array<{ value: ThemeMode; label: string; description: string }> = [
  {
    value: 'light',
    label: THEME_LABELS.light,
    description: 'Светлый интерфейс с фиолетовым акцентом',
  },
  {
    value: 'dark',
    label: THEME_LABELS.dark,
    description: 'Тёмный интерфейс, удобнее ночью',
  },
  {
    value: 'system',
    label: THEME_LABELS.system,
    description: 'Следовать настройкам операционной системы',
  },
]
</script>
