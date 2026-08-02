<template>
  <Teleport to="body">
    <div
      v-if="profileStore.isOpen"
      class="user-profile-sheet"
      role="dialog"
      aria-modal="true"
      aria-labelledby="user-profile-title"
      @click.self="profileStore.close()"
    >
      <div class="user-profile-sheet__panel">
        <header class="user-profile-sheet__header">
          <h2 id="user-profile-title">Профиль</h2>
          <BaseIconButton label="Закрыть" variant="ghost" @click="profileStore.close()">
            <AppIcon name="close" />
          </BaseIconButton>
        </header>

        <div v-if="profileStore.isLoading" class="user-profile-sheet__state">Загрузка…</div>

        <div v-else-if="profileStore.error" class="user-profile-sheet__state is-error">
          {{ profileStore.error }}
        </div>

        <template v-else-if="profileStore.profile">
          <section class="user-profile-sheet__identity">
            <BaseAvatar :name="profileStore.profile.username" size="lg" :color="0" />
            <div>
              <strong>@{{ profileStore.profile.username }}</strong>
              <span>Открытый профиль Megan</span>
            </div>
          </section>

          <BaseButton
            class="is-block"
            :disabled="!profileStore.profile.canMessage"
            @click="startChat()"
          >
            Написать сообщение
          </BaseButton>

          <p v-if="!profileStore.profile.canMessage" class="user-profile-sheet__hint">
            Пользователь не принимает новые сообщения.
          </p>
        </template>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import AppIcon from '@/shared/ui/AppIcon.vue'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseIconButton from '@/shared/ui/BaseIconButton.vue'
import { useUserProfileStore } from '@/features/user-profile/model/user-profile.store'
import { useGlobalSearchStore } from '@/features/global-search/model/global-search.store'
import { startDirectMessage } from '@/features/start-dm/model/start-dm'

const profileStore = useUserProfileStore()
const searchStore = useGlobalSearchStore()

async function startChat() {
  if (!profileStore.profile) return
  await startDirectMessage({
    userId: profileStore.profile.id,
    username: profileStore.profile.username,
    avatarUrl: profileStore.profile.avatarUrl,
  })
  profileStore.close()
  searchStore.clear()
}
</script>
