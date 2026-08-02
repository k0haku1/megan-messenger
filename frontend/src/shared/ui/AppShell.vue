<template>
  <main class="messenger-layout">
    <FolderSidebar />
    <slot />
    <ProfileModal />
    <UserProfileSheet />
  </main>
</template>

<script setup lang="ts">
import FolderSidebar from '@/widgets/folder-sidebar/ui/FolderSidebar.vue'
import ProfileModal from '@/features/profile/ui/ProfileModal.vue'
import UserProfileSheet from '@/features/user-profile/ui/UserProfileSheet.vue'
</script>
