<template>
  <div>
    <!-- Navbar -->
    <NavBarOn @toggle-sidebar="toggleSidebar" />
    <!-- Sidebar -->
    <Sidebar
      :isActive="isSidebarActive"
      :contactsList="contacts"
      @navigate-to="navigateTo"
    />
    <!-- Chatbox Content -->
    <div id="layout">
      <div class="chatbox-view-wrapper">
        <header class="chatbox-view-header">
          <div class="header-left-part">
            <div
              class="receiver-avatar"
              :style="{
                backgroundImage: `url(${user.avatar ? user.avatar : 'default-avatar.png'})`
              }"
            ></div>
            <h1 class="receiver-name">{{ user.name || "Unnamed User" }}</h1>
          </div>
          <!-- Only display follow-status for personal chats -->
          <p class="follow-status" v-if="type === 'PERSON'">{{ user.following }}</p>
        </header>

        <div class="chatbox-view-content" ref="contentDiv">
          <div
            class="message"
            v-for="(message, index) in allMessages"
            :style="msgPosition(message)"
            :key="index"
          >
            <div class="receiver-avatar-name">
              <div
  class="receiver-avatar-chat"
  v-if="displayName(message, index)"
  :style="{
    backgroundImage: `url(${type === 'GROUP'
      ? 'http://localhost:8081/' + (message.sender.avatar || 'default-avatar.png')
      : (user.avatar.startsWith('http') ? user.avatar : 'http://localhost:8081/' + (user.avatar || 'default-avatar.png'))
    })`
  }"
></div>
              <p class="message-author" v-if="displayName(message, index)">
                {{ message.sender.nickname }}
              </p>
            </div>
            <p :class="getClass(message)" class="message-content">
              {{ message.content }}
              <p class="message-timeStamp">{{ formatTime(message.createdAt) }}</p>
            </p>
          </div>
          <p class="seen" v-if="isLastMessageSeen">Seen by {{ user.name }}</p>
        </div>

        <form
          @submit.prevent="sendMessage"
          class="chatbox-view-input"
          autocomplete="off"
          @keyup.enter="sendMessage"
        >
          <input
            type="text"
            placeholder="Type your message here..."
            ref="sendMessageInput"
          />
          <button type="submit"><i class="uil uil-message"></i></button>
          <Emojis
            :input="$refs.sendMessageInput"
            :messagebox="$refs.contentDiv"
          />
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import NavBarOn from "@/components/NavBarOn.vue";
import Sidebar from "@/components/Sidebar.vue";
import Emojis from "../components/Chat/Emojis.vue";

export default {
  props: ["name", "receiverId", "type"],
  components: { NavBarOn, Sidebar, Emojis },
  data() {
    return {
      user: {
        name: "",
        avatar: "default-avatar.png",
        following: "",
      },
      previousMessages: [],
      isSidebarActive: false,
      contacts: [],
    };
  },
  computed: {
    allMessages() {
      const storeMessages = this.$store.getters.getMessages(
        this.receiverId,
        this.type
      );
      const uniqueMessages = storeMessages.filter(
        (msg) => !this.previousMessages.some((prevMsg) => prevMsg.id === msg.id)
      );
      return [...this.previousMessages, ...uniqueMessages];
    },
    isLastMessageSeen() {
      const lastMessage = this.allMessages[this.allMessages.length - 1];
      return (
        lastMessage &&
        lastMessage.senderId === this.myID &&
        lastMessage.isRead
      );
    },
    ...mapState({
      myID: (state) => state.id,
    }),
  },
  watch: {
    allMessages() {
      this.$nextTick(() => {
        this.scrollToBottom();
      });
    },
    receiverId: {
      immediate: true,
      handler(newId) {
        this.fetchUserDetails(newId);
        this.getPreviousMessages();
      },
    },
  },
  methods: {
    formatTime(isoString) {
      const date = new Date(isoString);
      if (isNaN(date.getTime())) {
        return "Now";
      }
      return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    },
    async fetchUserDetails(userId) {
      if (this.type === "GROUP") {
        // For group conversations, use the passed name and default group avatar.
        this.user = {
          name: this.name || "Group",
          avatar: require("@/assets/group.png"),
          following: "",
        };
        return;
      }
      try {
        const response = await fetch("http://localhost:8081/allUsers", {
          credentials: "include",
        });
        const data = await response.json();
        const foundUser = data.users.find((user) => user.id === userId);
        if (foundUser) {
          foundUser.following = foundUser.follower ? "Follows you" : "Not Following you";
          this.user = {
            name: foundUser.nickname,
            following: foundUser.following,
            avatar: `http://localhost:8081/${foundUser.avatar}` || "default-avatar.png",
          };
        } else {
          this.user = { name: "Unknown User", avatar: "default-avatar.png", following: "" };
        }
      } catch (error) {
        console.error("Error fetching user details:", error);
        this.user = { name: "Unknown User", avatar: "default-avatar.png", following: "" };
      }
    },
    updateLayoutWidth() {
      const sidebar = document.querySelector(".sidebar");
      const layout = document.getElementById("layout");
      if (sidebar?.classList.contains("sidebar--active")) {
        layout.style.width = "70%";
      } else {
        layout.style.width = "100%";
      }
    },
    toggleSidebar() {
        this.isSidebarActive = !this.isSidebarActive;
    },
    navigateTo(target) {
      if (target === "chatbot") {
        this.$router.push({ name: "mainpage" });
      } else if (target === "messages") {
        this.$router.push({ name: "messages" });
      }
    },
    async getPreviousMessages() {
      try {
        const response = await fetch("http://localhost:8081/messages", {
          credentials: "include",
          method: "POST",
          body: JSON.stringify({
            type: this.type,
            receiverId: this.receiverId,
          }),
        });
        const data = await response.json();
        const storeMessages = this.$store.getters.getMessages(
          this.receiverId,
          this.type
        );
        this.previousMessages = (data.chatMessage || [])
          .filter((msg) => !storeMessages.some((storeMsg) => storeMsg.id === msg.id))
          .map((msg) => ({
            ...msg,
            createdAt: this.removeZFromTimestamp(msg.createdAt),
            isRead: msg.isRead || false,
          }));
        this.scrollToBottom();
      } catch (error) {
        console.error("Error fetching messages:", error);
      }
    },
    async sendMessage() {
      const sendMessageInput = this.$refs.sendMessageInput;
      if (sendMessageInput.value === "") return;
      const msgObj = {
        receiverId: this.receiverId,
        content: sendMessageInput.value,
        createdAt: new Date(),
        type: this.type,
      };
      try {
        const response = await fetch("http://localhost:8081/newMessage", {
          body: JSON.stringify(msgObj),
          method: "POST",
          credentials: "include",
        });
        const data = await response.json();
        if (data.type === "Success") {
          this.$store.dispatch("addNewChatMessage", {
            ...msgObj,
            senderId: this.myID,
          });
          this.scrollToBottom();
        } else {
          this.$toast.open({
            message: data.message,
            type: "warning",
          });
        }
        sendMessageInput.value = "";
      } catch (error) {
        console.error("Error sending message:", error);
      }
    },
    clearChatNewMessages() {
      if (this.type === "GROUP") {
        const msgs = this.$store.state.chat.newGroupChatMessages.filter(
          (msg) => msg.receiverId !== this.receiverId
        );
        this.$store.commit("updateNewGroupChatMessages", msgs);
      } else {
        const msgs = this.$store.state.chat.newChatMessages.filter(
          (msg) =>
            msg.receiverId !== this.receiverId &&
            msg.senderId !== this.receiverId
        );
        this.$store.commit("updateNewChatMessages", msgs);
      }
    },
    displayName(message, index) {
      if (message.senderId === this.myID) return false;
      if (index < 1) return true;
      return message.senderId !== this.allMessages[index - 1]?.senderId;
    },
    getClass(message) {
      return message.senderId === this.myID
        ? { "sent-message": true }
        : { "received-message": true };
    },
    msgPosition(message) {
      return {
        alignSelf: message.senderId === this.myID ? "flex-end" : "flex-start",
      };
    },
    markMessageAsSeen(messageID) {
      console.log(`Marking message ${messageID} as seen...`);
      fetch("http://localhost:8081/messageRead", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          messageID: messageID,
          type: this.type,
        }),
      })
        .then((response) => response.json())
        .then((data) => {
          if (data.type === "Success") {
            console.log(`Message ${messageID} marked as read successfully.`);
            this.$store.dispatch("markMessageAsSeen", { messageID });
          } else {
            console.error("Failed to mark message as read:", data.message);
          }
        })
        .catch((error) => {
          console.error("Error marking message as seen:", error);
        });
    },
    scrollToBottom() {
      this.$nextTick(() => {
        if (this.$refs.contentDiv) {
          this.$refs.contentDiv.scrollTop = this.$refs.contentDiv.scrollHeight;
          const lastMessage = this.allMessages[this.allMessages.length - 1];
          if (lastMessage && lastMessage.senderId !== this.myID && !lastMessage.isRead) {
            this.markMessageAsSeen(lastMessage.id);
          }
        }
      });
    },
    removeZFromTimestamp(timestamp) {
      if (!timestamp || typeof timestamp !== "string") {
        return "Invalid timestamp";
      }
      return timestamp.replace("Z", "");
    },
  },
  async mounted() {
    this.fetchUserDetails(this.receiverId);
    const sidebar = document.querySelector(".sidebar");
    this.updateLayoutWidth();
    this.sidebarObserver = new MutationObserver(() => {
      this.updateLayoutWidth();
    });
    this.sidebarObserver.observe(sidebar, {
      attributes: true,
      attributeFilter: ["class"],
    });
  },
  beforeUnmount() {
    if (this.sidebarObserver) {
      this.sidebarObserver.disconnect();
    }
  },
  created() {
    console.log("Component created: chatbox");
    this.getPreviousMessages();
  },
  unmounted() {
    console.log("Component unmounted: chatbox");
    this.clearChatNewMessages();
  },
};
</script>

<style scoped>
#layout {
  display: flex;
  height: 95vh;
  width: 100%;
  position: fixed;
  bottom: 0;
  align-items: center;
  justify-content: center;
  right: 0;
  transition: all 0.3s ease;
}
.chatbox-view-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--page-bg);
  width: 60%;
  outline: 1px solid var(--bg-neutral);
}
.chatbox-view-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  background-color: var(--bg-neutral);
  padding: 15px;
  border-radius: 0 0 15px 15px;
  box-shadow: 0 2px 10px rgb(0, 0, 0);
  z-index: 1;
  font-size: 1.5em;
}
.header-left-part {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: left;
  color: var(--color-white);
  gap: 15px;
  font-size: 1.5em;
}
.receiver-name {
  align-items: center;
  justify-content: center;
  text-align: center;
  margin-top: 15px;
  font-size: 23px;
}
.receiver-avatar {
  margin-top: 15px;
  margin-left: 10px;
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  background-color: var(--purple-color);
  border: 2px solid var(--color-white);
}
.receiver-avatar-name {
  display: flex;
  align-items: end;
  gap: 6px;
}
.receiver-avatar-chat {
  box-shadow: 0 2px 10px rgb(0, 0, 0);
  background-color: var(--purple-color);
  border: 1px solid var(--color-white);
  border-radius: 50%;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  width: 40px;
  height: 40px;
  margin-bottom: 15px;
}
.follow-status {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 18px;
  margin-right: 30px;
  color: var(--purple-color);
  font-size: 16px;
}
.chatbox-view-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  border-radius: 15px;
  gap: 10px;
}
.message-author {
  font-size: 0.9em;
  color: var(--purple-color);
  margin-bottom: 15px;
}
.message-content {
  padding: 10px;
  border-radius: 10px;
  word-break: break-word;
}
.message-timeStamp {
  font-size: 0.7em;
  color: var(--color-grey);
  opacity: 0.5;
  text-align: right;
  margin-top: 5px;
}
.message {
  max-width: 80%;
}
.sent-message {
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.3);
  background-color: var(--purple-color);
  color: var(--color-white);
}
.received-message {
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.3);
  background-color: var(--bg-neutral);
  color: var(--color-white);
}
.chatbox-view-input {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  border-radius: 15px 15px 0 0;
  background-color: var(--bg-neutral);
  box-shadow: 0 2px 10px rgb(0, 0, 0);
  transition: all 0.3s ease;
}
.chatbox-view-input:hover {
  transform: scale(1.02);
}
.chatbox-view-input input {
  flex: 1;
  padding: 20px;
  border-radius: 15px;
  height: 45px;
  border: 1px solid var(--color-grey-light);
}
.chatbox-view-input button {
  background-color: var(--purple-color);
  color: var(--color-white);
  border: none;
  padding: 10px 15px;
  border-radius: 15px;
  font-size: 1.2em;
  cursor: pointer;
  transition: all 0.3s ease;
}
.chatbox-view-input button:hover {
  background-color: var(--hover-color);
}
.seen {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  color: grey;
}
</style>
