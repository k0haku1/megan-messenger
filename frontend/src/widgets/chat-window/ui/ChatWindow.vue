<template>
  <section class="chat-window" aria-label="Переписка">
    <template v-if="conversation">
      <header class="chat-header">
        <BaseAvatar
          :name="conversation.title ?? 'Личный диалог'"
          :color="conversation.type === 'group' ? 2 : 0"
        />
        <div class="chat-header__meta">
          <strong>{{ conversation.title ?? "Личный диалог" }}</strong>
          <span>{{
            conversation.type === "group" ? "групповой чат" : "был(а) недавно"
          }}</span>
        </div>
        <div class="chat-header__actions">
          <BaseIconButton label="Поиск по сообщениям"
            ><span class="action-glyph">⌕</span></BaseIconButton
          >
          <BaseIconButton label="Информация о чате"
            ><span class="action-glyph">⋮</span></BaseIconButton
          >
        </div>
      </header>

      <div class="message-area">
        <div v-if="visibleMessages.length === 0" class="chat-empty">
          <span class="chat-empty__icon">✉</span>
          <strong>Начните общение</strong>
          <span>Сообщения этого чата появятся здесь</span>
        </div>
        <div v-else class="message-stack">
          <article
            v-for="message in visibleMessages"
            :key="message.id"
            class="message-bubble"
          >
            <strong>{{ message.sender.username }}</strong>
            <p>{{ message.content }}</p>
            <time>{{
              new Date(message.createdAt).toLocaleTimeString("ru", {
                hour: "2-digit",
                minute: "2-digit",
              })
            }}</time>
          </article>
        </div>
      </div>

      <footer class="composer">
        <BaseIconButton label="Прикрепить файл"
          ><span class="action-glyph">⌕</span></BaseIconButton
        >
        <textarea
          v-model="draft"
          rows="1"
          placeholder="Сообщение"
          aria-label="Сообщение"
        />
        <BaseIconButton label="Эмодзи"
          ><span class="action-glyph">☺</span></BaseIconButton
        >
        <BaseIconButton label="Отправить" variant="primary"
          ><span class="send-glyph">➤</span></BaseIconButton
        >
      </footer>
    </template>

    <div v-else class="chat-placeholder">
      <div class="chat-placeholder__mark">M</div>
      <strong>Выберите чат</strong>
      <span>Откройте диалог слева, чтобы начать общение</span>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { storeToRefs } from "pinia";
import { db } from "@/shared/model/db";
import { useLiveQuery } from "@/shared/lib/use-live-query";
import { useConversationSelectionStore } from "@/features/conversation-selection/model/conversation-selection.store";
import BaseAvatar from "@/shared/ui/BaseAvatar.vue";
import BaseIconButton from "@/shared/ui/BaseIconButton.vue";

const navigation = useConversationSelectionStore();
const { selectedConversationId } = storeToRefs(navigation);
const conversations = useLiveQuery(() => db.conversations.toArray(), []);
const messages = useLiveQuery(
  () => db.messages.orderBy("createdAt").toArray(),
  [],
);
const draft = ref("");

const conversation = computed(() =>
  conversations.value.find(({ id }) => id === selectedConversationId.value),
);
const visibleMessages = computed(() =>
  messages.value.filter(
    ({ conversationId }) => conversationId === selectedConversationId.value,
  ),
);
</script>
