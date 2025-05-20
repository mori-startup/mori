package handlers

import (
	"encoding/json"
	"mori/pkg/models"
	"mori/pkg/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (handler *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	if r.Method != http.MethodPost {
		// We can still return a 400 or 405, but to keep your current logic:
		utils.RespondWithError(w, "Error on form submission", 200)
		return
	}

	// 1) Decode incoming JSON
	var client models.User
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		utils.RespondWithError(w, "Error on form submission", 200)
		return
	}

	// 2) Find user by email
	dbUser, errDb := handler.repos.UserRepo.FindUserByEmail(client.Email)
	if errDb != nil {
		// Email not found or DB error
		utils.RespondWithError(w, "Wrong credentials", 200)
		return
	}

	// 4) Compare password
	errPwd := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(client.Password))
	if errPwd != nil {
		utils.RespondWithError(w, "Wrong credentials", 200)
		return
	}

	// 3) Check if user is verified
	if !dbUser.Verified {
		utils.RespondWithError(w, "Please verify your email to sign in", 200)
		return
	}

	// 5) User valid -> create/update session
	_, errSession := handler.repos.SessionRepo.GetByUser(dbUser.ID)
	newSession := utils.SessionStart(w, r, dbUser.ID)
	var errOnSave error

	// 6) Update or create new row in db based on existing session
	if errSession != nil {
		// create new session
		errOnSave = handler.repos.SessionRepo.Set(newSession)
	} else {
		// update existing session
		errOnSave = handler.repos.SessionRepo.Update(newSession)
	}

	if errOnSave != nil {
		utils.RespondWithError(w, "Error on creating new session", 200)
		return
	}

	// 7) Respond success
	utils.RespondWithSuccess(w, "Login successful", 200)
}

// SessionActive remains unchanged...
func (handler *Handler) SessionActive(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	sessionId, errCookie := utils.GetCookie(r)
	if errCookie != nil {
		utils.RespondWithError(w, "Session not active", 200)
		return
	}

	session, errSession := handler.repos.SessionRepo.Get(sessionId)
	if errSession != nil {
		utils.RespondWithError(w, "Session not active", 200)
		return
	}

	sessionValid := utils.CheckSessionExpiration(session)
	if !sessionValid {
		handler.repos.SessionRepo.Delete(session)
		utils.DeleteCookie(w)
		utils.RespondWithError(w, "Session not active", 200)
		return
	}
	// Prolong session
	session.ExpirationTime = time.Now().Add(30 * time.Minute)
	handler.repos.SessionRepo.Update(session)
	utils.RespondWithSuccess(w, "Session active", 200)
}
