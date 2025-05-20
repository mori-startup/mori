<template>
  <div class="chatbot-container">
    <div :class="['chatbot-box', { 'chatbot-box--active': hasMessages }]">
      <div class="mori-img" v-if="!hasMessages">
        <div class="moriImg"></div>
      </div>
      <div class="mori" id="moriChatBot" v-if="!hasMessages">Mori</div>
      <div class="chatbot-message" v-if="!hasMessages">How can I help you?</div>
      <div class="chatbot-messages" v-if="hasMessages">
        <div
          v-for="(message, index) in allMessages"
          :key="index"
          :class="[
            'message',
            message.sender === 'Utilisateur' ? 'Utilisateur' : 'LLM',
          ]"
        >
          <div v-if="message.sender !== 'Utilisateur'" class="bot-logo">
            <img src="../assets/mori.png" alt="Bot Logo" />
          </div>
          <div
            v-if="message.sender === 'Utilisateur'"
            class="markdown-container-Utilisateur"
          >
            <Markdown :source="message.text" />
          </div>
          <div v-if="message.sender === 'LLM'" class="markdown-container-LLM">
            <Markdown class="markdownLLM" :source="message.text" />
          </div>
          <div class="timestamp">{{ message.timestamp }}</div>
        </div>
      </div>

      <div
        :class="[
          'chatbot-input-container',
          { 'chatbot-input-container--active': hasMessages },
        ]"
      >
        <textarea
          ref="textarea"
          :rows="rows"
          class="chatbot-textarea"
          v-model="userInput"
          @keydown="handleKeydown"
          @click="handleKeydown"
          placeholder="Type your message here..."
        ></textarea>
        <button @click="sendMessage">Send</button>
      </div>
    </div>
  </div>
</template>

<script>
import Markdown from "vue3-markdown-it";
const { v4: uuidv4 } = require("uuid");

export default {
  components: { Markdown },

  data() {
    return {
      userInput: "",
      messages: [],
      current_convID: "",
      rows: 10,
      sourceLLM: "",
      sourceUtilisateur: "",
      markdownText: "",
      conversation: {
        user_id: "",
        conversation_id: "",
        convo: [],
        new_conversation: false,
      },
    };
  },
  computed: {
    hasMessages() {
      return this.$store.getters.allMessages.length > 0;
    },
    allMessages() {
      this.messages = this.$store.getters.allMessages;
      return this.messages;
    },
  },
  mounted() {
    // this.loadCurrentConvo();
  },
  methods: {
    //Méthode pour récupérer les messages de la conversation selectionné
    getCurrentMessages() {
      this.messages = this.$store.getters.allMessages;
    },
    async loadCurrentConvo() {
      const convo_id = localStorage.getItem("current_convo_id");
      console.log("convo_id in loadConvo: ", convo_id);
    
      if (convo_id === null) {
        console.log("No conversation selected");
        return;
      }
    
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
        console.log("Current convo: ", resp.convo);
        this.$store.dispatch("clearMessages");
        this.convertMessages(resp);
      }
    },

    //Méthode pour récupérer l'ID de l'utilisateur
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
    appendMessage(sender, text) {
      let dict = {};
      dict = {
        sender,
        text,
        conversation_id: "",
      };

      this.$store.dispatch("addMessage", dict);

      this.$nextTick(() => {
        const chatBox = this.$el.querySelector(".chatbot-messages");
        chatBox.scrollTop = chatBox.scrollHeight;
      });
    },
    async sendMessage() {
      if (this.userInput.trim() === "") return;
      this.appendMessage("Utilisateur", this.userInput);
      this.conversation.convo.push({
        user_request: this.userInput,
        llm_response: "",
      });
      this.userInput = "";

      try {
        await this.sendData();
      } catch (error) {
        console.error("Erreur lors de l'envoi du message :", error);
      }
    },
    async sendData() {
      let accumulatedText = "";
      const response = await fetch(`http://localhost:8081/llmConvo`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(this.conversation),
      });

      if (!response.ok) {
        console.error(
          "Erreur lors de l'envoi des données :",
          response.statusText
        );
        return;
      }
      const reader = response.body.getReader();
      const decoder = new TextDecoder("utf-8");
      this.appendMessage("LLM", "");
      let lastLLMMessage =
        this.$store.getters.allMessages[
          this.$store.getters.allMessages.length - 1
        ]; // Référence au dernier message LLM

      try {
        accumulatedText = "";
        let { done, value } = await reader.read();

        while (!done) {
          const chunk = decoder.decode(value, { stream: true });
          const lines = chunk.split("\n");
          lines.forEach((line) => {
            if (line.startsWith("data: ")) {
              const jsonData = line.replace("data: ", "").trim();
              try {
                const parsedData = JSON.parse(jsonData);

                accumulatedText += parsedData.response;
                lastLLMMessage.text = accumulatedText;
              } catch (error) {
                console.error("Erreur de parsing JSON :", error);
              }
            }
          });
          ({ done, value } = await reader.read());
        }

        // this.conversation.llm_response = accumulatedText;

        this.appendMessage("LLM", accumulatedText);
        this.messages.splice(this.messages.length - 2, 1);
        this.$store.dispatch("deletingMessage", this.messages.length - 2);

        lastLLMMessage.text = "";
      } catch (error) {
        console.error("Erreur de lecture du flux", error);
      } finally {
        reader.releaseLock();
      }

      this.conversation.user_id = await this.getMyUserID();
      this.conversation.convo[this.conversation.convo.length - 1].llm_response =
        accumulatedText;

      console.log(" Number of allMessages(): ", this.messages.length);
      if (this.messages.length <= 2) {
        this.conversation.conversation_id = uuidv4();
        this.messages.forEach((message) => {
          message.conversation_id = this.conversation.conversation_id;
        });
        this.conversation.new_conversation = true;
        this.$store.dispatch("addConversation", this.conversation);
      } else {
        console.log("Messages before saving: ", this.messages);
        this.conversation.conversation_id = this.messages[0].conversation_id;
        this.conversation.new_conversation = false;
        this.addMessageToExistingConversation(
          this.$store.getters.allConversations,
          this.conversation
        );
      }

      console.log("Conversation : ", this.conversation);
      this.sendConversation();
    },
    async sendConversation() {
      const response = await fetch(`http://localhost:8081/llmConvoSave`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(this.conversation),
      });

      if (!response.ok) {
        console.error(
          "Erreur lors de l'envoi de la conversation :",
          response.statusText
        );
        return;
      }
      console.log("Conversation envoyée avec succès");
      this.conversation = {
        user_id: "",
        conversation_id: "",
        convo: [],
        new_conversation: false,
      };
    },
    //Méthode pour gérer les messages dans une conversation existante
    addMessageToExistingConversation(allConversation, currentConversation) {
      console.log(
        "Conversation in addMessageToExistingConversation: ",
        allConversation
      );
      console.log("Current conversation: ", currentConversation);
      console.log(
        "Current conversation ID: ",
        currentConversation.conversation_id
      );

      for (let t = 0; t < allConversation.length; t++) {
        const convo = allConversation[t];
        if (convo.conversation_id === currentConversation.conversation_id) {
          console.log("Conversation found");
          console.log("Current conversation: ", currentConversation.convo);
          currentConversation.convo.push(...convo.convo);
        }
      }
      console.log("Current conversation after: ", currentConversation.convo);
    },

    // Méthode pour gérer les événements de touche
    handleKeydown(event) {
      if (event.shiftKey && event.key === "Enter") {
        event.preventDefault();
        this.userInput += "\n";
        let textarea = this.$el.querySelector("textarea");
        textarea.style.height = `${textarea.scrollHeight + 10}px`;
      } else if (event.key === "Enter") {
        event.preventDefault();
        this.sendMessage();
      }
      // let textarea = this.$el.querySelector("textarea");
      let textarea = document.querySelector("textarea");
      const textLength = textarea.value.length;
      if (event.key === "Backspace" && textLength > 0) {
        const cursorPosition = textarea.selectionEnd; // Position actuelle du curseur

        // Vérifie si le caractère à supprimer est un retour chariot
        if (textarea.value[cursorPosition - 1] === "\n") {
          // Réduit la hauteur du textarea
          textarea.style.height = `${textarea.scrollHeight - 22}px`;
        }
      } else if (event.key === "Backspace" && textLength === 1) {
        textarea.style.height = `50px`;
      }
    },
    convertMessages(convo) {
      console.log("convo in convertMessages: ", convo);

      convo.convo.forEach((message) => {
        this.messages.push({
          sender: "Utilisateur",
          text: message.user_request,
          conversation_id: convo.conversation_id,
        });
        this.messages.push({
          sender: "LLM",
          text: message.llm_response,
          conversation_id: convo.conversation_id,
        });
      });
    },
  },
};
</script>

<style scoped>
  .Utilisateur {
    align-self: flex-end;
    background-color: var(--purple-color);
    color: var(--color-white);
  }
  
  .LLM {
    align-self: flex-start;
    text-align: left;
    background-color: var(--page-bg);
    color: var(--color-white);
  }
  
  .chatbot-container {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
    background-color: var(--page-bg);
    font-family: Arial, sans-serif;
  }
  
  .chatbot-box {
    display: flex;
    flex-direction: column;
    width: 100%;
    max-width: 800px;
    border-radius: 10px;
    padding: 30px;
    text-align: center;
    gap: 20px;
    margin-bottom: 130px;
    transition: all 0.5s ease;
  }
  
  .chatbot-box--active {
    justify-content: space-between;
    width: 80%;
    height: 85vh;
  }

  .bot-logo {
    display: inline-block;
    vertical-align: top;
    margin-right: 10px;
    margin-top: -5px;
    margin-left: -5px;
  }
  
  .bot-logo img {
    width: 35px; /* Adjust size as needed */
    height: 35px; /* Adjust size as needed */
    border-radius: 50%; /* Optional: Make the image circular */
    background-color: var(--purple-color);
    object-fit: cover; /* Ensure the image scales properly */
  }
  
  
  .mori-img {
    display: flex;
    justify-content: center;
    align-items: center;
    transition: opacity 0.5s ease;
  }
  
  #moriChatBot {
    user-select: none;
    font-size: 50px;
    font-weight: bold;
    transition: opacity 0.5s ease;
  }
  
  .chatbot-message {
    user-select: none;
    margin-bottom: 20px;
    font-size: 20px;
    color: var(--color-white);
    transition: opacity 0.5s ease;
  }
  
  .chatbot-messages {
    flex: 1;
    overflow-y: auto;
    padding: 10px;
    border-radius: 10px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  
  .message {
    max-width: 70%;
    padding: 10px;
    border-radius: 10px;
    font-size: 16px;
    position: relative;
  }
  
  .user {
    align-self: flex-end;
    background-color: var(--purple-color);
    color: var(--color-white);
  }
  
  .bot {
    align-self: flex-start;
    background-color: var(--bg-neutral);
    color: var(--color-white);
  }
  
  .timestamp {
    font-size: 12px;
    color: var(--color-grey);
    opacity: 0.5;
    text-align: right;
    margin-top: 5px;
  }
  
  /* Input field animation */
  /* 
  1. The container that slides down with an animation 
     (replaces .chatbot-input in your old code)
*/
.chatbot-input-container {
  display: flex;
  gap: 10px;
  align-items: center;
  position: absolute;
  top: 65%; /* Initially below the greeting message */
  left: 50%;
  transform: translate(-50%, -50%);
  width: calc(40% - 40px);
  border-radius: 10px;
  padding: 10px 20px;
  transition: top 0.7s ease, transform 0.7s ease, width 0.7s ease;
}

.chatbot-input-container--active {
  width: calc(50% - 40px); /* Widen the container */
  position: absolute;
  top: calc(97% - 80px);   /* Slide to bottom of viewport */
  transform: translateX(-50%);
}

/* 
  2. The textarea itself: 
     (new .chatbot-textarea class)
*/
.chatbot-textarea {
  flex: 1;
  border: 1px solid var(--color-grey);
  border-radius: 10px;
  font-size: 16px;
  min-height: 50px;   /* Ensure it matches your old input height */
  padding: 13px;
  resize: none;       /* Optional: remove manual resize handle */
}

/* 
  3. The Send button 
  (same rules as your old .chatbot-input button style)
*/
.chatbot-input-container button {
  padding: 15px 20px;
  background-color: var(--purple-color);
  color: var(--color-white);
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s;
}

.chatbot-input-container button:hover {
  background-color: var(--hover-background-color);
}

  
  .chatbot-input-container input:focus {
    outline: none;
    border-color: var(--color-primary);
  }
</style>
