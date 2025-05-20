<template>
  <div class="forgot-password-page">
    <div class="forgot-password-container">
      <h1 class="mori">
        Mori <span class="adder">- Reset Password</span>
      </h1>
      <h2>Forgot Your Password ?</h2>
      <p class="textfp">Please enter your email to receive a reset link.</p>

      <div class="form-row">
        <input
          v-model="email"
          type="email"
          placeholder="Enter your email"
          class="email-input"
        />
        <button class="reset-button" @click="requestResetPassword">
          Send Link
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "ForgotPassword",
  data() {
    return {
      email: "",
    };
  },
  methods: {
    async requestResetPassword() {
      try {
        const response = await fetch("http://localhost:8081/request-password-reset", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ email: this.email }),
        });
        const data = await response.json();

        if (response.ok) {
          // Toast success message
          this.$toast.open({
            message: "A password reset link has been sent to your email!",
            type: "success",
          });
          // Redirect to sign-in page
          this.$router.push("/sign-in");
        } else {
          // Show toast for the error
          this.$toast.open({
            message: data.message || "Could not send reset link. Please try again.",
            type: "error",
          });
        }
      } catch (err) {
        console.error("Network error or server unreachable", err);
        this.$toast.open({
          message: "Network error. Please try again later.",
          type: "error",
        });
      }
    },
  },
};
</script>

<style scoped>
.mori {
  text-align: left;
  margin-bottom: 20px;
}

h2 {
  text-align: left;
  font-weight: bolder;
  color: var(--color-white);
  font-size: 20px;
}

.textfp {
  text-align: left;
  color: var(--purple-color);
  font-size: 17px;
}

.forgot-password-page {
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--page-bg);
}

.forgot-password-container {
  display: flex;
  gap: 15px;
  flex-direction: column;
  width: 600px;
  margin: 50px auto;
  text-align: center;
  font-family: Arial, sans-serif;
  background-color: var(--bg-neutral);
  border-radius: 20px;
  box-shadow: 0 4px 15px rgb(0, 0, 0);
  padding: 30px;
}

.form-row {
  display: flex;
  flex-direction: row;
  gap: 10px;
  margin-top: 20px;
}

.email-input {
  padding: 10px;
  font-size: 16px;
  flex: 1; /* so the input expands, button stays fixed width */
}

.reset-button {
  width: 150px;
  padding: 10px;
  font-size: 16px;
  background-color: #9146bc;
  color: #fff;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.reset-button:hover {
  background-color: #7d3aa6;
}
</style>
