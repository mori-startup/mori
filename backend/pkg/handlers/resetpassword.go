package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"mori/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type ResetRequest struct {
	Email string `json:"email"`
}

// generateRandomToken creates a 32-char hex string
func generateRandomToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// RequestPasswordReset handles POST /request-password-reset
func (handler *Handler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	// 1) Set the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// 2) If it’s an OPTIONS preflight request, respond with allowed methods/headers
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req ResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 1) Find user by email
	user, errFind := handler.repos.UserRepo.FindUserByEmail(req.Email)
	if errFind != nil {
		// If user not found, you could respond success anyway (for security),
		// or respond 404. It's up to you:
		utils.RespondWithError(w, "Email not found", http.StatusNotFound)
		return
	}

	// 2) Generate reset token
	token, errToken := generateRandomToken()
	if errToken != nil {
		utils.RespondWithError(w, "Couldn't generate token", http.StatusInternalServerError)
		return
	}

	// 3) Set token & expiration in DB
	expires := time.Now().Add(1 * time.Hour)
	errSet := handler.repos.UserRepo.SetResetToken(user.ID, token, expires)
	if errSet != nil {
		utils.RespondWithError(w, "Couldn't update reset token in DB", http.StatusInternalServerError)
		return
	}

	fmt.Println(req.Email)

	errEmail := sendResetEmail(req.Email, token)
	if errEmail != nil {
		utils.RespondWithError(w, "Couldn't send email", http.StatusInternalServerError)
		return
	}

	// 5) Return success
	utils.RespondWithSuccess(w, "Reset email sent", http.StatusOK)
}

// Example: sendResetEmail function

func sendResetEmail(toEmail, token string) error {
	// SMTP config
	from := "mori.team.contact@gmail.com"
	password := "qeey kngz gmyn bzwi" // Put this in env variables for production
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Email subject
	subject := "Reset Your Mori Password"

	// Build the reset link
	resetLink := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", token)

	// HTML email body (similar style to verification email)
	// We'll greet with 'Hello toEmail'
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            padding: 20px;
            background-color: #2c2c2c;
            border-radius: 10px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            color: #e4e4e4;
        }
        .header {
            text-align: center;
            padding: 10px 0;
        }
        .header h1 {
            font-size: 24px;
            margin: 0;
            color: #9146bc;
        }
        .content {
            margin: 20px 0;
            font-size: 16px;
            line-height: 1.6;
        }
        .content p {
            color: #e4e4e4 !important;
            font-size: 16px !important;
        }
        .button {
            display: inline-block;
            padding: 12px 20px;
            margin: 20px 0;
            font-size: 16px;
            color: #fff !important;
            background-color: #9146bc;
            border-radius: 5px;
            text-decoration: none;
        }
        .footer {
            margin-top: 80px;
            text-align: center;
            font-size: 12px;
            color: #aaa;
        }
        .footer a {
            color: #9146bc;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Mori Team</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            <p>We received a request to reset your Mori account password.
            If you didn't request this, you can safely ignore this message.</p>
            <p>To reset your password, please click the button below:</p>
            <p>
                <a href="%s" class="button">Reset Password</a>
            </p>
            <p>If the button doesn't work, please copy and paste the following link into your browser:</p>
            <p>%s</p>
            <p>Thank you, and take care!</p>
        </div>
        <div class="footer">
            <p>&copy; 2024 Mori Team. All Rights Reserved.</p>
            <p>
                <a href="#">Privacy Policy</a> | <a href="#">Contact Us</a>
            </p>
        </div>
    </div>
</body>
</html>`,
		toEmail, resetLink, resetLink)

	// Build MIME headers + body
	msg := strings.Join([]string{
		fmt.Sprintf("Subject: %s", subject),
		"MIME-version: 1.0;",
		`Content-Type: text/html; charset="UTF-8";`,
		"",
		body,
	}, "\n")

	// Auth for Gmail (App password recommended)
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email
	toList := []string{toEmail}
	errSend := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toList, []byte(msg))
	if errSend != nil {
		return fmt.Errorf("failed to send reset email: %w", errSend)
	}

	return nil
}

type ResetPasswordPayload struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

func (handler *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	// 1) Set the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// 2) If it’s an OPTIONS preflight request, respond with allowed methods/headers
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		utils.RespondWithError(w, "Invalid method", http.StatusBadRequest)
		return
	}

	// 3) Parse the JSON from body
	var req ResetPasswordPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 4) Find user by the token
	user, errUser := handler.repos.UserRepo.FindUserByResetToken(req.Token)
	if errUser != nil {
		utils.RespondWithError(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	// 5) Check if token is expired
	if user.ResetTokenExpires == nil || user.ResetTokenExpires.Before(time.Now()) {
		utils.RespondWithError(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	// 6) Hash new password
	hashedPwd, errHash := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if errHash != nil {
		utils.RespondWithError(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	// 7) Update DB: set new password, clear token
	errUpdate := handler.repos.UserRepo.UpdatePasswordAndClearToken(user.ID, string(hashedPwd))
	if errUpdate != nil {
		utils.RespondWithError(w, "Could not update password", http.StatusInternalServerError)
		return
	}

	// 8) Respond success
	utils.RespondWithSuccess(w, "Password reset successful", http.StatusOK)
}
