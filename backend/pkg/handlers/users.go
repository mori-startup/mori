package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"mori/pkg/models"
	"mori/pkg/utils"
	ws "mori/pkg/wsServer"
)

/* -------------------------------------------------------------------------- */
/*                                    users                                   */
/* -------------------------------------------------------------------------- */
// Find all users and they relation with current user
func (handler *Handler) AllUsers(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// request all users exccept current + relations
	users, errUsers := handler.repos.UserRepo.GetAllAndFollowing(userId)
	if errUsers != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, users, 200)
}

// Returns user nickname, id and path to avatar
func (handler *Handler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	user, err := handler.repos.UserRepo.GetDataMin(userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, []models.User{user}, 200)
}

// Returns user data based on public / private profile and user_id from request
// waits for GET request with query "userId" ->user client is looking for
//
//	can be used both on own profile and other users
func (handler *Handler) UserData(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	// Access user ID from the request context
	currentUserId := r.Context().Value(utils.UserKey).(string)

	// Get the userId from the query
	query := r.URL.Query()
	userId := query.Get("userId")

	// Check profile visibility status
	status, err := handler.repos.UserRepo.ProfileStatus(userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting profile status", 200)
		return
	}

	// Determine if the current user is accessing their own profile
	currentUser := (currentUserId == userId)

	var following bool
	if !currentUser {
		// Check if the current user is following the profile user
		following, err = handler.repos.UserRepo.IsFollowing(userId, currentUserId)
		if err != nil {
			utils.RespondWithError(w, "Error on checking following status", 200)
			return
		}
	}

	// Fetch user details based on their profile visibility
	var user models.User
	if currentUser || following || status == "PUBLIC" {
		// Fetch full profile details including avatar and name
		user, err = handler.repos.UserRepo.GetProfileMax(userId)
	} else {
		// Fetch minimal profile details
		user, err = handler.repos.UserRepo.GetProfileMin(userId)
		if err == nil {
			// Check if a follow request is pending
			notif := models.Notification{Type: "FOLLOW", Content: currentUserId, TargetID: userId}
			user.FollowRequestPending, err = handler.repos.NotifRepo.CheckIfExists(notif)
		}
	}

	if err != nil {
		utils.RespondWithError(w, "Error on getting profile data", 200)
		return
	}

	// Ensure avatar and name are included in the response
	user.Following = following
	user.CurrentUser = currentUser
	user.Status = status

	// Respond with the user data
	utils.RespondWithUsers(w, []models.User{user}, 200)
}

// changes user status in db return status
// in case of turning to PUBLIC -> also accept follow requests
func (handler *Handler) UserStatus(w http.ResponseWriter, r *http.Request) {
	statusList := []string{"PUBLIC", "PRIVATE"} // possible status
	var client models.User

	w = utils.ConfigHeader(w)
	// access user id
	client.ID = r.Context().Value(utils.UserKey).(string)
	// get status from request
	query := r.URL.Query()
	reqStatus := strings.ToUpper(query.Get("status"))

	// check if valid value and asign to user
	if reqStatus == statusList[0] {
		client.Status = statusList[0]
	} else if reqStatus == statusList[1] {
		client.Status = statusList[1]
	} else {
		utils.RespondWithError(w, "Requested status not valid", 200)
		return
	}
	// request current status from db
	currentStatus, err := handler.repos.UserRepo.GetStatus(client.ID)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// check if requested status is not the same as current
	if currentStatus == client.Status {
		utils.RespondWithError(w, "Status change not valid", 200)
		return
	}
	// Set new status
	err = handler.repos.UserRepo.SetStatus(client)
	if err != nil {
		utils.RespondWithError(w, "Error on saving status", 200)
		return
	}
	// if new status is public -> also accept pending follow requests
	// responds with success and newly created status
	utils.RespondWithSuccess(w, client.Status, 200)
}

/* -------------------------------------------------------------------------- */
/*                                  followers                                 */
/* -------------------------------------------------------------------------- */
// Find all followers
func (handler *Handler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// get userId from request
	query := r.URL.Query()
	userId := query.Get("userId")
	// request all  following users
	followers, errUsers := handler.repos.UserRepo.GetFollowers(userId)
	if errUsers != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, followers, 200)
}

// Find all who clinet is following
func (handler *Handler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// get userId from request
	query := r.URL.Query()
	userId := query.Get("userId")
	// request all  following users
	followers, errUsers := handler.repos.UserRepo.GetFollowing(userId)
	if errUsers != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, followers, 200)
}

func (handler *Handler) Follow(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	currentUserId := r.Context().Value(utils.UserKey).(string)
	// get status from request
	query := r.URL.Query()
	reqUserId := query.Get("userId")
	/* ----------------- safety check -> if request already made ---------------- */
	alreadyFollowing, _ := handler.repos.UserRepo.IsFollowing(reqUserId, currentUserId)
	if alreadyFollowing {
		utils.RespondWithError(w, "User already is following", 200)
		return
	}
	// get target user profile status -> public or private
	reqUserStatus, err := handler.repos.UserRepo.GetStatus(reqUserId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if reqUserStatus == "PUBLIC" {
		// SAVE AS FOLLOWER
		err := handler.repos.UserRepo.SaveFollower(reqUserId, currentUserId)
		if err != nil {
			utils.RespondWithError(w, "Error on saving follower", 200)
			return
		}
	} else if reqUserStatus == "PRIVATE" {
		// SAVE IN NOTIFICATIONS as pending folllow request
		notification := models.Notification{
			ID:       utils.UniqueId(),
			TargetID: reqUserId,
			Type:     "FOLLOW",
			Content:  currentUserId,
			Sender:   currentUserId,
		}
		err := handler.repos.NotifRepo.Save(notification)
		if err != nil {
			utils.RespondWithError(w, "Error on save", 200)
			return
		}
		// if user online send notification about follow request
		for client := range wsServer.Clients {
			if client.ID == reqUserId {
				client.SendNotification(notification)
			}
		}

	}
	utils.RespondWithSuccess(w, "Following successful", 200)
}

func (handler *Handler) CancelFollowRequest(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	currentUserId := r.Context().Value(utils.UserKey).(string)
	// get status from request
	query := r.URL.Query()
	reqUserId := query.Get("userId")
	// delete notification corresponding to follow request
	notif := models.Notification{
		Type:     "FOLLOW",
		TargetID: reqUserId,
		Content:  currentUserId,
	}
	if err := handler.repos.NotifRepo.DeleteByType(notif); err != nil {
		utils.RespondWithError(w, "Error on canceling request", 200)
		return
	}
	utils.RespondWithSuccess(w, "Follow request canceled successfuly", 200)
}

func (handler *Handler) Unfollow(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	currentUserId := r.Context().Value(utils.UserKey).(string)
	// get status from request
	query := r.URL.Query()
	reqUserId := query.Get("userId")

	if err := handler.repos.UserRepo.DeleteFollower(reqUserId, currentUserId); err != nil {
		utils.RespondWithError(w, "Error on deleting follower", 200)
		return
	}
	utils.RespondWithSuccess(w, "Unfollowing successful", 200)
}

// not tested
// wait for POST request with notification Id and response -"ACCEPT" or "DECLINE"
func (handler *Handler) ResponseFollowRequest(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	// Try to decode the JSON request to a new response
	type Response struct {
		RequestID string `json:"requestId"`
		Response  string `json:"response"` // ACCEPT or DECLINE
	}
	var resp Response
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	// get other user id from notification
	followerId, err := handler.repos.NotifRepo.GetUserFromRequest(resp.RequestID)
	userId := r.Context().Value(utils.UserKey).(string)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	if strings.ToUpper(resp.Response) == "ACCEPT" {
		err = handler.repos.UserRepo.SaveFollower(userId, followerId)
		if err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
	}
	/* ----------------------- delete pending notification ---------------------- */
	err = handler.repos.NotifRepo.Delete(resp.RequestID)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	// notify websocket about notification changes
	utils.RespondWithSuccess(w, "Response successful", 200)
}

/* -------------------------------------------------------------------------- */
/*                                  chat List                                 */
/* -------------------------------------------------------------------------- */
func (handler *Handler) ChatList(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	// Récupération du userId depuis la query
	userId := r.URL.Query().Get("userId")

	// 1) Récupérer les "amis" ou "following"
	followers, errUsers := handler.repos.UserRepo.GetFollowing(userId)
	if errUsers != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}

	// 2) Récupérer l'historique DM (seulement type='PERSON') => un set d'IDs
	ids, errIds := handler.repos.MsgRepo.GetChatHistoryIds(userId)
	if errIds != nil {
		utils.RespondWithError(w, "Error on getting chat history", 200)
		return
	}

	// 3) Parcourir ces IDs
	for currentId := range ids {
		// Si on a déjà ce user dans "followers", on ne fait rien
		isPresent := ContainsUser(followers, currentId)
		if isPresent {
			continue
		}

		// Tenter de récupérer l'utilisateur (table "users")
		user, err := handler.repos.UserRepo.GetDataMin(currentId)
		if err != nil {
			// Si erreur => peut-être un group_id ou user inexistant => on skip
			// log.Println("Skipping ID", currentId, ":", err)
			continue
		}
		// Ajouter ce user à la liste
		followers = append(followers, user)
	}

	// 4) Répondre au client avec la liste d'utilisateurs
	utils.RespondWithUsers(w, followers, 200)
}

/* --------------------------------- helper --------------------------------- */
func ContainsUser(list []models.User, id string) bool {
	for _, value := range list {
		if value.ID == id {
			return true
		}
	}
	return false
}

/* -------------------------------------------------------------------------- */

// ChangeNickname allows the user to update their nickname.
func (handler *Handler) ChangeNickname(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	log.Println("Début de ChangeNickname")
	// Retrieve the user ID from the request context
	userId, ok := r.Context().Value(utils.UserKey).(string)
	if !ok || userId == "" {
		log.Println("Erreur : utilisateur non authentifié")
		utils.RespondWithError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Read request body
	type RequestBody struct {
		Nickname string `json:"nickname"`
	}
	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.Nickname == "" {
		log.Printf("Erreur de décodage ou nickname invalide : %v", err)
		utils.RespondWithError(w, "Invalid nickname provided", http.StatusBadRequest)
		return
	}
	log.Printf("Tentative de mise à jour du nickname pour l'utilisateur %s avec le nickname %s", userId, body.Nickname)
	// Update nickname in the database
	err = handler.repos.UserRepo.UpdateNickname(userId, body.Nickname)
	if err != nil {
		log.Printf("Erreur lors de la mise à jour du nickname dans la base de données : %v", err)
		utils.RespondWithError(w, "Error updating nickname", http.StatusInternalServerError)
		return
	}
	log.Println("Nickname mis à jour avec succès")
	utils.RespondWithSuccess(w, "Nickname updated successfully", http.StatusOK)
}

// ChangeAvatar allows the user to update their avatar.
func (handler *Handler) ChangeAvatar(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	userId, ok := r.Context().Value(utils.UserKey).(string)
	if !ok || userId == "" {
		log.Println("Utilisateur non authentifié")
		utils.RespondWithError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	log.Printf("Tentative de mise à jour de l'avatar pour l'utilisateur %s", userId)
	avatarPath := utils.SaveAvatar(r)
	if avatarPath == "" {
		log.Println("Aucun fichier valide reçu")
		utils.RespondWithError(w, "Invalid avatar file", http.StatusBadRequest)
		return
	}

	log.Printf("Chemin de l'avatar sauvegardé : %s", avatarPath)
	err := handler.repos.UserRepo.UpdateAvatar(userId, avatarPath)
	if err != nil {
		log.Printf("Erreur lors de la mise à jour de l'avatar dans la base de données : %v", err)
		utils.RespondWithError(w, "Error updating avatar", http.StatusInternalServerError)
		return
	}

	utils.RespondWithSuccess(w, "Avatar updated successfully", http.StatusOK)
}

// DeleteAccount removes the user account from the database.
func (handler *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	// Récupération de l'ID utilisateur depuis le contexte de la requête
	userId, ok := r.Context().Value(utils.UserKey).(string)
	if !ok || userId == "" {
		utils.RespondWithError(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}

	// Appel de la méthode du repository pour supprimer le compte utilisateur
	err := handler.repos.UserRepo.DeleteUser(userId)
	if err != nil {
		utils.RespondWithError(w, "Erreur lors de la suppression du compte", http.StatusInternalServerError)
		return
	}

	// Réponse de succès
	utils.RespondWithSuccess(w, "Compte supprimé avec succès", http.StatusOK)
}
