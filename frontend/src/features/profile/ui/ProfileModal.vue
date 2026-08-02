<template>
  <Teleport to="body">
    <div
      v-if="profileUi.isOpen"
      class="profile-modal"
      role="dialog"
      aria-modal="true"
      aria-labelledby="profile-modal-title"
      @click.self="profileUi.close()"
    >
      <div class="profile-modal__panel">
        <header class="profile-modal__header">
          <div class="profile-modal__header-start">
            <BaseIconButton
              v-if="profileUi.canGoBack"
              label="Назад"
              variant="ghost"
              @click="profileUi.back()"
            >
              <AppIcon name="chevron-left" />
            </BaseIconButton>
            <h2 id="profile-modal-title">{{ profileUi.title }}</h2>
          </div>
          <BaseIconButton label="Закрыть" variant="ghost" @click="profileUi.close()">
            <AppIcon name="close" />
          </BaseIconButton>
        </header>

        <ProfileMainView v-if="profileUi.view === 'main'" />
        <ProfilePasswordView v-else-if="profileUi.view === 'password'" />
        <ProfileThemeView v-else-if="profileUi.view === 'theme'" />
        <ProfileUsernameView v-else-if="profileUi.view === 'username'" />
        <ProfilePrivacyView v-else-if="profileUi.view === 'privacy'" />
        <ProfilePrivacyFindView v-else-if="profileUi.view === 'privacy-find'" />
        <ProfilePrivacyDmView v-else-if="profileUi.view === 'privacy-dm'" />
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import AppIcon from '@/shared/ui/AppIcon.vue'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import { useProfileUiStore } from '../model/profile-ui.store'
import ProfileMainView from './ProfileMainView.vue'
import ProfilePasswordView from './ProfilePasswordView.vue'
import ProfileThemeView from './ProfileThemeView.vue'
import ProfileUsernameView from './ProfileUsernameView.vue'
import ProfilePrivacyView from './ProfilePrivacyView.vue'
import ProfilePrivacyFindView from './ProfilePrivacyFindView.vue'
import ProfilePrivacyDmView from './ProfilePrivacyDmView.vue'

const profileUi = useProfileUiStore()
</script>
