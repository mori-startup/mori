// JavaScript to toggle between login and register forms
function toggleForms() {
  const loginForm = document.getElementById("login-form");
  const registerForm = document.getElementById("register-form");

  if (loginForm.style.display === "none") {
    loginForm.style.display = "flex";
    registerForm.style.display = "none";
  } else {
    loginForm.style.display = "none";
    registerForm.style.display = "flex";
  }
}

// Function to handle login form submission
document.getElementById("login-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const errorDiv = document.getElementById("login-error");
  errorDiv.style.display = "none";

  const formData = new FormData(e.target);
  const response = await fetch("/login", {
    method: "POST",
    body: formData,
  });

  const result = await response.json();
  if (!response.ok) {
    errorDiv.textContent = result.error || "An error occurred during login.";
    errorDiv.style.display = "block";
  } else {
    window.location.href = "/dashboard"; // Redirect to dashboard on success
  }
});

function showPopup(message) {
  const popup = document.getElementById("popup");
  const popupMessage = document.getElementById("popup-message");

  popupMessage.textContent = message;
  popup.classList.remove("hidden");
  popup.classList.add("show");

  // Hide the popup after 3 seconds
  setTimeout(() => {
    popup.classList.remove("show");
    popup.classList.add("hidden");
  }, 3000);
}

document
  .getElementById("register-form")
  .addEventListener("submit", async (e) => {
    e.preventDefault();
    const errorDiv = document.getElementById("register-error");
    errorDiv.style.display = "none";

    const formData = new FormData(e.target);
    const response = await fetch("/register", {
      method: "POST",
      body: formData,
    });

    const result = await response.json();
    if (!response.ok) {
      errorDiv.textContent =
        result.error || "An error occurred during registration.";
      errorDiv.style.display = "block";
    } else {
      errorDiv.style.display = "none";
      errorDiv.textContent = "";
      showPopup("Registration successful, please verify your email.");
      toggleForms(); // Switch to login form
    }
  });

  // Function to reload captcha
function reloadCaptcha() {
    const captchaImage = document.getElementById("captcha-image");
    captchaImage.src = "/captcha?" + new Date().getTime(); // Append timestamp to prevent caching
  }
  
  // Handle registration form submission
  document
  .getElementById("register-form")
  .addEventListener("submit", async (e) => {
    e.preventDefault();
    const errorDiv = document.getElementById("register-error");
    const submitButton = e.target.querySelector("button[type='submit']");
    errorDiv.style.display = "none";

    // Disable the button to prevent multiple submissions
    submitButton.disabled = true;

    const formData = new FormData(e.target);
    const response = await fetch("/register", {
      method: "POST",
      body: formData,
    });

    const result = await response.json();
    if (!response.ok) {
      // Enable the button again if there's an error
      submitButton.disabled = false;

      errorDiv.textContent =
        result.error || "An error occurred during registration.";
      errorDiv.style.display = "block";

      if (result.reloadCaptcha) {
        reloadCaptcha();
      }
    } else {
      errorDiv.style.display = "none";
      errorDiv.textContent = "";
      showPopup("Registration successful, please Login.");
      toggleForms(); // Switch to login form
    }
  });

