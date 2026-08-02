<template>
  <BaseModalSheet
    :open="profileStore.isOpen"
    title="Профиль"
    title-id="user-profile-title"
    compact
    @close="profileStore.close()"
  >
    <div v-if="profileStore.isLoading" class="modal-sheet__state">Загрузка…</div>

    <div v-else-if="profileStore.error" class="modal-sheet__state is-error">
      {{ profileStore.error }}
    </div>

    <template v-else-if="profileStore.profile">
      <section class="modal-sheet__identity">
        <BaseAvatar :name="profileStore.profile.username" size="lg" :color="0" />
        <div class="modal-sheet__identity-meta">
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

      <p v-if="!profileStore.profile.canMessage" class="modal-sheet__hint">
        Пользователь не принимает новые сообщения.
      </p>
    </template>
  </BaseModalSheet>
</template>

<script setup lang="ts">
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import BaseButton from '@/shared/ui/BaseButton.vue'
import BaseModalSheet from '@/shared/ui/BaseModalSheet.vue'
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
