<template>
  <div class="contacts-wrapper">
    <h1 class="mori">Mori <span class="adder">- Messaging</span></h1>

    <!-- Section Friends -->
    <h3 class="sous_titres">
      Friends (<span class="purple-strong">{{ chatUserList.length }}</span>)
    </h3>
    <ul class="horizontal-list">
      <li
        v-for="contact in chatUserList"
        :key="contact.id"
        @click="selectContact(contact, 'PERSON')"
        class="contact-item-horizontal"
      >
        <div
          class="user-picture small"
          :style="{
            backgroundImage: `url(http://localhost:8081/${contact.avatar})`,
          }"
        ></div>
        <div class="contact-name">{{ contact.nickname }}</div>
      </li>
    </ul>

    <!-- Section Personal Conversations -->
    <h3 class="sous_titres">
      Conversations (<span class="purple-strong">{{ friends.length }}</span>)
    </h3>
    <div class="conversation-card-wrapper">
      <div
        v-for="convMsg in friends"
        :key="convMsg.id"
        class="conversation-card"
        @click="selectContact(convMsg, 'PERSON')"
      >
        <div
          class="avatar"
          :style="{
            backgroundImage: `url(http://localhost:8081/${convMsg.avatar})`,
          }"
        ></div>
        <div class="content">
          <div class="header">
            <span class="name">{{ convMsg.name }}</span>
            <span class="time">{{ formatTime(convMsg.lastMessageTime) }}</span>
          </div>
          <div class="message-preview">
            {{
              convMsg.lastMessageSenderId === myId
                ? "You: "
                : convMsg.name + ": "
            }}
            {{
              convMsg.lastMessage && convMsg.lastMessage.length > 40
                ? convMsg.lastMessage.slice(0, 40) + "…"
                : convMsg.lastMessage
            }}
            <span
              v-if="unreadCount(convMsg.id) > 0"
              class="unread-badge"
            >
              {{ unreadCount(convMsg.id) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Section Groups -->
    <h3 class="sous_titres">
      Your groups (<span class="purple-strong">{{ userGroups.length }}</span>)
    </h3>
    <NewGroup />
    <ul class="horizontal-list">
      <li
        v-for="group in userGroups"
        :key="group.id"
        @click="selectContact(group, 'GROUP')"
        class="contact-item-horizontal"
      >
        <div
          class="user-picture small"
          :style="{ backgroundImage: `url(${defaultGroupLogo})` }"
        ></div>
        <div class="contact-name">{{ group.name }}</div>
      </li>
    </ul>

    <!-- Section Group Conversations -->
    <h3 class="sous_titres">
      Group conversations (<span class="purple-strong">{{ groups.length }}</span>)
    </h3>
    <div class="conversation-card-wrapper">
      <div
        v-for="convMsg in groups"
        :key="convMsg.id"
        class="conversation-card"
        @click="selectContact(convMsg, 'GROUP')"
      >
        <div
          class="avatar"
          :style="{ backgroundImage: `url(${defaultGroupLogo})` }"
        ></div>
        <div class="content">
          <div class="header">
            <span class="name">{{ convMsg.name }}</span>
            <span class="time">{{ formatTime(convMsg.lastMessageTime) }}</span>
          </div>
          <div class="message-preview">
            {{
              convMsg.lastMessageSenderId === myId
                ? "You: "
                : convMsg.name + ": "
            }}
            {{
              convMsg.lastMessage && convMsg.lastMessage.length > 40
                ? convMsg.lastMessage.slice(0, 40) + "…"
                : convMsg.lastMessage
            }}
            <span
              v-if="unreadGroupCount(convMsg.id) > 0"
              class="unread-badge"
            >
              {{ unreadGroupCount(convMsg.id) }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import NewGroup from "@/components/NewGroup.vue";
import defaultGroupLogo from "@/assets/group.png";

export default {
  name: "ContactsForChatBotView",
  components: {
    NewGroup,
  },
  data() {
    return {
      // Tracks the conversation currently open
      currentChatId: null,
    };
  },
  computed: {
    ...mapState({
      chatUserList: (state) => state.chat.chatUserList || [],
      userGroups: (state) => state.groups.userGroups || [],
      conversationsMsg: (state) => state.conversationsMsg,
      newChatMessages: (state) => state.chat.newChatMessages,
      unreadMsgsStatsFromDB: (state) => state.chat.unreadMsgsStatsFromDB,
      openChats: (state) => state.chat.openChats,
      newGroupChatMessages: (state) => state.chat.newGroupChatMessages,
      myId: (state) => state.id, // Adjust according to your store
    }),
    ...mapGetters([
      "getUnreadMessagesCount",
      "getUnreadGroupMessagesCount",
      "getUnreadMsgsCountFromDB",
      "getMessages",
    ]),
    // Update personal conversations with the latest message dynamically.
    friends() {
      if (!Array.isArray(this.conversationsMsg)) return [];
      return this.conversationsMsg
        .filter((c) => c.type === "PERSON")
        .map((conv) => {
          const msgs = this.getMessages(conv.id, "PERSON");
          if (msgs && msgs.length > 0) {
            const lastMsg = msgs[msgs.length - 1];
            conv.lastMessage = lastMsg.content;
            conv.lastMessageTime = lastMsg.time;
            conv.lastMessageSenderId = lastMsg.senderId;
          }
          return conv;
        });
    },
    // Update group conversations similarly.
    groups() {
      if (!Array.isArray(this.conversationsMsg)) return [];
      return this.conversationsMsg
        .filter((c) => c.type === "GROUP")
        .map((conv) => {
          const msgs = this.getMessages(conv.id, "GROUP");
          if (msgs && msgs.length > 0) {
            const lastMsg = msgs[msgs.length - 1];
            conv.lastMessage = lastMsg.content;
            conv.lastMessageTime = lastMsg.time;
            conv.lastMessageSenderId = lastMsg.senderId;
          }
          return conv;
        });
    },
    // Make the default group logo available to the template.
    defaultGroupLogo() {
      return defaultGroupLogo;
    },
  },
  created() {
    this.$store.dispatch("fetchConversationsMsg");
  },
  methods: {
    selectContact(contact, type) {
      // Set the active conversation ID so that unread count for this chat becomes 0
      this.currentChatId = contact.id;
      this.$emit("select-contact", {
        id: contact.id,
        name: contact.nickname || contact.name,
        type,
      });
      // Clear unread messages for this conversation from the store
      this.$store.dispatch("removeUnreadMessages", { receiverId: contact.id, type });
    },
    formatTime(isoString) {
      const date = new Date(isoString);
      // Return a default value if the date is invalid.
      if (isNaN(date.getTime())) {
        return "Now";
      }
      return date.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });
    },
    // Helper method for personal conversations.
    unreadCount(conversationId) {
      if (this.currentChatId === conversationId) {
        return 0;
      }
      return this.getUnreadMessagesCount(conversationId);
    },
    // Helper method for group conversations.
    unreadGroupCount(conversationId) {
      if (this.currentChatId === conversationId) {
        return 0;
      }
      return this.getUnreadGroupMessagesCount(conversationId);
    },
  },
};
</script>

<style scoped>
.purple-strong {
  color: var(--purple-color);
}
.contacts-wrapper {
  margin-top: -5px;
  background-color: var(--bg-neutral);
}

.titre {
  color: white;
  font-size: 1.8em;
}

.sous_titres {
  color: white;
  font-size: 1.2em;
  margin-top: 2vh;
  margin-bottom: 10px;
}

/* Horizontal list */
.horizontal-list {
  display: flex;
  flex-wrap: nowrap;
  overflow-x: auto;
  gap: 10px;
  list-style: none;
  padding-top: 10px;
  padding-bottom: 10px;
  padding-left: 5px;
  padding-right: 5px;
  margin: 10px 0;
}

.contact-item-horizontal {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  cursor: pointer;
  width: 60px;
}

.user-picture.small {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background-size: cover;
  background-position: center;
  margin-bottom: 5px;
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.6);
  transition: all 0.3s ease;
}

.user-picture.small:hover {
  transform: scale(1.1);
  border: 1px solid var(--purple-color);
}

.contact-name {
  font-size: 12px;
  color: var(--text-primary);
}

/* Conversation cards */
.conversation-card-wrapper {
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-height: 150px;
  overflow-y: scroll;
  padding-right: 10px;
}

.conversation-card {
  display: flex;
  align-items: center;
  background-color: #3a3a3a;
  border-radius: 10px;
  padding: 10px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.conversation-card:hover {
  background-color: #4a4a4a;
}

.avatar {
  width: 50px;
  height: 50px;
  background-color: var(--purple-color);
  border: 1px solid var(--color-white);
  border-radius: 50%;
  background-size: cover;
  background-repeat: no-repeat;
  background-position: center;
  margin-right: 15px;
}

.content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
}

.name {
  font-size: 14px;
  font-weight: bold;
  color: white;
}

.time {
  font-size: 12px;
  color: #ccc;
}

/* Unread badge */
.unread-badge {
  padding: 2px 7px;
  color: var(--color-white);
  background-color: brown;
  font-size: 12px;
  border-radius: 5px;
  align-self: flex-end;
}

/* Message preview */
.message-preview {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: #ccc;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 300px;
}

/* Filter Options */
.filter-options {
  margin-bottom: 10px;
  display: flex;
  gap: 20px;
  align-items: center;
}
.filter-options .search-bar {
  flex: 1;
}
.search-bar input {
  width: 100%;
  padding: 8px;
  border-radius: 10px;
  border: 1px solid #ccc;
}
.extra-filters {
  display: flex;
  gap: 15px;
  margin-top: 5px;
}
.extra-filters label {
  color: var(--color-white);
}

.user-picture.small {
  background-color: var(--purple-color);
  border: 1px solid var(--color-white);
  background-repeat: no-repeat;
  background-size: cover;
  background-position: center;
}
</style>
