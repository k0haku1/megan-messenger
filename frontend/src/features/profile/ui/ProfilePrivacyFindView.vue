<template>
  <div class="profile-view">
    <section class="profile-view__section" aria-label="Кто может найти меня">
      <button
        v-for="option in PRIVACY_FIND_OPTIONS"
        :key="String(option.value)"
        class="privacy-picker__option"
        :class="{ 'is-selected': selected === option.value }"
        type="button"
        :disabled="isSaving"
        @click="select(option.value)"
      >
        <span class="privacy-picker__body">
          <strong>{{ option.label }}</strong>
          <small>{{ option.description }}</small>
        </span>
        <AppIcon v-if="selected === option.value" class="privacy-picker__check" name="check" size="sm" />
      </button>
    </section>

    <p v-if="error" class="auth-form__error">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import { PRIVACY_FIND_OPTIONS } from '../model/privacy-labels'
import { usePrivacySettings } from '../model/use-privacy-settings'
import { useProfileUiStore } from '../model/profile-ui.store'

const profileUi = useProfileUiStore()
const { session, isSaving, error, save } = usePrivacySettings()

const selected = computed(() => session.user?.usernameSearchable !== false)

async function select(usernameSearchable: boolean) {
  if (usernameSearchable === selected.value) {
    profileUi.back()
    return
  }

  const dmPolicy = session.user?.dmPolicy ?? 'everyone'
  const ok = await save(usernameSearchable, dmPolicy)
  if (ok) profileUi.back()
}
</script>
