<template>
  <div>
    <!-- NavBar with sidebar toggle support -->
    <NavBarOn @toggle-sidebar="toggleSidebar" />
    <!-- Sidebar -->
    <Sidebar
      :isActive="isSidebarActive"
      :contactsList="contacts"
      @select-contact="handleContactSelection"
      @navigate-to="navigateTo"
    />
    <!-- Layout container that adjusts its width based on sidebar state -->
    <div id="layout">
      <div class="dropfiles-view-wrapper">
        <UpgradeIA />
      </div>
    </div>
  </div>
</template>

<script>
import NavBarOn from "@/components/NavBarOn.vue";
import Sidebar from "@/components/Sidebar.vue";
import UpgradeIA from "@/components/UpgradeIA.vue";

export default {
  name: "DropFiles",
  components: {
    NavBarOn,
    Sidebar,
    UpgradeIA,
  },
  data() {
    return {
      isSidebarActive: false, // Set sidebar open by default for this view.
      contacts: [], // Load contacts as needed (or from Vuex)
    };
  },
  methods: {
    toggleSidebar(value) {
      // If a boolean value is passed, set isSidebarActive to that value.
      // Otherwise, toggle the state.
      if (typeof value === "boolean") {
        this.isSidebarActive = value;
      } else {
        this.isSidebarActive = !this.isSidebarActive;
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
    navigateTo(target) {
      // Navigate based on target.
      if (target === "chatbot") {
        this.$router.push({ name: "mainpage" });
      } else if (target === "messages") {
        this.$router.push({ name: "messages" });
      }
    },
    updateLayoutWidth() {
      const sidebar = document.querySelector(".sidebar");
      const layout = document.getElementById("layout");
      if (sidebar && sidebar.classList.contains("sidebar--active")) {
        layout.style.width = "80%";
      } else {
        layout.style.width = "100%";
      }
    },
  },
  mounted() {
    // Set the initial layout width.
    this.updateLayoutWidth();

    // Observe sidebar class changes to update layout width dynamically.
    const sidebar = document.querySelector(".sidebar");
    this.sidebarObserver = new MutationObserver(() => {
      this.updateLayoutWidth();
    });
    if (sidebar) {
      this.sidebarObserver.observe(sidebar, {
        attributes: true,
        attributeFilter: ["class"],
      });
    }
  },
  beforeUnmount() {
    if (this.sidebarObserver) {
      this.sidebarObserver.disconnect();
    }
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
  right: 0;
  transition: width 0.3s ease;
}

.dropfiles-view-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-top: 50px;
  justify-content: flex-start;
  align-items: center;
  background-color: var(--page-bg);
}
</style>
