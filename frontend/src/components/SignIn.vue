<template>
  <div class="sign-in__wrapper">
    <div class="bg-forest" :style="{ backgroundImage: `url(${backgroundImage})` }"></div>
    <div class="image-div-login" :style="{ backgroundImage: `url(${backgroundImage})` }"></div>

    <div class="sign-in">
      <h1 class="mori">Mori <span class="adder">- sign in</span></h1>
      <form class="form-group" @submit.prevent="signSubmit" id="sign-in__form">
        <div class="form-input">
          <label for="username">Email</label>
          <input type="email" id="email" v-model="signInForm.login" required>
        </div>
        <div class="form-input">
          <label for="password">Password</label>
          <input type="password" id="password" v-model="signInForm.password" required>
        </div>
      </form>
      <div class="button-or-register">
        <button class="btn" form="sign-in__form" type="submit">Sign in</button>
        <p>Need an account?
          <router-link to="/register" id="sign-up">Register here</router-link>
        </p>
        <p>Forgot your password ? 
          <router-link to="/forgotpassword" id="forgotPassword">click here</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "SignIn",
  data() {
    return {
      signInForm: {
        login: "",
        password: "",
      },
      backgroundImage: "", // Propriété pour stocker l'image sélectionnée
    };
  },
  methods: {
    setRandomImage() {
      const images = [
        require("../assets/images/forest_1.png"),
        require("../assets/images/sakura.webp"),
        require("../assets/images/automn.webp"),
        require("../assets/images/fantastic.webp"),
        require("../assets/images/fairytail.webp"),
      ];
      this.backgroundImage = images[Math.floor(Math.random() * images.length)];
    },
    async signSubmit() {
      try {
        await fetch("http://localhost:8081/signin", {
          credentials: "include",
          method: "POST",
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
          },
          body: JSON.stringify(this.signInForm),
        })
          .then((response) => response.json())
          .then((json) => {
            if (json.message === "Login successful") {
              this.$toast.open({
                message: "Login successful!",
                type: "success",
              });

              this.$store.dispatch("createWebSocketConn").then(() =>
                this.$router.push("/main")
              );
            } else {
              this.$router.push("/");
              this.$toast.open({
                message: json.message,
                type: "error",
              });
            }
          });
      } catch {
        // Handle errors
      }
    },
  },
  created() {
    this.setRandomImage(); // Choisir une image aléatoire à la création
  },
};
</script>

<style>

.bg-forest {
  position: absolute;
  height: 100vh;
  width: 100vw;
  left: 0;
  z-index: -1;
  filter: blur(8px) brightness(70%);
  -webkit-filter: blur(8px) brightness(70%);
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.sign-in__wrapper {
  display: flex;
  background-color: var(--bg-neutral);
  border-radius: 20px;
  color: var(--color-white);
  box-shadow: 0 4px 15px rgb(0, 0, 0);
  overflow: hidden;
  align-items: center;
  width: 810px;
  max-width: 90%; /* Added for smaller screens */
}

.image-div-login {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 50%;
  height: 100%;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.sign-in {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 40px;
  margin: 0 auto;
  padding: 8%;
  width: 60%;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 15px;
  color: var(--color-white);
  width: 100%;
}

.form-input {
  color: var(--color-white);
}

.sign-in button {
  margin-bottom: 10px;
  width: 80px;
}

#sign-up {
  color: var(--purple-color);
  text-decoration: underline;
  font-size: 16.5px;
  transition: all 0.3s ease;
}

#sign-up:hover {
  color: var(--hover-background-color);
  text-decoration: underline;
  font-size: 16.5px;
}

#forgotPassword{
  color: var(--purple-color);
  text-decoration: underline;
  font-size: 16.5px;
  transition: all 0.3s ease;
}

#forgotPassword:hover {
  color: var(--hover-background-color);
  text-decoration: underline;
  font-size: 16.5px;
}

.button-or-register {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Media Queries */

/* Tablet and Phone View (768px and below) */
@media (max-width: 900px) {
  html, body {
    overflow-y: hidden; /* Prevent horizontal scrolling */
  }

  .bg-forest {
    filter: blur(8px) brightness(60%); /* Softer background for clarity */
    top: 0;
  }

  .sign-in__wrapper {
    flex-direction: column;
    align-items: center;
    width: 95%;
    max-width: 600px;
    margin: 0 auto;
    border-radius: 10px;
  }

  .image-div-login {
    width: 100%;
    height: 200px;
    display: block;
  }

  .sign-in {
    width: 90%;
    padding: 5%;
  }

  h1.mori {
    font-size: 1.8rem; /* Adjust heading size for smaller screens */
  }
  .sign-in button {
    display: flex;
    width: 100%; /* Wider button for better touch interaction */
    align-items: center;
    justify-content: center;
    text-align: center;
    font-size: 15px;
  }
}

/* Phone View (480px and below) */
@media (max-width: 480px) {
  .sign-in__wrapper {
    flex-direction: column; /* Keep the same column layout */
    width: 95%;
    margin: 0 auto;
  }

  .image-div-login {
    width: 100%; /* Ensure image spans full width */
    height: 150px; /* Slightly smaller height for phone screens */
    display: block;
  }

  .sign-in {
    width: 95%; /* Slightly wider form for phones */
    padding: 4%;
  }

  h1.mori {
    font-size: 1.5rem; /* Adjust heading size for better readability */
  }

  .bg-forest {
    filter: blur(8px) brightness(60%); /* Softer background for clarity */
    top: 0;
  }

  .sign-in button {
    display: flex;
    width: 100%; /* Wider button for better touch interaction */
    align-items: center;
    justify-content: center;
    text-align: center;
    font-size: 15px;
  }

  p {
    font-size: 14px; /* Smaller font size for better readability */
  }
}

</style>