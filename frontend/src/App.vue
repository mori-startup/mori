<template>
  <router-view @hide-chat="toggleChatVisibility"></router-view>
  <Chat v-if="showSmallChat" />
</template>

<script>
import Chat from "./components/Chat/Chat.vue";
import { mapState } from "vuex";

export default {
  name: "App",
  components: { Chat },
  data() {
    return {
      showSmallChat: true, // Controls visibility of the small chatbox
    };
  },
  computed: {
    ...mapState({
      currentRoute: (state) => state.route, // Map current route from Vuex or Vue Router
    }),
    shouldHideSmallChat() {
      // Hide small chatbox on specific routes
      const hiddenRoutes = [
        "/sign-in",
        "/register",
        "/messages",
        "/verified",
        "/forgotpassword",
        "/reset-password",
        "/upgrade-ia"
      ];
      return hiddenRoutes.includes(this.$route.path);
    },
  },
  watch: {
    $route() {
      this.toggleChatVisibility();
    },
  },
  methods: {
    toggleChatVisibility() {
      this.showSmallChat = !this.shouldHideSmallChat;
    },
    createWebSocketConn() {
      const excludedPaths = ["/sign-in", "/register", "/verified", "/forgotpassword", "/reset-password"];
      if (excludedPaths.includes(this.$route.path)) {
        return;
      }
      this.$store.dispatch("createWebSocketConn");
    },
  },
  created() {
    this.toggleChatVisibility();
    this.createWebSocketConn();
  },
};
</script>
