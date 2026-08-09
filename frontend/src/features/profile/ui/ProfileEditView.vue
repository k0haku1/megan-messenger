<template>
  <div class="profile-view">
    <section class="profile-edit__avatar">
      <BaseAvatar
        :name="displayName"
        :src="session.user?.avatarUrl"
        size="lg"
        :color="avatarColor"
      />
      <label class="profile-edit__avatar-action">
        <span>{{ uploadingAvatar ? 'Загрузка…' : 'Сменить фото' }}</span>
        <input
          type="file"
          accept="image/jpeg,image/png,image/webp,image/gif"
          :disabled="uploadingAvatar"
          @change="onAvatarPick"
        />
      </label>
      <p v-if="avatarError" class="profile-edit__error">{{ avatarError }}</p>
      <p v-else-if="avatarSuccess" class="modal-sheet__success">{{ avatarSuccess }}</p>
    </section>

    <p class="modal-sheet__hint">
      Фото будет видно другим участникам в чатах и в вашем профиле.
    </p>

    <ProfileNavRow
      icon="username"
      label="Username"
      :value="usernameValue"
      @select="profileUi.openView('username')"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import { useSessionStore } from '@/entities/session/model/session.store'
import { getUserAvatarColor, getUserDisplayName } from '@/entities/session/lib/display'
import { userApi } from '@/entities/user/api/user.api'
import { getApiErrorMessage } from '@/shared/lib/get-api-error-message'
import { useProfileUiStore } from '../model/profile-ui.store'
import ProfileNavRow from './ProfileNavRow.vue'

const session = useSessionStore()
const profileUi = useProfileUiStore()

const uploadingAvatar = ref(false)
const avatarError = ref('')
const avatarSuccess = ref('')

const displayName = computed(() =>
  getUserDisplayName(session.user?.username, session.user?.phone),
)
const avatarColor = computed(() => getUserAvatarColor(displayName.value))
const usernameValue = computed(() =>
  session.user?.username ? `@${session.user.username}` : 'Не задан',
)

async function onAvatarPick(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

  uploadingAvatar.value = true
  avatarError.value = ''
  avatarSuccess.value = ''
  try {
    const { avatarUrl } = await userApi.uploadAvatar(file)
    if (session.user) {
      session.user.avatarUrl = avatarUrl
    }
    avatarSuccess.value = 'Фото обновлено'
  } catch (err) {
    avatarError.value = getApiErrorMessage(err, { fallback: 'Не удалось обновить фото' })
  } finally {
    uploadingAvatar.value = false
  }
}
</script>

<style scoped>
.profile-edit__avatar {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}
.profile-edit__avatar-action {
  position: relative;
  display: inline-flex;
  color: var(--color-accent);
  cursor: pointer;
  font-size: var(--font-size-sm);
  font-weight: 600;
}
.profile-edit__avatar-action input {
  position: absolute;
  inset: 0;
  opacity: 0;
  cursor: pointer;
}
.profile-edit__error {
  margin: 0;
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}
</style>
