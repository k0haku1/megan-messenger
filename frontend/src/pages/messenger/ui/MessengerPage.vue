<template>
  <main class="messenger-layout">
    <FolderSidebar />
    <ConversationList />
    <ChatWindow />
    <ProfileModal />
    <UserProfileSheet />
  </main>
</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useConversations } from '@/entities/conversation/api/conversation.queries'
import ConversationList from '@/widgets/conversation-list/ui/ConversationList.vue'
import FolderSidebar from '@/widgets/folder-sidebar/ui/FolderSidebar.vue'
import ChatWindow from '@/widgets/chat-window/ui/ChatWindow.vue'
import ProfileModal from '@/features/profile/ui/ProfileModal.vue'
import UserProfileSheet from '@/features/user-profile/ui/UserProfileSheet.vue'
import { useUserProfileStore } from '@/features/user-profile/model/user-profile.store'

const route = useRoute()
const profileStore = useUserProfileStore()

useConversations()

function openProfileFromRoute() {
  const username = route.params.username
  if (typeof username === 'string' && username.length > 0) {
    void profileStore.open(username)
  }
}

onMounted(openProfileFromRoute)
watch(() => route.params.username, openProfileFromRoute)
</script>
