import { useToast } from "vue-toast-notification";
const toast = useToast();

export default {
  async getMyUserID({ commit }) {
    await fetch("http://localhost:8081/currentUser", {
      credentials: "include",
    })
      .then((r) => r.json())
      .then((json) => {
        commit("updateMyUserID", json.users[0].id);
      });
  },

  async getMyProfileInfo({ dispatch, state, commit }) {
    await dispatch("getMyUserID");
    const userID = state.id;
    await fetch("http://localhost:8081/userData?userId=" + userID, {
      credentials: "include",
    })
      .then((r) => r.json())
      .then((json) => {
        const userInfo = json.users[0];
        commit("updateProfileInfo", userInfo);
      })
      .catch((err) => {
        toast.open({
          message: "Error fetching profile info: " + err.message,
          type: "error",
        });
      });
  },

  async getAllUsers({ commit }) {
    await fetch("http://localhost:8081/allUsers", {
      credentials: "include",
    })
      .then((r) => r.json())
      .then((json) => {
        const users = json.users;
        commit("updateAllUsers", users);
      })
      .catch((err) => {
        toast.open({
          message: "Error fetching all users: " + err.message,
          type: "error",
        });
      });
  },

  async getAllGroups({ commit }) {
    await fetch("http://localhost:8081/allGroups", {
      credentials: "include",
    })
      .then((r) => r.json())
      .then((json) => {
        const groups = json.groups;
        commit("updateAllGroups", groups);
      })
      .catch((err) => {
        toast.open({
          message: "Error fetching all groups: " + err.message,
          type: "error",
        });
      });
  },

  async getUserGroups({ commit }) {
    try {
      const response = await fetch(`http://localhost:8081/userGroups`, {
        credentials: "include",
      });
      const data = await response.json();
      commit("updateUserGroups", data.groups);
      commit("updateDataLoaded", "userGroups");
    } catch (err) {
      toast.open({
        message: "Error fetching user groups: " + err.message,
        type: "error",
      });
    }
  },

  addUserGroup({ state, commit }, userGroup) {
    let userGroups = state.groups.userGroups;
    console.log("userGroups state", userGroups);
    if (userGroups === null) {
      userGroups = [];
    }
    userGroups.push(userGroup);
    console.log("userGroup", userGroup);
    commit("updateUserGroups", userGroups);
  },

  async getMyFollowers({ dispatch, commit, state }) {
    await dispatch("getMyProfileInfo");
    const myID = state.profileInfo.id;
    try {
      const response = await fetch(
        `http://localhost:8081/followers?userId=${myID}`,
        {
          credentials: "include",
        }
      );
      const data = await response.json();
      commit("updateMyFollowers", data.users);
    } catch (err) {
      toast.open({
        message: "Error fetching followers: " + err.message,
        type: "error",
      });
    }
  },

  async isLoggedIn() {
    const response = await fetch("http://localhost:8081/sessionActive", {
      credentials: "include",
    });
    const data = await response.json();
    return data.message === "Session active";
  },

  markMessageAsSeen({ commit, state }, { messageID }) {
    if (!Array.isArray(state.newChatMessages)) {
      console.error("newChatMessages is not an array or undefined.");
      return;
    }
    const updatedMessages = state.newChatMessages.map((msg) =>
      msg.id === messageID ? { ...msg, isRead: true } : msg
    );
    commit("updateNewChatMessages", updatedMessages);
  },

  async fetchConversationsMsg({ commit }) {
    const resp = await fetch("http://localhost:8081/conversationsMsg", {
      credentials: "include",
    });
    const data = await resp.json();
    commit("setConversationsMsg", data.conversationsMsg);
  },

  async changeNickname({ dispatch }, newNickname) {
    if (!newNickname.trim()) {
      toast.open({
        message: "Please select a valid Nickname.",
        type: "warning",
      });
      return;
    }
    try {
      const response = await fetch("http://localhost:8081/updateNickname", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ nickname: newNickname }),
      });
      const result = await response.json();
      if (!response.ok) {
        throw new Error(
          result.message || "Erreur lors de la modification du pseudo"
        );
      }
      toast.open({
        message: "Pseudo modifié avec succès !",
        type: "success",
      });
      await dispatch("getMyProfileInfo");
    } catch (err) {
      console.error("Erreur lors de la modification du pseudo :", err.message);
      toast.open({
        message: err.message,
        type: "error",
      });
    }
  },

  async changeAvatar({ dispatch }, avatarFile) {
    if (!avatarFile) {
      toast.open({
        message: "Please select a valid avatar.",
        type: "warning",
      });
      return;
    }
    const formData = new FormData();
    formData.append("avatar", avatarFile);
    try {
      const response = await fetch("http://localhost:8081/updateAvatar", {
        method: "POST",
        credentials: "include",
        body: formData,
      });
      if (!response.ok) {
        throw new Error("Erreur lors de la modification de l'avatar.");
      }
      toast.open({
        message: "Avatar modifié avec succès !",
        type: "success",
      });
      await dispatch("getMyProfileInfo");
    } catch (err) {
      console.error("Erreur lors de la modification de l'avatar:", err);
      toast.open({
        message: err.message,
        type: "error",
      });
    }
  },

  async deleteAccount({ dispatch }) {
    try {
      const response = await fetch("http://localhost:8081/DeleteAccount", {
        method: "DELETE",
        credentials: "include",
      });
      if (!response.ok) {
        throw new Error("Error deleting account.");
      }
      toast.open({
        message: "Account deleted successfully.",
        type: "success",
      });

      document.cookie.split(";").forEach((cookie) => {
        document.cookie =
          cookie.replace(/^ +/, "").split("=")[0] +
          "=;expires=Thu, 01 Jan 1970 00:00:00 UTC;path=/;";
      });

      // Optionally dispatch a logout action
      await dispatch("logout");

      // Redirect to sign-in page and force a reload
      window.location.href = "/sign-in";
      window.location.reload();
    } catch (err) {
      console.error("Error deleting account:", err);
      toast.open({
        message: err.message,
        type: "error",
      });
    }
  },

  createWebSocketConn({ commit, dispatch }) {
    const ws = new WebSocket("ws://localhost:8081/ws");

    ws.addEventListener("message", (e) => {
      const data = JSON.parse(e.data);
      console.log("WebSocket message received:", e.data);
      if (data.action === "chat") {
        const message = data.chatMessage;
        dispatch("addNewChatMessage", message);
        if (data.message === "NEW") {
          dispatch("fetchChatUserList");
        }
        if (message.type === "PERSON" || message.type === "GROUP") {
          dispatch("addUnreadChatMessage", message);
        }
      } else if (data.action === "notification") {
        dispatch("addNewNotification", data.notification);
      } else if (data.action === "groupAccept") {
        dispatch("getUserGroups");
      }
    });

    commit("updateWebSocketConn", ws);
  },

  addConversation({ commit }, message) {
    commit("addConversation", message);
  },
  deleteConversation({ commit }, message) {
    commit("deleteConversation", message);
  },
  deleteConversationById({ commit }, conversation_id) {
    commit('deleteConversationById', conversation_id);
  },
  clearChatHistory({ commit }) {
    commit("clearChatHistory");
  },
  addMessage({ commit }, message) {
    commit("addMessage", message);
  },
  // deletingMessage({ commit }, message) {
  //   commit("deletingMessage", message);
  // },
  clearMessages({ commit }) {
    commit("clearMessages");
  },
  getCurrentConvoID({ commit }) {
    commit("getCurrentConvoID");
  },
  updateCurrentConvoID({ commit }, id) {
    commit("updateCurrentConvoID", id);
  },
};
