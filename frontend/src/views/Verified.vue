<template>
  <div class="container">
    <div class="card">
      <div class="logo">
        <!-- Use your Mori logo or any image you want -->
        <img src="../assets/mori.png" alt="Logo">
      </div>

      <!-- Step 1: "verifying" state -->
      <template v-if="verificationStatus === 'verifying'">
        <h1>Verifying your email ...</h1>
        <p>Please wait a moment.</p>
      </template>

      <!-- Step 2: "success" state -->
      <template v-else-if="verificationStatus === 'success'">
        <h1>Email Verified! 🎉</h1>
        <p>Your email has been successfully verified. You can now log in to your account.</p>
        <button class="btn" @click="goToLogin">Go to Login</button>
      </template>

      <!-- Step 3: "failure" state -->
      <template v-else-if="verificationStatus === 'failure'">
        <h1>Invalid or Expired Token ❌</h1>
        <p>Please try again or contact support.</p>
        <button class="btn" @click="goToLogin">Go to Login</button>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  name: "Verified",
  data() {
    return {
      verificationStatus: "verifying",
    };
  },
  created() {
    // Grab the token from the URL: /verified?token=xxx
    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get("token");

    if (!token) {
      // No token -> mark as failure
      this.verificationStatus = "failure";
    } else {
      // If we have a token, attempt to verify on the server
      this.verifyEmail(token);
    }
  },
  methods: {
    async verifyEmail(token) {
      try {
        const res = await fetch(`http://localhost:8081/verified?token=${token}`, {
          method: "GET",
          credentials: "include",
        });

        // Attempt to parse JSON (your server returns JSON on success/failure)
        const data = await res.json().catch(() => null);

        switch (res.status) {
          case 200:
            console.log("✅ Verification succeeded:", data);
            this.verificationStatus = "success";
            break;

          case 400:
            console.warn("❌ Verification failed (400):", data);
            this.verificationStatus = "failure";
            break;

          default:
            console.error(`Unexpected status ${res.status}:`, data);
            this.verificationStatus = "failure";
        }
      } catch (error) {
        console.error("❌ Network or fetch error:", error);
        this.verificationStatus = "failure";
      }
    },

    goToLogin() {
      this.$router.push("/sign-in");
    },
  },
};
</script>

<style scoped>
/* Body styling is now handled by the container/card, so we don't 
   do a direct "body" { ... } reset as we did in the old HTML. 
   But we can replicate the same background color by styling .container. */

/* General container styling */
.container {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  text-align: center;
  animation: fadeIn 1.5s ease-in-out;
  min-height: 100vh;
  background-color: #1a1a1a; /* Dark background */
  color: #e4e4e4; /* Light text */
  margin: 0; /* reset any default margin */
  padding: 0; /* reset any default padding */
}

/* Card styling */
.card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #2c2c2c; /* Slightly lighter dark background */
  padding: 40px;
  border-radius: 15px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.6);
  max-width: 400px;
  width: 100%;
  text-align: center;
  animation: popUp 0.5s ease;
  transition: all 0.3s ease;
}

.card:hover {
  transform: translateY(-5px) scale(1.01); /* Lift effect and slight scale */
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.8);
}

/* Logo styling */
.logo img {
  width: 100px;
  height: 100px;
  margin-bottom: 10px;
  border-radius: 50%;
}

/* Heading styling */
.card h1 {
  font-size: 28px;
  margin: 10px 0;
  font-weight: bold;
  color: #ffffff;
}

/* Paragraph styling */
.card p {
  font-size: 16px;
  margin-bottom: 40px;
  color: #bfbfbf;
  line-height: 1.5;
}

/* Button styling */
.card .btn {
  text-decoration: none;
  background: #9146bc;
  color: #ffffff;
  padding: 10px 20px;
  border-radius: 5px;
  font-weight: bold;
  transition: background 0.3s ease;
  cursor: pointer;
  border: none;
}

.card .btn:hover {
  background: #7d3aa6;
}

/* Animations */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes popUp {
  from {
    transform: scale(0.8);
  }
  to {
    transform: scale(1);
  }
}
</style>
