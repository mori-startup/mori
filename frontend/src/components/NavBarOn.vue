<template>
  <div id="navbar">
    <div id="menu-btn" @click="toggleSidebar"></div>
    <div id="nav-titleSearch">
      <div class="smallMoriLogo"></div>
      <router-link to="/main" class="mori" id="nav-title">Mori</router-link>
      <Search />
    </div>
    <ul class="nav-links">
      <li id="notifications-link">
        <Notifications />
      </li>
      <li>
        <router-link
          v-if="typeof user.id !== 'undefined'"
          :to="{ name: 'Profile', params: { id: user.id } }"
        >
          My profile
        </router-link>
      </li>
      <li @click="navigateToUpgradeIA" style="cursor: pointer">Upgrade IA</li>
      <li @click="logout">Log out</li>
    </ul>
    <!-- Sidebar -->
    <Sidebar
      :isActive="isSidebarActive"
      :contactsList="contactsList"
      @select-component="$emit('select-component', $event)"
      @select-contact="$emit('select-contact', $event)"
    />
  </div>
</template>

<script>
import Search from "./Search.vue";
import Notifications from "./Notifications.vue";
import Sidebar from "./Sidebar.vue";

export default {
  name: "NavBarOn",
  components: { Notifications, Search, Sidebar },
  props: {
    contactsList: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      user: {},
      isSidebarActive: true, // Controls sidebar visibility
    };
  },
  created() {
    this.getUserInfo();
  },
  methods: {
    async getUserInfo() {
      const response = await fetch("http://localhost:8081/currentUser", {
        credentials: "include",
      });
      const json = await response.json();
      this.user = json.users[0];
    },
    async logout() {
      await fetch("http://localhost:8081/logout", {
        credentials: "include",
        headers: { Accept: "application/json" },
      });
      this.$store.state.wsConn.close(1000, "user logged out");
      this.$router.push("/");
    },
    toggleSidebar() {
      this.isSidebarActive = !this.isSidebarActive;
    },
    navigateToUpgradeIA() {
      this.$router.push({ name: "UpgradeIA" });
    },
    navigateToChatbot() {
      this.$router.push({ name: "mainpage" }); // Navigate to the Chatbot view
    },
    navigateToMessages() {
      if (this.contactsList.length > 0) {
        const firstContact = this.contactsList[0];
        this.$router.push({
          name: "messages",
          query: {
            name: firstContact.nickname,
            receiverId: firstContact.id,
            type: "PERSON",
          },
        });
      }
    },
  },
};
</script>

<style scoped>
#menu-btn {
  width: 30px;
  height: 30px;
  margin-right: 20px;
  background-image: url("../assets/menu.png");
  background-size: cover;
  cursor: pointer;
  transition: all 0.3s;
}

#menu-btn:hover {
  transform: scale(1.05);
}

#navbar {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 3;
  width: 100%;
  min-width: min-content;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 30px;
  background-color: var(--purple-color);
  color: var(--color-white);
  position: sticky;
}

#navbar a {
  color: var(--color-white);
}

#nav-title {
  user-select: none;
  position: relative;
}

.nav-links li {
  user-select: none;
  font-weight: 300;
  display: inline-block;
  margin-left: 20px;
  cursor: pointer;

  position: relative;
}

#nav-titleSearch {
  display: flex;
  gap: 25px;
  flex-grow: 1;
  align-items: center;
}

#navbar li:not(#notifications-link)::after,
#nav-title::after {
  content: "";
  height: 2.5px;
  width: 0;
  display: block;
  position: absolute;

  transition: all 0.35s ease-out;
}

#navbar li:not(#notifications-link):hover::after,
#nav-title:hover::after {
  width: 100%;
  background-color: var(--hover-background-color);
}

a:link {
  text-decoration: none;
}

a:visited {
  text-decoration: none;
}

a:hover {
  text-decoration: none;
}

a:active {
  text-decoration: none;
}
</style>
