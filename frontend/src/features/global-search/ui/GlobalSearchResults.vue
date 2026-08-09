<template>
  <section v-if="hasResults" class="global-search-results" aria-label="Результаты поиска">
    <header class="global-search-results__heading">Пользователи</header>

    <button
      v-for="(user, index) in searchStore.userResults"
      :key="user.id"
      class="global-search-results__item"
      type="button"
      @click="openProfile(user.username)"
    >
      <BaseAvatar :name="user.username" :src="user.avatarUrl" :color="index" size="md" />
      <span class="global-search-results__meta">
        <strong>@{{ user.username }}</strong>
        <small>Написать сообщение</small>
      </span>
    </button>

    <p v-if="searchStore.isSearching" class="global-search-results__hint">Ищем пользователей…</p>
    <p v-else-if="searchStore.searchError" class="global-search-results__error">
      {{ searchStore.searchError }}
    </p>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BaseAvatar from '@/shared/ui/BaseAvatar.vue'
import { normalizeUsernameQuery } from '@/entities/user/lib/username'
import { useGlobalSearchStore } from '@/features/global-search/model/global-search.store'
import { useUserProfileStore } from '@/features/user-profile/model/user-profile.store'

const searchStore = useGlobalSearchStore()
const profileStore = useUserProfileStore()

const hasResults = computed(
  () =>
    normalizeUsernameQuery(searchStore.query).length >= 2
    && (searchStore.userResults.length > 0
      || searchStore.isSearching
      || Boolean(searchStore.searchError)),
)

function openProfile(username: string) {
  void profileStore.open(username)
}
</script>
