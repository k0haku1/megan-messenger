<template>
  <BaseModalSheet
    :open="profileUi.isOpen"
    :title="profileUi.title"
    title-id="profile-modal-title"
    scrollable
    @close="profileUi.close()"
  >
    <template #header-start>
      <BaseIconButton
        v-if="profileUi.canGoBack"
        label="Назад"
        variant="ghost"
        @click="profileUi.back()"
      >
        <AppIcon name="chevron-left" />
      </BaseIconButton>
    </template>

    <ProfileMainView v-if="profileUi.view === 'main'" />
    <ProfileEditView v-else-if="profileUi.view === 'edit'" />
    <ProfilePasswordView v-else-if="profileUi.view === 'password'" />
    <ProfileThemeView v-else-if="profileUi.view === 'theme'" />
    <ProfileUsernameView v-else-if="profileUi.view === 'username'" />
    <ProfilePrivacyView v-else-if="profileUi.view === 'privacy'" />
    <ProfilePrivacyFindView v-else-if="profileUi.view === 'privacy-find'" />
    <ProfilePrivacyDmView v-else-if="profileUi.view === 'privacy-dm'" />
  </BaseModalSheet>
</template>

<script setup lang="ts">
import AppIcon from '@/shared/ui/AppIcon.vue'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'
import { useProfileUiStore } from '../model/profile-ui.store'
import ProfileMainView from './ProfileMainView.vue'
import ProfileEditView from './ProfileEditView.vue'
import ProfilePasswordView from './ProfilePasswordView.vue'
import ProfileThemeView from './ProfileThemeView.vue'
import ProfileUsernameView from './ProfileUsernameView.vue'
import ProfilePrivacyView from './ProfilePrivacyView.vue'
import ProfilePrivacyFindView from './ProfilePrivacyFindView.vue'
import ProfilePrivacyDmView from './ProfilePrivacyDmView.vue'

const profileUi = useProfileUiStore()
</script>
