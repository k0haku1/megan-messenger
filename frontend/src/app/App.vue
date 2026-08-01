<template>
  <main class="messenger-layout">
    <FolderSidebar />
    <ConversationList />
    <ChatWindow />
  </main>
</template>

<script setup lang="ts">
import { useConversationsSync } from "@/entities/conversation/api/conversation.queries";
import ConversationList from "@/widgets/conversation-list/ui/ConversationList.vue";
import FolderSidebar from "@/widgets/folder-sidebar/ui/FolderSidebar.vue";
import ChatWindow from "@/widgets/chat-window/ui/ChatWindow.vue";

useConversationsSync();
</script>
