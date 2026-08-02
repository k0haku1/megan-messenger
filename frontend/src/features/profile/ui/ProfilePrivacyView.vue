<template>
  <div class="profile-view">
    <section class="profile-view__section" aria-label="Конфиденциальность">
      <ProfileNavRow
        label="Кто может найти меня"
        :value="findLabel"
        @select="profileUi.openView('privacy-find')"
      />
      <ProfileNavRow
        label="Кто может написать мне"
        :value="dmLabel"
        @select="profileUi.openView('privacy-dm')"
      />
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useSessionStore } from '@/entities/session/model/session.store'
import { privacyDmLabel, privacyFindLabel } from '../model/privacy-labels'
import { useProfileUiStore } from '../model/profile-ui.store'
import ProfileNavRow from './ProfileNavRow.vue'

const session = useSessionStore()
const profileUi = useProfileUiStore()

const findLabel = computed(() => privacyFindLabel(session.user?.usernameSearchable))
const dmLabel = computed(() => privacyDmLabel(session.user?.dmPolicy))
</script>
