<template>
  <div class="profile-view">
    <section class="profile-modal__identity">
      <BaseAvatar :name="displayName" size="lg" :color="avatarColor" />
      <div class="profile-modal__identity-meta">
        <strong>{{ displayName }}</strong>
        <span>{{ session.user?.phone }}</span>
      </div>
    </section>

    <section class="profile-view__section" aria-label="Настройки">
      <ProfileNavRow
        icon="username"
        label="Username"
        :value="usernameValue"
        @select="profileUi.openView('username')"
      />
      <ProfileNavRow
        icon="privacy"
        label="Конфиденциальность"
        :value="privacyStatus"
        @select="profileUi.openView('privacy')"
      />
      <ProfileNavRow
        icon="theme"
        label="Тема оформления"
        description="Светлая, тёмная или как в системе"
        :value="themeLabel"
        @select="profileUi.openView('theme')"
      />
      <ProfileNavRow
        icon="password"
        label="Облачный пароль"
        description="Дополнительная защита после SMS"
        :value="passwordStatus"
        @select="profileUi.openView('password')"
      />
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import { useSessionStore } from '@/entities/session/model/session.store'
import { getUserAvatarColor, getUserDisplayName } from '@/entities/session/lib/display'
import { useThemeStore } from '@/features/theme/model/theme.store'
import { privacyFindLabel, privacyDmLabel } from '../model/privacy-labels'
import { useProfileUiStore } from '../model/profile-ui.store'
import ProfileNavRow from './ProfileNavRow.vue'

const session = useSessionStore()
const profileUi = useProfileUiStore()
const themeStore = useThemeStore()

const displayName = computed(() =>
  getUserDisplayName(session.user?.username, session.user?.phone),
)
const avatarColor = computed(() => getUserAvatarColor(displayName.value))
const passwordStatus = computed(() => (session.user?.hasPassword ? 'Вкл' : 'Выкл'))
const themeLabel = computed(() => themeStore.label)
const usernameValue = computed(() => (session.user?.username ? `@${session.user.username}` : 'Не задан'))
const privacyStatus = computed(() => {
  const find = privacyFindLabel(session.user?.usernameSearchable)
  const dm = privacyDmLabel(session.user?.dmPolicy)
  if (find === 'Все' && dm === 'Все') return 'Все'
  return `${find} · ${dm}`
})
</script>
