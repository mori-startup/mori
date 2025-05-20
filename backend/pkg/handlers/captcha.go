package handlers

import (
	"net/http"

	"github.com/dchest/captcha"
)

// ServeCaptcha génère et sert une nouvelle captcha
func (handler *Handler) ServeCaptcha(w http.ResponseWriter, r *http.Request) {
	// Génère un nouvel ID de captcha
	captchaID := captcha.New()

	// Retourne l'ID au frontend
	http.SetCookie(w, &http.Cookie{
		Name:     "captcha_id",
		Value:    captchaID,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
	captcha.WriteImage(w, captchaID, 240, 80) // Image de taille 240x80 pixels
}

// VerifyCaptcha vérifie si l'utilisateur a correctement rempli le captcha
func VerifyCaptcha(captchaID, userInput string) bool {
	return captcha.VerifyString(captchaID, userInput)
}
