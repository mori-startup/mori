<template>
  <div class="register__wrapper">
    <div
      class="bg-forest-reg"
      :style="{ backgroundImage: `url(${backgroundImage})` }"
    ></div>
    <div
      class="image-div-reg"
      :style="{ backgroundImage: `url(${backgroundImage})` }"
    ></div>

    <div class="register">
      <h1 class="mori">Mori <span class="adder">- Register</span></h1>
      <form
        class="form-group"
        @submit.prevent="submitRegData"
        id="register__form"
      >
        <div class="form-row">
          <div class="form-input">
            <label for="firstname">First name</label>
            <input
              v-model="form.firstname"
              type="text"
              name="firstname"
              id="firstname"
              required
            />
          </div>
          <div class="form-input">
            <label for="lastname">Last name</label>
            <input
              v-model="form.lastname"
              type="text"
              name="lastname"
              id="lastname"
              required
            />
          </div>
        </div>

        <div class="form-row">
          <div class="form-input">
            <label for="email">Email</label>
            <input
              v-model="form.email"
              type="email"
              name="email"
              id="email"
              required
            />
          </div>
          <div class="form-input">
            <label for="password">Password</label>
            <input
              v-model="form.password"
              type="password"
              name="password"
              id="password"
              required
            />
          </div>
        </div>

        <div class="form-row">
          <div class="form-input">
            <label for="date">Date of Birth</label>
            <input
              v-model="form.dateofbirth"
              type="date"
              name="date"
              id="date"
              required
            />
          </div>
          <div class="form-input">
            <label for="nickname">Nickname</label>
            <input
              v-model="form.nickname"
              type="text"
              name="nickname"
              id="nickname"
            />
          </div>
        </div>

        <div class="form-row">
          <div class="form-input">
            <label for="aboutme">About me</label>
            <textarea
              v-model="form.aboutme"
              id="aboutme"
              name="aboutme"
              rows="4"
            ></textarea>
          </div>
          <div class="form-input">
            <FileUpload v-model:file="form.avatar" labelName="Avatar" />
          </div>
        </div>

        <!-- Captcha Row -->
        <div class="form-row">
          <div class="form-input">
            <!-- Display the captcha image -->
            <img
              :src="captchaImageUrl"
              alt="Captcha"
              style="margin-bottom: 5px; display: block"
            />
            <div id="captchaRow">
              <input
                v-model="form.captchaValue"
                type="text"
                name="captcha"
                id="captcha"
                placeholder="Enter the code above"
                required
              />
              <button
                type="button"
                class="btn"
                @click="reloadCaptcha"
                style="font-size: 1.2rem; height: 44px; width: 44px; text-align: center"
              >
              ⟳
              </button>
            </div>
          </div>
        </div>
        <!-- End Captcha Row -->

        <div class="button-or-signIn">
          <button class="btn" form="register__form" type="submit">
            Create account
          </button>
          <p>
            Already have an account?
            <router-link to="/sign-in" id="sign-up">Sign-in</router-link>
          </p>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import forest1 from "../assets/images/forest_1.png";
import sakura from "../assets/images/sakura.webp";
import automn from "../assets/images/automn.webp";
import fantastic from "../assets/images/fantastic.webp";
import fairytail from "../assets/images/fairytail.webp";
import FileUpload from "./FileUpload.vue";

export default {
  name: "Register",
  components: { FileUpload },
  data() {
    return {
      form: {
        email: "",
        password: "",
        firstname: "",
        lastname: "",
        dateofbirth: null,
        nickname: "",
        avatar: null,
        aboutme: "",
        // We'll store captcha input in this field
        captchaValue: "",
      },
      backgroundImage: "", // Randomly chosen background image
      // We'll build this dynamic URL so we can reload the captcha
      captchaImageUrl: "http://localhost:8081/captcha",
    };
  },
  methods: {
    setRandomImage() {
      const images = [forest1, sakura, automn, fantastic, fairytail];
      this.backgroundImage = images[Math.floor(Math.random() * images.length)];
    },

    // Reload Captcha: append a timestamp so the browser doesn't cache
    reloadCaptcha() {
      this.captchaImageUrl = `http://localhost:8081/captcha?${Date.now()}`;
    },

    async submitRegData() {
      // Basic password checks (as you already do)
      const pwd = this.form.password;
      if (pwd.length < 10) {
        this.$toast.open({
          message: "Password must be at least 10 characters.",
          type: "error",
        });
        return;
      }
      if (!/[A-Z]/.test(pwd)) {
        this.$toast.open({
          message: "Password must contain at least one uppercase letter.",
          type: "error",
        });
        return;
      }
      if (!/\d/.test(pwd)) {
        this.$toast.open({
          message: "Password must contain at least one digit.",
          type: "error",
        });
        return;
      }
      if (!/[^a-zA-Z0-9]/.test(pwd)) {
        this.$toast.open({
          message: "Password must contain at least one special character.",
          type: "error",
        });
        return;
      }

      // Build FormData
      const formData = new FormData();
      formData.set("avatar", this.form.avatar);
      formData.set("email", this.form.email);
      formData.set("password", this.form.password);
      formData.set("firstname", this.form.firstname);
      formData.set("lastname", this.form.lastname);
      formData.set("dateofbirth", this.form.dateofbirth);
      formData.set("nickname", this.form.nickname);
      formData.set("aboutme", this.form.aboutme);
      // Set the captcha user input
      formData.set("captchaValue", this.form.captchaValue);

      try {
        // Because the server sets a "captcha_id" cookie, we need "credentials: include"
        const res = await fetch("http://localhost:8081/register", {
          credentials: "include",
          method: "POST",
          body: formData,
        });

        if (res.status === 409) {
          this.$toast.open({
            message: "Email already taken",
            type: "error",
          });
        } else if (res.status === 400) {
          this.$toast.open({
            message: "Invalid Captcha",
            type: "error",
          });
          // Optionally, you can reload the captcha if the user typed it wrong
          this.reloadCaptcha();
        } else {
          this.$toast.open({
            message: "Successfully registered, verify your email to login.",
            type: "success",
          });
          this.$router.push("/");
        }
      } catch (err) {
        this.$toast.open({
          message: "Server error. Please try again.",
          type: "error",
        });
        console.error("Error while registering:", err);
      }
    },
  },
  created() {
    this.setRandomImage();
    // Immediately load the first captcha on page creation
    this.reloadCaptcha();
  },
};
</script>

<style>
.bg-forest-reg {
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

.register__wrapper {
  display: flex;
  background-color: var(--bg-neutral);
  border-radius: 20px;
  color: var(--color-white);
  box-shadow: 0 4px 15px rgb(0, 0, 0);
  overflow: hidden;
  align-items: center;
  width: 85vw;
  max-width: 900px;
}

.image-div-reg {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 50%;
  height: 100%;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.register {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 25px;
  margin: 0 auto;
  padding: 5%;
  width: 70%;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 20px;
  color: var(--color-white);
  width: 100%;
}

.form-row {
  display: flex;
  gap: 20px;
  width: 100%;
}

.form-input {
  flex: 1;
  color: var(--color-white);
}

textarea {
  resize: none;
  width: 100%;
}

.register button {
  width: fit-content;
  text-align: center;
}

.button-or-signIn {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

#captchaRow{
  display: flex;
  gap: 10px;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}

/* Media Queries */

/* Tablet and Phone View (768px and below) */
@media (max-width: 900px) {
  html,
  body {
    overflow-y: hidden; /* Prevent horizontal scrolling */
  }

  .bg-forest-reg {
    filter: blur(8px) brightness(60%);
    top: 0;
  }

  .register__wrapper {
    flex-direction: column;
    align-items: center;
    width: 95%;
    max-width: 600px;
    margin: 0 auto;
    border-radius: 10px;
  }

  .image-div-reg {
    width: 100%;
    height: 200px;
    display: block;
  }

  .register {
    width: 90%;
    padding: 5%;
  }

  h1.mori {
    font-size: 1.8rem;
  }

  .register button {
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
  .register__wrapper {
    flex-direction: column;
    width: 95%;
    margin: 0 auto;
  }

  .image-div-reg {
    width: 100%;
    height: 150px; /* Slightly smaller height for phone screens */
    display: block;
  }

  .register {
    width: 100%;
    padding: 4%;
  }

  h1.mori {
    font-size: 1.5rem;
  }

  .bg-forest-reg {
    filter: blur(8px) brightness(60%);
    top: 0;
  }

  .register button {
    display: flex;
    width: 100%; /* Wider button for better touch interaction */
    align-items: center;
    justify-content: center;
    text-align: center;
    font-size: 15px;
  }

  p {
    font-size: 14px;
  }
}
</style>
