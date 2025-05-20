<template>
    <div>
      <!-- Navbar with sidebar toggle support -->
      <NavBarOn @toggle-sidebar="toggleSidebar" />
      <!-- Sidebar (shows contacts, groups, etc.) -->
      <Sidebar
        :isActive="isSidebarActive"
        :contactsList="contacts"
        @select-contact="handleContactSelection"
      />
      <!-- Layout container that adjusts its width based on the sidebar -->
      <div id="layout">
        <div class="profile-view-wrapper">
          <Profile />
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import NavBarOn from "@/components/NavBarOn.vue";
  import Sidebar from "@/components/Sidebar.vue";
  import Profile from "@/components/Profile.vue";
  
  export default {
    name: "ProfileView",
    components: { NavBarOn, Sidebar, Profile },
    data() {
      return {
        isSidebarActive: false, // Controls sidebar visibility
        contacts: [], // If needed, you can load your contacts here or from Vuex
      };
    },
    methods: {
      toggleSidebar(value) {
        // If a value is passed, set the sidebar state accordingly.
        // Otherwise, simply toggle.
        if (typeof value === "boolean") {
          this.isSidebarActive = value;
        } else {
          this.isSidebarActive = !this.isSidebarActive;
        }
      },
      handleContactSelection({ id, name, type }) {
        // Navigate to the selected conversation (if desired)
        this.$router.push({
          name: "messages",
          query: { name, receiverId: id, type },
        });
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
      // Set initial layout width
      this.updateLayoutWidth();
  
      // Observe sidebar changes to update layout width dynamically
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
  
  .profile-view-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding-bottom: 30px;
    justify-content: flex-start;
    align-items: center;
    background-color: var(--page-bg);
    overflow-y: scroll;
  }
  </style>
  