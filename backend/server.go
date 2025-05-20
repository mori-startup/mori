package main

import (
	"fmt"
	"net/http"
	"time"

	sqlite "mori/pkg/db/PostgreSql"
	"mori/pkg/handlers"
	"mori/pkg/middleware"
	"mori/pkg/utils"
	ws "mori/pkg/wsServer"
)

func main() {
	// initialize database
	db := sqlite.InitDB()
	defer db.Close()
	// temp
	// initialize repositories
	repos := sqlite.InitRepositories(db)
	// initialize handlers with connection to repositories
	handler := handlers.InitHandlers(repos)
	// initialize wsServer
	wsServer := ws.StartServer(repos)

	// set up server address and routes
	server := &http.Server{
		Addr:         ":8081",
		Handler:      middleware.RateLimit(setRoutes(handler, wsServer)),
		ReadTimeout:  5 * time.Second, // limite de lecture requête
		WriteTimeout: 2 * time.Minute,
	}

	fmt.Printf("Server started at http://localhost" + server.Addr + "\n")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error", err)
	}
}

// Set up all routes
func setRoutes(handler *handlers.Handler, wsServer *ws.Server) http.Handler {
	mux := http.NewServeMux()
	/* ------------------------------ image server ------------------------------ */
	fs := http.FileServer(http.Dir("./imageUpload"))
	mux.Handle("/imageUpload/", http.StripPrefix("/imageUpload/", utils.ConfigFSHeader(fs)))
	/* ------------------------------- auth route ------------------------------- */
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/signin", handler.Signin)
	mux.HandleFunc("/logout", handler.Auth(handler.Logout))
	mux.HandleFunc("/captcha", handler.ServeCaptcha)
	mux.HandleFunc("/verified", handler.VerifyEmail)
	mux.HandleFunc("/sessionActive", handler.SessionActive)
	mux.HandleFunc("/request-password-reset", handler.RequestPasswordReset)
	mux.HandleFunc("/reset-password", handler.ResetPassword)

	/* ------------------------------- LLM_conv ------------------------------- */
	mux.HandleFunc("/llmConvo", handler.Auth(handler.LLMHandler))
	mux.HandleFunc("/llmConvoSave", handler.Auth(handler.LLMConvoSave))
	mux.HandleFunc("/llmConvoGet", handler.Auth(handler.LLMConvoGet))
	mux.HandleFunc("/llmConvoSelected", handler.Auth(handler.LLMConvoSelected))
	mux.HandleFunc("/llmConvoGetLast", handler.Auth(handler.LLMConvoGetLast))
	mux.HandleFunc("/llmConvoDelete", handler.Auth(handler.LLMConvoDelete))

	/* ---------------------------------- users --------------------------------- */
	mux.HandleFunc("/allUsers", handler.Auth(handler.AllUsers))       // all users + info except current
	mux.HandleFunc("/followers", handler.Auth(handler.GetFollowers))  // follower list
	mux.HandleFunc("/following", handler.Auth(handler.GetFollowing))  // following list
	mux.HandleFunc("/currentUser", handler.Auth(handler.CurrentUser)) // current user data
	mux.HandleFunc("/userData", handler.Auth(handler.UserData))       // userd data based on following status
	mux.HandleFunc("/changeStatus", handler.Auth(handler.UserStatus)) // change status

	mux.HandleFunc("/follow", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.Follow(wsServer, w, r)
	})) // follow user
	mux.HandleFunc("/cancelFollowRequest", handler.Auth(handler.CancelFollowRequest))
	mux.HandleFunc("/unfollow", handler.Auth(handler.Unfollow))
	mux.HandleFunc("/responseFollowRequest", handler.Auth(handler.ResponseFollowRequest))

	/* --------------------------------- groups --------------------------------- */
	mux.HandleFunc("/allGroups", handler.Auth(handler.AllGroups))             // group list
	mux.HandleFunc("/userGroups", handler.Auth(handler.UserGroups))           // group list of user groups
	mux.HandleFunc("/otherUserGroups", handler.Auth(handler.OtherUserGroups)) // group list for specific user

	mux.HandleFunc("/groupInfo", handler.Auth(handler.GroupInfo))       // get group info
	mux.HandleFunc("/groupMembers", handler.Auth(handler.GroupMembers)) // get group members

	mux.HandleFunc("/groupRequests", handler.Auth(handler.GroupRequests))             // get group member requests
	mux.HandleFunc("/cancelGroupRequests", handler.Auth(handler.CancelGroupRequests)) // cancel request or joing group

	mux.HandleFunc("/newGroup", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.NewGroup(wsServer, w, r)
	})) // create new group

	mux.HandleFunc("/newGroupInvite", handler.Auth(func(w http.ResponseWriter, r *http.Request) { // invite new users to group
		handler.NewGroupInvite(wsServer, w, r)
	}))
	mux.HandleFunc("/newGroupRequest", handler.Auth(func(w http.ResponseWriter, r *http.Request) { // invite new users to group
		handler.NewGroupRequest(wsServer, w, r)
	}))
	mux.HandleFunc("/responseGroupRequest", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.ResponseGroupRequest(wsServer, w, r)
	})) // response to join request
	mux.HandleFunc("/responseInviteRequest", handler.Auth(handler.ResponseInviteRequest)) // response to invite request

	/* ------------------------------ notifications ----------------------------- */
	mux.HandleFunc("/notifications", handler.Auth(handler.Notifications)) // get all notifs from db on login

	/* ------------------------------ chat messages ----------------------------- */
	mux.HandleFunc("/messages", handler.Auth(handler.Messages))             // get all chat messages for specific chat
	mux.HandleFunc("/unreadMessages", handler.Auth(handler.UnreadMessages)) // get list of messages that isn't read

	mux.HandleFunc("/messageRead", handler.Auth(handler.MessageRead)) // mark message as read
	mux.HandleFunc("/newMessage", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.NewMessage(wsServer, w, r)
	})) // new chat message
	mux.HandleFunc("/chatList", handler.Auth(handler.ChatList))                       // get list of users to display in chatbox
	mux.HandleFunc("/responseChatRequest", handler.Auth(handler.ResponseChatRequest)) // response to chat request
	/* ---------------------------- ConversationMsg Sidebar ---------------------------- */
	mux.HandleFunc("/conversationsMsg", handler.Auth(handler.ConversationsMsg)) // get list of users to display in chatbox
	/* ---------------------------- Update Profile ---------------------------- */
	mux.HandleFunc("/updateNickname", handler.Auth(handler.ChangeNickname)) // Update user nickname
	mux.HandleFunc("/updateAvatar", handler.Auth(handler.ChangeAvatar))     // Update user avatar
	mux.HandleFunc("/DeleteAccount", handler.Auth(handler.DeleteAccount))

	/* ---------------------------- websocket server ---------------------------- */
	mux.HandleFunc("/ws", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.SocketHandler(wsServer, w, r)
	}))

	// File upload endpoints.
	mux.HandleFunc("/api/upload", handler.Auth(handler.UploadFiles))
	mux.HandleFunc("/api/files", handler.Auth(handler.ListFiles))
	mux.HandleFunc("/api/files/", handler.Auth(handler.DeleteFile))

	return mux
}
