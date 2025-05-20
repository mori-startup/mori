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
          v-for="(message, index) in messages"
          :key="index"
          :class="['message', message.sender === 'Utilisateur' ? 'Utilisateur' : 'LLM']"
        >
          <div v-if="message.sender !== 'Utilisateur'" class="bot-logo">
            <img src="../assets/mori.png" alt="Bot Logo" />
          </div>
          <pre>{{ message.text }}</pre>
          <div class="timestamp">{{ message.timestamp }}</div>
        </div>
      </div>

      <div :class="['chatbot-input', { 'chatbot-input--active': hasMessages }]">
        <textarea
          :rows="rows"
          class="chatbot-input"
          type="text"
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
export default {
  data() {
    return {
      userInput: "",
      messages: [],
      rows: 10,
      count: 0,
    };
  },
  computed: {
    hasMessages() {
      return this.messages.length > 0;
    },
  },
  mounted() {
    this.initializeConversation();
  },
  methods: {
    initializeConversation() {
      this.conversation = {
        user_id:"",
        conversation_id: "",
        user_request: "",
        llm_response: "",
        session: "",
        new_conversation: false,
        history: [],
      };
    },
    appendMessage(sender, text) {
      const timestamp = new Date().toLocaleTimeString();
      this.messages.push({ sender, text, timestamp });
      this.$nextTick(() => {
        const chatBox = this.$el.querySelector(".chatbot-messages");
        chatBox.scrollTop = chatBox.scrollHeight;
      });
    },
    async sendMessage() {
      if (this.userInput.trim() === "") return;

      this.appendMessage("Utilisateur", this.userInput);
      this.conversation.user_request = this.userInput;
      this.userInput = "";

      try {
        await this.sendData();
        // this.userInput = ""; // Clear the input
      } catch (error) {
        console.error("Erreur lors de l'envoi du message :", error);
      }
    },
    async sendData() {
      const response = await fetch(`http://localhost:8081/llmConvo`, {
        method: "POST",
        credentials: 'include',
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(this.conversation),
      });
      
      if (!response.ok) {
        console.error("Erreur lors de l'envoi des données :", response.statusText);
        return;
      }
      const reader = response.body.getReader();
      const decoder = new TextDecoder("utf-8");
      
      try {
      // const llm_bubble = document.createElement("div");
      // const chatBox = this.$el.querySelector(".chatbot-messages");
      // llm_bubble.classList.add("message","LLM");

      // const bot_image = document.createElement("div");
      // bot_image.classList.add("bot-logo");
      // const bot_img = document.createElement("img");
      // bot_img.src = "../assets/mori.png";
      // bot_image.appendChild(bot_img);
      // llm_bubble.appendChild(bot_image);
      
      // const innerPre = document.createElement("pre");

      // const timestamp = new Date().toLocaleTimeString();
      // const divTimestamp = document.createElement("div");
      // divTimestamp.classList.add("timestamp");
      // divTimestamp.innerText = timestamp;
      // llm_bubble.appendChild(divTimestamp);
      
      // llm_bubble.appendChild(innerPre);
      
      let { done, value } = await reader.read();
      
      while (!done) {
        const chunk = decoder.decode(value, { stream: true });
        const lines = chunk.split("\n");
        lines.forEach((line) => {
          if (line.startsWith("data: ")) {
            const jsonData = line.replace("data: ", "").trim();
            try {
              const parsedData = JSON.parse(jsonData);
              
              innerPre.innerText += parsedData.response.message.content;
              
            } catch (error) {
              console.error("Erreur de parsing JSON :", error);
            }
          }
        });
        ({ done, value } = await reader.read());
      }
      
      
      // chatBox.appendChild(llm_bubble);
      this.conversation.llm_response = innerPre.innerText;
      this.conversation.history.push(this.conversation.llm_response);
      await this.sendConversation();
    } catch (error) {
        console.error("Erreur de lecture du flux", error);
      } finally {
        reader.releaseLock();
      }
    },
    async sendConversation() {
      const response = await fetch(`http://localhost:8081/llmConvo`, {
        method: "POST",
        credentials: 'include',
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(this.conversation),
      });

      if (!response.ok) {
        console.error("Erreur lors de l'envoi de la conversation :", response.statusText);
        return;
      }

      console.log("Conversation envoyée avec succès !");
    },
    handleKeydown(event) {
      
      if (event.shiftKey && event.key === "Enter" ) {
        event.preventDefault();
        this.userInput += "\n";
        let textarea = this.$el.querySelector("textarea");
        textarea.style.height = `${textarea.scrollHeight+10}px`;
        // this.count = this.userInput.split('\n').length;

      } else if (event.key === "Enter") {
        event.preventDefault();
        this.sendMessage();
      }
      // event.preventDefault();
      // let textarea = this.$el.querySelector("textarea");
      // if (event.key === "Backspace" && textarea.selectionEnd ===  ) {
      //   this.userInput += "LOL"
      //   textarea.style.height = `${textarea.scrollHeight-10}px`;
      // }
      
    },
  },
};
</script>

<style scoped>
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
  
  .message{
    max-width: 70%;
    padding: 10px;
    border-radius: 10px;
    font-size: 16px;
    position: relative;
  }
  
  .Utilisateur {
    align-self: flex-end;
    background-color: var(--purple-color);
    color: var(--color-white);
  }
  
  .LLM {
    align-self: flex-start;
    background-color: red;
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
  .chatbot-input textarea {
    display: flex;
    gap: 10px;
    align-items: center;
    position: absolute;
    top: 65%; /* Initially positioned below the "How can I help?" */
    left: 50%;
    transform: translate(-50%, -50%);
    width: calc(40% - 40px);
    border-radius: 10px;
    padding: 10px 20px;
    transition: top 0.7s ease, transform 0.7s ease, width 0.7s ease; /* Added width for smooth transition */
  }
  
  .chatbot-input--active {
    width: calc(50% - 40px); /* New width for active state */
    position: fixed;
    top: calc(97% - 80px); /* Slide to the bottom of the viewport */
    transform: translateX(-50%);
  }
  
  .chatbot-input {
    flex: 1;
    padding: 10px;
    border: 1px solid var(--color-grey);
    border-radius: 10px;
    height: 50px;
    font-size: 16px;
  }

  textarea {
    resize: none;
    word-wrap: break-word;
    overflow: hidden;
  }
  
  .chatbot-input button {
    padding: 15px 20px;
    background-color: var(--purple-color);
    color: var(--color-white);
    border: none;
    border-radius: 10px;
    cursor: pointer;
    font-size: 16px;
    transition: background-color 0.3s;
  }
  
  .chatbot-input button:hover {
    background-color: var(--hover-background-color);
  }
  pre {
    white-space: pre-wrap;
    word-wrap: break-word;
  }
  
  .chatbot-input textarea:focus {
    outline: none;
    border-color: var(--color-primary);
  }</style>