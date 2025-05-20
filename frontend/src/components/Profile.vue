<template>
  <div v-if="user && $store.state.id !== ''" class="user-profil">
    <div id="layout-profile">
      <!-- Left Section -->
      <div class="left-section">
        <div class="user-profile__public">
          <div
            class="user-picture"
            :style="{ backgroundImage: `url(http://localhost:8081/${user.avatar})` }"
          ></div>
          <div class="user-profile__info">
            <h2 class="username">{{ user.firstName }} {{ user.lastName }}</h2>
            <h3 v-if="showNickname" class="username">{{ user.nickname }}</h3>
            <p class="user-email" v-if="user.login">{{ user.login }}</p>
            <p class="user-dateOfBirth" v-if="user.dateOfBirth">{{ user.dateOfBirth }}</p>
          </div>

          <div class="profile-btns">
            <PrivacyBtn v-if="isMyProfile" :status="user.status" />
            <!-- Follow/unfollow button -->
            <component
              v-else
              :is="displayBtn"
              v-bind="{ user }"
              @follow="checkFollowRequest"
              @unfollow="unfollow"
            ></component>
          </div>
        </div>

        <div class="multiple-item-list" v-if="showProfileData">
          <Following :following="following" />
          <Followers :followers="followers" />
        </div>

        <Groups :groups="profileGroups" v-if="showProfileData" />
      </div>

      <!-- Middle Section -->
      <div class="middle-section" v-if="showProfileData">
        <div class="about" v-if="user.about !== ''">
          <h2 class="about-title">About me</h2>
          <p class="about-text">{{ user.about }}</p>
        </div>

        <!-- Edit Profile Banner (only for own profile) -->
        <div v-if="isMyProfile" class="edit-profile-banner">
          <button class="toggle-edit-banner btn-primary" @click="toggleEditBanner">
            {{ editBannerOpen ? "Close Edit Profile" : "Edit Profile" }}
          </button>
          <transition name="slide">
            <div v-if="editBannerOpen" class="edit-profile">
              <h2 class="edit-profile-title">Modify profile</h2>

              <!-- Change Nickname -->
              <div class="edit-section">
                <h3 class="edit-section-title">New Nickname</h3>
                <input
                  v-model="newNickname"
                  type="text"
                  placeholder="Entrez votre nouveau pseudo"
                  class="input-field"
                />
                <button @click="changeNickname" class="btn-primary">
                  Modify Nickname
                </button>
              </div>

              <!-- Change Avatar -->
              <div class="edit-section">
                <h3 class="edit-section-title">New Avatar</h3>
                <input
                  type="file"
                  ref="avatarInput"
                  class="file-input"
                  accept="image/*"
                />
                <button @click="changeAvatar" class="btn-primary">
                  Modify Avatar
                </button>
              </div>

              <!-- Delete Account Button -->
              <div class="edit-section">
                <button @click="openDeleteModal" class="btn-delete">
                  Delete Account
                </button>
              </div>
            </div>
          </transition>
        </div>
      </div>

      <!-- If profile is private -->
      <p v-else class="additional-info large">
        This profile is private
      </p>
    </div>

    <!-- Delete Account Modal -->
    <div v-if="showDeleteModal" class="modal-overlay">
      <div class="modal-content">
        <h3 id="deleteText">Are you sure you want to delete your account?</h3>
        <div class="modal-buttons">
          <button @click="confirmDelete" class="btn-confirm">Yes</button>
          <button @click="cancelDelete" class="btn-cancel">No</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Following from "./Following.vue";
import Followers from "./Followers.vue";
import FollowBtn from "./FollowBtn.vue";
import PrivacyBtn from "./PrivacyBtn.vue";
import UnfollowBtn from "./UnfollowBtn.vue";
import Groups from "./Groups.vue";

export default {
  name: "Profile",
  components: { Following, Followers, FollowBtn, PrivacyBtn, UnfollowBtn, Groups },
  data() {
    return {
      user: null,
      isMyProfile: false,
      followers: [],
      following: [],
      profileGroups: null,
      newNickname: "", // New nickname
      selectedAvatar: null, // File for avatar
      editBannerOpen: false, // Controls visibility of the edit-profile banner
      showDeleteModal: false, // Controls visibility of the delete account modal
    };
  },
  created() {
    this.updateProfileData();
  },
  computed: {
    showProfileData() {
      return this.user && (this.user.following || this.isMyProfile || this.user.status === "PUBLIC");
    },
    showNickname() {
      return this.user && this.user.nickname !== `${this.user.firstName} ${this.user.lastName}`;
    },
    displayBtn() {
      return this.user && this.user.following ? UnfollowBtn : FollowBtn;
    },
  },
  methods: {
    updateProfileData() {
      this.getUserData();
      this.getFollowers();
      this.getFollowing();
      this.checkProfile();
      this.getProfileGroups();
    },
    async getUserData() {
      await fetch("http://localhost:8081/userData?userId=" + this.$route.params.id, {
        credentials: "include",
      })
        .then((r) => r.json())
        .then((json) => {
          console.log("/getUserData", json);
          console.log("id", this.$route.params.id);
          this.user = json.users[0];
        });
    },
    async getProfileGroups() {
      const response = await fetch("http://localhost:8081/otherUserGroups?userId=" + this.$route.params.id, {
        credentials: "include",
      });
      const data = await response.json();
      if (data.type == "Error") {
        console.log("/getProfileGroups error: ", data.message);
      } else {
        this.profileGroups = data.groups;
      }
    },
    async getMyUserID() {
      await this.$store.dispatch("getMyUserID");
    },
    async checkProfile() {
      await this.getMyUserID();
      const myID = this.$store.state.id;
      const profileID = this.$route.params.id;
      this.isMyProfile = (profileID === myID);
    },
    checkFollowRequest(action) {
      if (action === "followedUser") {
        this.$store.dispatch("fetchChatUserList");
        this.updateProfileData();
        this.toggleFollowingThisUser();
      }
    },
    toggleFollowingThisUser() {
      this.user.following = !this.user.following;
    },
    unfollow() {
      this.updateProfileData();
      this.$store.dispatch("fetchChatUserList");
    },
    async getFollowers() {
      await fetch("http://localhost:8081/followers?userId=" + this.$route.params.id, {
        credentials: "include",
      })
        .then((response) => response.json())
        .then((json) => {
          this.followers = json.users;
        });
    },
    async getFollowing() {
      await fetch("http://localhost:8081/following?userId=" + this.$route.params.id, {
        credentials: "include",
      })
        .then((response) => response.json())
        .then((json) => {
          this.following = json.users;
        });
    },
    async changeNickname() {
      await this.$store.dispatch("changeNickname", this.newNickname);
      this.updateProfileData();
    },
    async changeAvatar() {
      const input = this.$refs.avatarInput;
      if (!input || !input.files || input.files.length === 0) {
        this.$toast.open({
          message: "Please select a file.",
          type: "warning",
        });
        return;
      }
      this.selectedAvatar = input.files[0];
      await this.$store.dispatch("changeAvatar", this.selectedAvatar);
      this.updateProfileData();
    },
    addChat() {
      let newChat = {
        name: this.user.nickname,
        receiverId: this.user.id,
        type: "PERSON",
      };
      this.$store.dispatch("addNewChat", newChat);
    },
    toggleEditBanner() {
      this.editBannerOpen = !this.editBannerOpen;
    },
    // Delete Account logic
    openDeleteModal() {
      this.showDeleteModal = true;
    },
    cancelDelete() {
      this.showDeleteModal = false;
    },
    async confirmDelete() {
      try {
        await this.$store.dispatch("deleteAccount");
        this.$toast.open({
          message: "Account deleted successfully.",
          type: "success",
        });
      } catch (err) {
        console.error("Error deleting account:", err);
        this.$toast.open({
          message: err.message,
          type: "error",
        });
      }
    },
  },
  watch: {
    $route() {
      if (this.$route.name === "Profile") {
        this.updateProfileData();
      }
    },
  },
};
</script>

<style scoped>
h3 {
  color: white;
}

#deleteText{
  color: rgb(17, 17, 17);
}
.user-profile {
  overflow: scroll;
}

#layout-profile {
  display: grid;
  grid-template-columns: 1fr minmax(min-content, 550px) 1fr;
  column-gap: 50px;
  margin-top: 50px;
}

.middle-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 50px;
}

.left-section {
  display: flex;
  flex-direction: column;
  gap: 35px;
  max-width: 250px;
  justify-self: flex-end;
}

.user-profile__public,
.user-profile__private {
  display: flex;
  flex-direction: column;
  padding: var(--container-padding);
  background-color: var(--bg-neutral);
  box-shadow: 0 2px 10px rgb(0, 0, 0);
  border-radius: var(--container-border-radius);
  align-items: center;
  text-align: center;
  gap: 25px;
}

.user-profile__public p,
.user-profile__private p,
.user-profile__public h3,
.user-profile__private h3 {
  color: var(--color-white);
}

.user-profile__public h2,
.user-profile__private h2 {
  color: var(--purple-color);
}

.user-profile__privacy {
  color: var(--color-white);
}

:is(.user-profile__public, .user-profile__private) .user-picture {
  height: 185px;
  width: 185px;
}

.user-profile__info {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.profile-btns {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.additional-info {
  text-align: center;
}

/* Edit Profile Banner Styles */
.edit-profile-banner {
  width: 100%;
  margin-top: 20px;
}
.toggle-edit-banner {
  margin-bottom: 10px;
  padding: 10px 20px;
  background-color: var(--purple-color);
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}
.slide-enter-active,
.slide-leave-active {
  transition: max-height 0.3s ease, opacity 0.3s ease;
}
.slide-enter,
.slide-leave-to {
  max-height: 0;
  opacity: 0;
  overflow: hidden;
}

/* Edit Profile Section */
.edit-profile {
  padding: 20px;
  background-color: var(--bg-neutral);
  border-radius: var(--container-border-radius);
  box-shadow: 0 2px 10px rgb(0, 0, 0);
  width: 550px;
}

.edit-section {
  margin-bottom: 20px;
}

.input-field {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  margin-bottom: 10px;
}

.file-input {
  margin-bottom: 10px;
  border-radius: 5px;
}

.btn-primary {
  padding: 10px 20px;
  background-color: var(--purple-color);
  color: #fff;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-primary:hover {
  background-color: var(--hover-background-color);
}

.btn-delete {
  padding: 10px 20px;
  background-color: #e74c3c;
  color: #fff;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-delete:hover {
  background-color: #c0392b;
}

.edit-profile-title {
  text-align: center;
  color: var(--purple-color);
  margin-bottom: 20px;
}

.edit-section-title {
  margin-bottom: 15px;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: #fff;
  padding: 30px;
  border-radius: 10px;
  text-align: center;
  width: 90%;
  max-width: 400px;
}

.modal-buttons {
  margin-top: 20px;
  display: flex;
  justify-content: space-around;
}

.btn-confirm {
  padding: 10px 20px;
  background-color: #e74c3c;
  color: #fff;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.btn-cancel {
  padding: 10px 20px;
  background-color: #95a5a6;
  color: #fff;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}
</style>
