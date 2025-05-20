<template>
  <!-- <h3>Historique du Chatbot</h3> -->
  <button class="btn_new_convo" @click.stop="newConversation">
    New Conversation
  </button>
  <ul id="list_chat_convo">
    <li
      v-for="(convo, index) in conversationsUpdate"
      :key="convo.conversation_id"
      :class="{ selected: selectedConvoId === convo.conversation_id }"
    >
     
      <button @click.stop="loadConvo(convo.conversation_id)" class="btn_elmt_history">
          {{ convo.convo[convo.convo.length - 1].user_request }}  
      </button>
      <button
        class="btn_delete_convo"
        @click.stop="deleteConvo(convo.conversation_id)"
      >
        X
      </button>
    </li>
  </ul>
</template>

<script>
export default {
  data() {
    return {
      userInput: "",
      selectedConvoId: null, // Track selected conversation
    };
  },
  computed: {
    conversationsUpdate() {
      // Always return a new reversed array, don't mutate chatHistory directly
      console.log("computed conversationsUpdate: ", this.$store.getters.allConversations);
      // return [...this.$store.getters.allConversations].reverse();
      return [...this.$store.getters.allConversations].reverse();
    },
  },
  mounted() {
    this.getHistory();
  },
  methods: {
    async getMyUserID() {
      const response = await fetch("http://localhost:8081/currentUser", {
        credentials: "include",
        headers: new Headers({
          "Content-Type": "application/json",
        }),
        method: "POST",
      });
      if (!response.ok) {
        console.error(
          "Erreur lors de la récupération de l'ID de l'utilisateur :",
          response.statusText
        );
        return;
      } else {
        const resp = await response.json();
        console.log(resp);
        console.log(resp.users[0].id);
        return resp.users[0].id;
      }
    },

    newConversation() {
      this.selectedConvoId = null;
      this.$store.dispatch("clearMessages");
    },

    async deleteConvo(conversation_id) {
      console.log("conversation_id FOR DELETING: ", conversation_id);
  
      const response = await fetch("http://localhost:8081/llmConvoDelete", {
        credentials: "include",
        headers: new Headers({
          "Content-Type": "application/json",
        }),
        method: "POST",
        body: JSON.stringify({
          user_id: await this.getMyUserID(),
          conversation_id: conversation_id,
        }),
      });
      if (!response.ok) {
        console.error(
          "Erreur lors de la suppression de la conversation de l'utilisateur :",
          response.statusText
        );
        return;
      } else {
        // Supprime la conversation du store par son ID
        this.$store.dispatch("deleteConversationById", conversation_id);
        this.$store.dispatch("clearMessages");
        // Si la conversation supprimée était sélectionnée, on désélectionne
        if (this.selectedConvoId === conversation_id) {
          this.selectedConvoId = null;
        }
        console.log("Conversation supprimée");
      }
    },

    //Méthode pour obtenir la discussion selctionnée
    async loadConvo(convo_id) {
      console.log("convo_id in loadConvo: ", convo_id);
      const response = await fetch("http://localhost:8081/llmConvoSelected", {
        credentials: "include",
        headers: new Headers({
          "Content-Type": "application/json",
        }),

        method: "POST",
        body: JSON.stringify({
          user_id: await this.getMyUserID(),
          conversation_id: convo_id,
        }),
      });
      if (!response.ok) {
        console.error(
          "Erreur lors de la récupération de la conversation de l'utilisateur :",
          response.statusText
        );
        return;
      } else {
        const resp = await response.json();
        console.log("RESPONSE loadConvo: ", resp);
        console.log("Current convo: ", resp.convo);
        console.log("Current convo_id loadConvo: ", resp.conversation_id);
        this.$store.dispatch("clearMessages");
        this.convertMessages(resp);
      }
    },

    convertMessages(convo) {
      console.log("convo in convertMessages: ", convo);

      if (convo === null) {
        console.log("convo is null ");
        return;
      }
      convo.convo.reverse().forEach((message) => {
        this.$store.dispatch("addMessage", {
          sender: "Utilisateur",
          text: message.user_request,
          conversation_id: convo.conversation_id,
        });
        this.$store.dispatch("addMessage", {
          sender: "LLM",
          text: message.llm_response,
          conversation_id: convo.conversation_id,
        });
      });
    },
    convertLastMessages(convo, conversation_id) {
      console.log("convo in convertLastMessages: ", convo);

      convo.forEach((c) => {
        if (c.conversation_id === conversation_id) {
          c.convo.forEach((message) => {
            this.$store.dispatch("addMessage", {
              sender: "Utilisateur",
              text: message.user_request,
              conversation_id: message.conversation_id,
            });
            this.$store.dispatch("addMessage", {
              sender: "LLM",
              text: message.llm_response,
              conversation_id: message.conversation_id,
            });
          });
        }
      });
    },

    async getHistory() {
      const response = await fetch("http://localhost:8081/llmConvoGet", {
        credentials: "include",
        headers: new Headers({
          "Content-Type": "application/json",
        }),
        method: "POST",
        body: JSON.stringify({ user_id: await this.getMyUserID() }),
      });
      if (!response.ok) {
        console.error(
          "Erreur lors de la récupération des conversations de l'utilisateur :",
          response.statusText
        );
        return;
      }
      const resp = await response.json();
      console.log("RESPONSE getHistory: ", resp);
      this.$store.dispatch("clearMessages");
      this.$store.dispatch("clearChatHistory");
      resp.forEach((convo) => {
        this.$store.dispatch("addConversation", convo);
      });
      // Ne touche plus à this.chatHistory ici !
      console.log("Conversations ajoutées au store.");
    },
  },
};
</script>

<style>
#app {
  display: flex;
  height: 100vh;
}

.sidebar {
  width: 25%;
  background: #f4f4f4;
  border-right: 1px solid #ddd;
  padding: 10px;
}

.sidebar ul {
  list-style: none;
  padding: 0;
}

.sidebar li {
  cursor: pointer;
  padding: 5px;
  border-bottom: 1px solid #ddd;
}

.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 10px;
}

.chat-content {
  flex: 1;
  overflow-y: auto;
  margin-bottom: 10px;
}

.user-message {
  text-align: right;
  background: #d1ffd1;
  margin: 5px;
  padding: 5px 10px;
  border-radius: 5px;
}

.bot-message {
  text-align: left;
  background: #f0f0f0;
  margin: 5px;
  padding: 5px 10px;
  border-radius: 5px;
}

.btn_elmt_history {
  padding: 10px;
  border: 1px solid #ddd;
  margin: 5px;
  border-radius: 5px;
  background-color: black;
  color: white;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.btn_elmt_history:hover {
  background-color: none;
  border-color: none;
}

input[type="text"] {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  outline: none;
}
.selected .elmt_history {
  background-color: #333;
  color: #fff;
}
</style>

<style>
#list_chat_convo {
  list-style: none;
  padding: 0;
}

#list_chat_convo li {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  margin-bottom: 5px;
  border-radius: 10px;
  background-color: #333; /* Couleur de fond sombre */
  transition: background-color 0.3s ease;
  color: #fff; /* Couleur du texte blanche */
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.3); /* Ombre pour donner du relief */
}

#list_chat_convo li.selected {
  background-color: #444; /* Couleur légèrement plus claire pour l'élément sélectionné */
}

.btn_new_convo {
  display: inline-block;
  padding: 10px 20px;
  margin-bottom: 10px;
  border-radius: 10px;
  background-color: var(--purple-color); /* Utilisation de la couleur violette */
  color: #fff;
  border: none;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.3s ease;
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.3);
}

.btn_new_convo:hover {
  background-color: var(--hover-color); /* Couleur de survol */
  transform: scale(1.05);
}

.btn_elmt_history {
  flex-grow: 1;
  text-align: left;
  background-color: transparent;
  border: none;
  color: inherit;
  cursor: pointer;
  padding: 0;
}

.btn_delete_convo {
  background-color: #ff4d4d;
  border: none;
  color: white;
  padding: 5px 10px;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.btn_delete_convo:hover {
  background-color: #ff1a1a;
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
  width: 50px; /* Ajusté pour correspondre à l'image */
  height: 50px; /* Ajusté pour correspondre à l'image */
  background-color: var(--purple-color);
  border-radius: 50%;
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
</style>
