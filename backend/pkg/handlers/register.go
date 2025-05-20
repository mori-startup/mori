package handlers

import (
	"net/http"

	"mori/pkg/models"
	"mori/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register user endpoint -> validate inputs / save in db / start session
func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	if r.Method != http.MethodPost {
		utils.RespondWithError(w, "Error on form submission", 400)
		return
	}
	// Parse up to 3MB of multipart form data
	err := r.ParseMultipartForm(3145728) // 3MB
	if err != nil {
		utils.RespondWithError(w, "Error in form validation", 400)
		return
	}

	// 1) Get captcha cookie & user input from form
	captchaCookie, err := r.Cookie("captcha_id")
	if err != nil {
		utils.RespondWithError(w, "captcha cookie missing", 400)
		return
	}
	captchaID := captchaCookie.Value
	captchaValue := r.PostFormValue("captchaValue")

	// 2) Create new user instance from form values
	newUser := models.User{
		Email:       r.PostFormValue("email"),
		FirstName:   r.PostFormValue("firstname"),
		LastName:    r.PostFormValue("lastname"),
		Password:    r.PostFormValue("password"),
		Nickname:    r.PostFormValue("nickname"),
		About:       r.PostFormValue("aboutme"),
		DateOfBirth: r.PostFormValue("dateofbirth"),
	}

	// 3) Validate user fields AND captcha together
	errValid := utils.ValidateNewUser(newUser, captchaID, captchaValue)
	if errValid != nil {
		utils.RespondWithError(w, errValid.Error(), 400)
		return
	}

	// 4) Check if email is already taken
	emailUnique, _ := handler.repos.UserRepo.EmailNotTaken(newUser.Email)
	if !emailUnique {
		utils.RespondWithError(w, "Email already taken", 409)
		return
	}

	// 5) Hash the user password
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	newUser.Password = string(hashedPwd)

	// 6) Create a unique user ID
	userID := utils.UniqueId()
	newUser.ID = userID

	// 7) If the user uploaded an avatar, save it
	newUser.ImagePath = utils.SaveAvatar(r)

	// 8) Save the user in the DB (and automatically send email)
	errSave := handler.repos.UserRepo.Add(newUser)
	if errSave != nil {
		utils.RespondWithError(w, "Couldn't save new user or send email", 500)
		return
	}

	// 9) Optionally start a session
	/*
	   session := utils.SessionStart(w, r, userID)
	   errSession := handler.repos.SessionRepo.Set(session)
	   if errSession != nil {
	       utils.RespondWithError(w, "Error on creating new session", 500)
	       return
	   }
	*/
}
