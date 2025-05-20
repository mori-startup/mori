<template>
    <div class="reset-password-page">
      <div class="reset-password-container">
        <h1 class="mori">Mori <span class="adder">- New Password</span></h1>
        <h2>Reset Your Password</h2>
  
        <!-- If we are in "loading" or "error" states, we can show a message -->
        <div v-if="status === 'verifying'">
          <p>Checking your token. Please wait...</p>
        </div>
  
        <!-- The actual form for a new password -->
        <div v-else-if="status === 'form'">
          <p>Please enter a new password below and confirm it.</p>
  
          <!-- First password input -->
          <div class="form-row">
            <input
              v-model="newPassword"
              type="password"
              placeholder="Enter a new password"
              class="password-input"
            />
          </div>
  
          <!-- Confirm password input -->
          <div class="form-row">
            <input
              v-model="confirmPassword"
              type="password"
              placeholder="Confirm new password"
              class="password-input"
            />
          </div>
  
          <div class="form-row">
            <button class="submit-button" @click="submitNewPassword">
              Reset Password
            </button>
          </div>
        </div>
  
        <!-- If success -->
        <div class="success" v-else-if="status === 'success'">
          <p>Your password has been updated. You can now sign in!</p>
          <button class="submit-button" @click="goToLogin">Go to Login</button>
        </div>
  
        <!-- If failure -->
        <div v-else-if="status === 'failure'">
          <p>{{ errorMessage }}</p>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: "ResetPassword",
    data() {
      return {
        token: "",
        newPassword: "",
        confirmPassword: "",
        status: "form", // can be "verifying", "form", "success", or "failure"
        errorMessage: "",
      };
    },
    created() {
      // Grab token from the URL (e.g., /reset-password?token=XYZ)
      const urlParams = new URLSearchParams(window.location.search);
      this.token = urlParams.get("token");
  
      if (!this.token) {
        // No token => can't proceed
        this.status = "failure";
        this.errorMessage = "Invalid link. No token found.";
      }
    },
    methods: {
      async submitNewPassword() {
        // Basic validation: ensure both fields are filled
        if (!this.newPassword.trim() || !this.confirmPassword.trim()) {
          this.$toast.open({
            message: "Please fill both password fields.",
            type: "error",
          });
          return;
        }
  
        // Check if passwords match
        if (this.newPassword !== this.confirmPassword) {
          this.$toast.open({
            message: "Passwords do not match. Please try again.",
            type: "error",
          });
          return;
        }
  
        // 1) Check min length: 10
        if (this.newPassword.length < 10) {
          this.$toast.open({
            message: "Password must be at least 10 characters.",
            type: "error",
          });
          return;
        }
  
        // 2) Check at least 1 uppercase
        if (!/[A-Z]/.test(this.newPassword)) {
          this.$toast.open({
            message: "Password must contain at least one uppercase letter.",
            type: "error",
          });
          return;
        }
  
        // 3) Check at least 1 digit
        if (!/\d/.test(this.newPassword)) {
          this.$toast.open({
            message: "Password must contain at least one digit.",
            type: "error",
          });
          return;
        }
  
        // 4) Check at least 1 special character
        if (!/[^a-zA-Z0-9]/.test(this.newPassword)) {
          this.$toast.open({
            message: "Password must contain at least one special character.",
            type: "error",
          });
          return;
        }
  
        // If all checks pass, send request to back-end
        try {
          const response = await fetch("http://localhost:8081/reset-password", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              token: this.token,
              newPassword: this.newPassword,
            }),
          });
          const data = await response.json();
  
          if (response.ok) {
            // success
            this.$toast.open({
              message: "Password reset successful!",
              type: "success",
            });
            this.status = "success";
          } else {
            // error
            this.status = "failure";
            this.errorMessage = data.message || "Could not reset password.";
            this.$toast.open({
              message: this.errorMessage,
              type: "error",
            });
          }
        } catch (err) {
          console.error("Network or server error:", err);
          this.status = "failure";
          this.errorMessage = "Network error. Please try again later.";
          this.$toast.open({
            message: this.errorMessage,
            type: "error",
          });
        }
      },
      goToLogin() {
        this.$router.push("/sign-in");
      },
    },
  };
  </script>
  
  <style scoped>
  .mori {
    text-align: left;
    margin-bottom: 20px;
  }
  .reset-password-page {
    width: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: var(--page-bg);
  }
  .reset-password-container {
    display: flex;
    gap: 15px;
    flex-direction: column;
    width: 500px;
    margin: 50px auto;
    text-align: center;
    font-family: Arial, sans-serif;
    background-color: var(--bg-neutral);
    border-radius: 20px;
    box-shadow: 0 4px 15px rgb(0, 0, 0);
    padding: 30px;
  }
  
  h2 {
    text-align: left;
    font-weight: bolder;
    color: var(--color-white);
    font-size: 20px;
  }
  
  p {
    text-align: left;
    color: var(--purple-color);
    font-size: 17px;
  }
  
  .form-row {
    display: flex;
    flex-direction: row;
    gap: 10px;
    margin-top: 20px;
  }
  
  .password-input {
    width: 90%;
    padding: 10px;
    font-size: 16px;
  }
  
  .submit-button {
    width: 160px;
    padding: 10px;
    font-size: 16px;
    background-color: #9146bc;
    color: #fff;
    border: none;
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.3s ease;
  }
  
  .submit-button:hover {
    background-color: #7d3aa6;
  }

  .success{
    gap: 25px;
    display: flex;
    flex-direction: column;
  }
  </style>
  