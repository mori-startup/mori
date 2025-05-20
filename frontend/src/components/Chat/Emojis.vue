<template>
  <i class="uil uil-smile" @click="toggleShowEmojis"></i>
  <div class="emojis" v-show="showEmojis">
    <p class="emoji" v-for="emoji in getEmojis" @click.stop="addEmoji(emoji)">
      {{ emoji }}
    </p>
  </div>
</template>

<script>
export default {
  props: ["input", "messagebox"],
  data() {
    return {
      emojiCodes: [128512, 128514, 128519, 128520],
      showEmojis: false,
    };
  },

  computed: {
    getEmojis() {
      return this.emojiCodes.map((code) => String.fromCodePoint(code));
    },
  },

  methods: {
    addEmoji(emoji) {
      this.input.value += emoji;
    },

    toggleShowEmojis() {
      this.showEmojis = !this.showEmojis;
    },
  },
};
</script>

<style scoped>
i {
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  font-size: 1.4em;
  width: 5%;
  min-width: 40px;
  height: 100%;
  border-radius: 10px;
  background-color: var(--purple-color);
  color: white;
  transition: all 0.3s ease;
}

i:hover {
  background-color: var(--hover-background-color);
}

.emojis {
  display: flex;
  padding-top: 5px;
  padding-left: 7px;
  width: 40%;
  gap: 5px;
  opacity: 0;
  animation: appearFromRight 0.6s forwards;
}

.emoji {
  font-size: 1.3em;
  cursor: pointer;
  transition: all 0.3s ease;
}

.emoji:hover {
  transform: scale(1.2);
}

@keyframes appearFromRight {
  0% {
    opacity: 0;
    width: 0%;
    transform: translateX(50px);
  }

  20% {
    opacity: 0;
  }

  100% {
    opacity: 1;
    width: 40%;
    transform: translateX(0);
  }
}
</style>
