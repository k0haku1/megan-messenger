import { defineStore } from "pinia";
import { useLocalStorage } from "@vueuse/core";

/** Ephemeral interface state. Domain records remain in Dexie. */
export const useConversationSelectionStore = defineStore(
  "conversation-selection",
  () => {
    const selectedConversationId = useLocalStorage<string | null>(
      "megan:selected-conversation",
      null,
    );
    const activeFolder = useLocalStorage("megan:active-folder", "all");

    function select(conversationId: string | null): void {
      selectedConversationId.value = conversationId;
    }

    function selectFolder(folderId: string): void {
      activeFolder.value = folderId;
    }

    return { selectedConversationId, activeFolder, select, selectFolder };
  },
);
