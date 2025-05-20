<template>
  <NavBarOn @toggle-sidebar="toggleSidebar" />
  <Sidebar
    :isActive="isSidebarActive"
    :contactsList="contacts"
    @navigate-to="navigateTo"
  />
  <div id="layout">
    <div class="main-content">
      <ChatbotConversation />
    </div>
  </div>
</template>

<script>
import NavBarOn from "@/components/NavBarOn.vue";
import Sidebar from "@/components/Sidebar.vue";
import ChatbotConversation from "@/components/ChatbotConversation.vue";

export default {
  components: {
    NavBarOn,
    Sidebar,
    ChatbotConversation,
  },
  data() {
    return {
      contacts: [],
      isSidebarActive: false,
    };
  },
  methods: {
    updateLayoutWidth() {
      const sidebar = document.querySelector(".sidebar");
      const layout = document.getElementById("layout");

      if (sidebar?.classList.contains("sidebar--active")) {
        layout.style.width = "80%";
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
  },
  mounted() {
    const sidebar = document.querySelector(".sidebar");

    // Ensure layout width is set initially
    this.updateLayoutWidth();

    // Observe changes to the sidebar class
    this.sidebarObserver = new MutationObserver(() => {
      this.updateLayoutWidth();
    });

    this.sidebarObserver.observe(sidebar, {
      attributes: true, // Watch for changes to attributes (like `class`)
      attributeFilter: ["class"], // Only observe the `class` attribute
    });
  },
  beforeUnmount() {
    // Disconnect the observer to prevent memory leaks
    if (this.sidebarObserver) {
      this.sidebarObserver.disconnect();
    }
  },
};
</script>

<style>
html,
body {
  overflow-y: hidden;
}

#layout {
  display: flex;
  height: 95vh;
  width: 100%;
  position: fixed;
  bottom: 0px;
  right: 0px;
  transition: width 0.3s ease;
}

.main-content {
  height: 100%;
  width: fit-content;
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: var(--page-bg);
}
</style>
