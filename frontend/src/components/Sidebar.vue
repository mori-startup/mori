<template>
  <div :class="['sidebar', { 'sidebar--active': isActive }]">
    <div class="sidebar-content">
      <!-- Icônes du haut (accès direct) -->
      <ul class="icon-container">
        <li @click="navigateToMessages" class="icon-wrapper">
          <div class="icon-circle">
            <img src="@/assets/icons/messages.png" alt="Messagerie" />
          </div>
          <span>Messages</span>
        </li>

        <li @click="navigateToChatBot" class="icon-wrapper">
          <div class="icon-circle">
            <img src="@/assets/icons/chat.png" alt="Chat" />
          </div>
          <span>Mori Chatbot</span>
        </li>
      </ul>

      <!-- Zone d'affichage des contacts (amis + groupes) -->
      <ContactsForChatBotView
        v-if="activeView === 'contacts'"
        @select-contact="handleContactSelection"
      /> -->
      <ChatHistory />
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import ContactsForChatBotView from "./ContactsForChatBoxView.vue";
import ChatHistory from "./ChatHistory.vue";

export default {
  name: "Sidebar",
  props: {
    isActive: {
      type: Boolean,
      required: true,
    },
    // Si tu n'utilises plus ce tableau "contactsList" directement,
    // tu peux le laisser ou le retirer selon ta logique
    contactsList: {
      type: Array,
      required: true,
    },
  },
  data() {
    if (this.$route.name === "messages") {
      return {
        activeView: "contacts",
      };
    }else{
      return {
        activeView: null,
      };
    }
  },
  components: { ContactsForChatBotView,ChatHistory },
  computed: {
    // Map additional state so we can determine the most recent conversation.
    ...mapState({
      conversationsMsg: (state) => state.conversationsMsg,
    }),
  },
  methods: {
    async navigateToMessages() {
      try {
        // Check if conversation data is valid
        if (Array.isArray(this.conversationsMsg) && this.conversationsMsg.length > 0) {
          // Create a shallow copy and sort descending by lastMessageTime.
          const convs = [...this.conversationsMsg];
          convs.sort(
            (a, b) => new Date(b.lastMessageTime) - new Date(a.lastMessageTime)
          );
          const recentConv = convs[0];
          this.activeView = "contacts";
          await this.$router.push({
            name: "messages",
            query: {
              name: recentConv.nickname || recentConv.name,
              receiverId: recentConv.id,
              type: recentConv.type,
            },
          });
        } else if (Array.isArray(this.contactsList) && this.contactsList.length > 0) {
          // Fallback: if no conversation data, navigate to the first contact.
          this.activeView = "contacts";
          const firstContact = this.contactsList[0];
          await this.$router.push({
            name: "messages",
            query: {
              name: firstContact.nickname,
              receiverId: firstContact.id,
              type: "PERSON",
            },
          });
        } else {
          // If no contacts and no conversations, display an error or notification.
          this.activeView = "contacts";
          console.error("No conversations or contacts available.");
          // Optionally, you can set a flag to display a message in your UI.
        }
      } catch (error) {
        console.error("Error navigating to messages:", error);
      }
    },

    async navigateToChatBot() {
      try {
        await this.$router.push({ name: "mainpage" });
      } catch (error) {
        console.error("Error navigating to ChatBot:", error);
      }
    },

    handleContactSelection({ id, name, type }) {
      // Navigate to the selected conversation.
      this.$router.push({
        name: "messages",
        query: {
          name,
          receiverId: id,
          type,
        },
      });
    },
    handleConversationSelection(conversation) {
      this.$router.push({
        name: "messages",
        query: {
          name: conversation.nickname,
          receiverId: conversation.id,
          type: "PERSON",
        },
      });
    },
  },
};
</script>

<style scoped>
.sidebar {
  position: fixed;
  top: 64.45px;
  left: -440px;
  width: 440px;
  height: calc(100% - 64.45px);
  background-color: var(--bg-neutral);
  transition: left 0.3s ease;
  z-index: 2;
  overflow-y: auto;
}

.sidebar--active {
  left: 0;
}

.sidebar-content {
  padding: 20px;
}

.icon-container {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 20px;
}

.icon-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  text-align: center;
  color: var(--text-color);
  transition: transform 0.3s ease;
}

.icon-wrapper:hover {
  transform: scale(1.1);
}

.icon-circle {
  width: 100px;
  height: 70px;
  background-color: var(--purple-color);
  border-radius: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
  transition: background-color 0.3s ease;
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.3);
}

.icon-circle:hover {
  background-color: var(--hover-color);
}

.icon-circle img {
  width: 32px;
  height: 32px;
}

.icon-wrapper span {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}
</style>
