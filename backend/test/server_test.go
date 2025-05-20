package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ws "mori/pkg/wsServer"
)

// FullHandlerInterface defines all the methods used by setRoutes.
type FullHandlerInterface interface {
	Register(http.ResponseWriter, *http.Request)
	Signin(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
	ServeCaptcha(http.ResponseWriter, *http.Request)
	VerifyEmail(http.ResponseWriter, *http.Request)
	SessionActive(http.ResponseWriter, *http.Request)
	RequestPasswordReset(http.ResponseWriter, *http.Request)
	ResetPassword(http.ResponseWriter, *http.Request)
	LLMHandler(http.ResponseWriter, *http.Request)
	AllUsers(http.ResponseWriter, *http.Request)
	GetFollowers(http.ResponseWriter, *http.Request)
	GetFollowing(http.ResponseWriter, *http.Request)
	CurrentUser(http.ResponseWriter, *http.Request)
	UserData(http.ResponseWriter, *http.Request)
	UserStatus(http.ResponseWriter, *http.Request)
	Follow(*ws.Server, http.ResponseWriter, *http.Request)
	CancelFollowRequest(http.ResponseWriter, *http.Request)
	Unfollow(http.ResponseWriter, *http.Request)
	ResponseFollowRequest(http.ResponseWriter, *http.Request)
	AllGroups(http.ResponseWriter, *http.Request)
	UserGroups(http.ResponseWriter, *http.Request)
	OtherUserGroups(http.ResponseWriter, *http.Request)
	GroupInfo(http.ResponseWriter, *http.Request)
	GroupMembers(http.ResponseWriter, *http.Request)
	GroupRequests(http.ResponseWriter, *http.Request)
	CancelGroupRequests(http.ResponseWriter, *http.Request)
	NewGroup(*ws.Server, http.ResponseWriter, *http.Request)
	NewGroupInvite(*ws.Server, http.ResponseWriter, *http.Request)
	NewGroupRequest(*ws.Server, http.ResponseWriter, *http.Request)
	ResponseGroupRequest(*ws.Server, http.ResponseWriter, *http.Request)
	ResponseInviteRequest(http.ResponseWriter, *http.Request)
	Notifications(http.ResponseWriter, *http.Request)
	Messages(http.ResponseWriter, *http.Request)
	UnreadMessages(http.ResponseWriter, *http.Request)
	MessageRead(http.ResponseWriter, *http.Request)
	NewMessage(*ws.Server, http.ResponseWriter, *http.Request)
	ChatList(http.ResponseWriter, *http.Request)
	ResponseChatRequest(http.ResponseWriter, *http.Request)
	ConversationsMsg(http.ResponseWriter, *http.Request)
	ChangeNickname(http.ResponseWriter, *http.Request)
	ChangeAvatar(http.ResponseWriter, *http.Request)
	SocketHandler(*ws.Server, http.ResponseWriter, *http.Request)
	UploadFiles(http.ResponseWriter, *http.Request)
	ListFiles(http.ResponseWriter, *http.Request)
	DeleteFile(http.ResponseWriter, *http.Request)
	Auth(http.HandlerFunc) http.HandlerFunc
}

// dummyHandler implements FullHandlerInterface with dummy responses.
// Each method simply writes its endpoint name as the response.
type dummyHandler struct{}

func (d *dummyHandler) Register(w http.ResponseWriter, r *http.Request) { w.Write([]byte("register")) }
func (d *dummyHandler) Signin(w http.ResponseWriter, r *http.Request)   { w.Write([]byte("signin")) }
func (d *dummyHandler) Logout(w http.ResponseWriter, r *http.Request)   { w.Write([]byte("logout")) }
func (d *dummyHandler) ServeCaptcha(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("captcha"))
}

func (d *dummyHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("verified"))
}

func (d *dummyHandler) SessionActive(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sessionActive"))
}

func (d *dummyHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("requestPasswordReset"))
}

func (d *dummyHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("resetPassword"))
}

func (d *dummyHandler) LLMHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("llmConvo"))
}
func (d *dummyHandler) AllUsers(w http.ResponseWriter, r *http.Request) { w.Write([]byte("allUsers")) }
func (d *dummyHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("followers"))
}

func (d *dummyHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("following"))
}

func (d *dummyHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("currentUser"))
}
func (d *dummyHandler) UserData(w http.ResponseWriter, r *http.Request) { w.Write([]byte("userData")) }
func (d *dummyHandler) UserStatus(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("userStatus"))
}

func (d *dummyHandler) Follow(ws *ws.Server, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("follow"))
}

func (d *dummyHandler) CancelFollowRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("cancelFollowRequest"))
}
func (d *dummyHandler) Unfollow(w http.ResponseWriter, r *http.Request) { w.Write([]byte("unfollow")) }
func (d *dummyHandler) ResponseFollowRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("responseFollowRequest"))
}

func (d *dummyHandler) AllGroups(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("allGroups"))
}

func (d *dummyHandler) UserGroups(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("userGroups"))
}

func (d *dummyHandler) OtherUserGroups(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("otherUserGroups"))
}

func (d *dummyHandler) GroupInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("groupInfo"))
}

func (d *dummyHandler) GroupMembers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("groupMembers"))
}

func (d *dummyHandler) GroupRequests(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("groupRequests"))
}

func (d *dummyHandler) CancelGroupRequests(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("cancelGroupRequests"))
}

func (d *dummyHandler) NewGroup(ws *ws.Server, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("newGroup"))
}

func (d *dummyHandler) NewGroupInvite(ws *ws.Server, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("newGroupInvite"))
}

func (d *dummyHandler) NewGroupRequest(ws *ws.Server, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("newGroupRequest"))
}

func (d *dummyHandler) ResponseGroupRequest(ws *ws.Server, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("responseGroupRequest"))
}

func (d *dummyHandler) ResponseInviteRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("responseInviteRequest"))
}

func (d *dummyHandler) Notifications(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("notifications"))
}
func (d *dummyHandler) Messages(w http.ResponseWriter, r *http.Request) { w.Write([]byte("messages")) }
func (d *dummyHandler) UnreadMessages(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("unreadMessages"))
}

func (d *dummyHandler) MessageRead(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("messageRead"))
}

func (d *dummyHandler) NewMessage(ws *ws.Server, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("newMessage"))
}
func (d *dummyHandler) ChatList(w http.ResponseWriter, r *http.Request) { w.Write([]byte("chatList")) }
func (d *dummyHandler) ResponseChatRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("responseChatRequest"))
}

func (d *dummyHandler) ConversationsMsg(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("conversationsMsg"))
}

func (d *dummyHandler) ChangeNickname(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("updateNickname"))
}

func (d *dummyHandler) ChangeAvatar(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("updateAvatar"))
}

func (d *dummyHandler) SocketHandler(ws *ws.Server, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ws"))
}

func (d *dummyHandler) UploadFiles(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("uploadFiles"))
}

func (d *dummyHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("listFiles"))
}

func (d *dummyHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("deleteFile"))
}

func (d *dummyHandler) Auth(fn http.HandlerFunc) http.HandlerFunc {
	return fn
}

// dummyWSServer is a dummy implementation for ws.Server.
type dummyWSServer struct{}

// newDummyWSServer returns a dummy instance of ws.Server.
func newDummyWSServer() *ws.Server {
	// Adjust this if ws.Server requires specific initialization.
	return &ws.Server{}
}

// setRoutesForTest sets up the routes as in your server.go file.
func setRoutesForTest(handler FullHandlerInterface, wsServer *ws.Server) http.Handler {
	mux := http.NewServeMux()

	// Omit the file server ("/imageUpload/") endpoint for testing.
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/signin", handler.Signin)
	mux.HandleFunc("/logout", handler.Auth(handler.Logout))
	mux.HandleFunc("/captcha", handler.ServeCaptcha)
	mux.HandleFunc("/verified", handler.VerifyEmail)
	mux.HandleFunc("/sessionActive", handler.SessionActive)
	mux.HandleFunc("/request-password-reset", handler.RequestPasswordReset)
	mux.HandleFunc("/reset-password", handler.ResetPassword)

	mux.HandleFunc("/llmConvo", handler.Auth(handler.LLMHandler))

	mux.HandleFunc("/allUsers", handler.Auth(handler.AllUsers))
	mux.HandleFunc("/followers", handler.Auth(handler.GetFollowers))
	mux.HandleFunc("/following", handler.Auth(handler.GetFollowing))
	mux.HandleFunc("/currentUser", handler.Auth(handler.CurrentUser))
	mux.HandleFunc("/userData", handler.Auth(handler.UserData))
	mux.HandleFunc("/changeStatus", handler.Auth(handler.UserStatus))

	mux.HandleFunc("/follow", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.Follow(wsServer, w, r)
	}))
	mux.HandleFunc("/cancelFollowRequest", handler.Auth(handler.CancelFollowRequest))
	mux.HandleFunc("/unfollow", handler.Auth(handler.Unfollow))
	mux.HandleFunc("/responseFollowRequest", handler.Auth(handler.ResponseFollowRequest))

	mux.HandleFunc("/allGroups", handler.Auth(handler.AllGroups))
	mux.HandleFunc("/userGroups", handler.Auth(handler.UserGroups))
	mux.HandleFunc("/otherUserGroups", handler.Auth(handler.OtherUserGroups))

	mux.HandleFunc("/groupInfo", handler.Auth(handler.GroupInfo))
	mux.HandleFunc("/groupMembers", handler.Auth(handler.GroupMembers))

	mux.HandleFunc("/groupRequests", handler.Auth(handler.GroupRequests))
	mux.HandleFunc("/cancelGroupRequests", handler.Auth(handler.CancelGroupRequests))

	mux.HandleFunc("/newGroup", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.NewGroup(wsServer, w, r)
	}))
	mux.HandleFunc("/newGroupInvite", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.NewGroupInvite(wsServer, w, r)
	}))
	mux.HandleFunc("/newGroupRequest", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.NewGroupRequest(wsServer, w, r)
	}))
	mux.HandleFunc("/responseGroupRequest", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.ResponseGroupRequest(wsServer, w, r)
	}))
	mux.HandleFunc("/responseInviteRequest", handler.Auth(handler.ResponseInviteRequest))

	mux.HandleFunc("/notifications", handler.Auth(handler.Notifications))

	mux.HandleFunc("/messages", handler.Auth(handler.Messages))
	mux.HandleFunc("/unreadMessages", handler.Auth(handler.UnreadMessages))
	mux.HandleFunc("/messageRead", handler.Auth(handler.MessageRead))
	mux.HandleFunc("/newMessage", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.NewMessage(wsServer, w, r)
	}))
	mux.HandleFunc("/chatList", handler.Auth(handler.ChatList))
	mux.HandleFunc("/responseChatRequest", handler.Auth(handler.ResponseChatRequest))
	mux.HandleFunc("/conversationsMsg", handler.Auth(handler.ConversationsMsg))

	mux.HandleFunc("/updateNickname", handler.Auth(handler.ChangeNickname))
	mux.HandleFunc("/updateAvatar", handler.Auth(handler.ChangeAvatar))

	mux.HandleFunc("/ws", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.SocketHandler(wsServer, w, r)
	}))

	mux.HandleFunc("/api/upload", handler.Auth(handler.UploadFiles))
	mux.HandleFunc("/api/files", handler.Auth(handler.ListFiles))
	mux.HandleFunc("/api/files/", handler.Auth(handler.DeleteFile))

	return mux
}

// TestAllEndpointsVisual tests all endpoints and prints a visual ASCII table.
func TestAllEndpointsVisual(t *testing.T) {
	dh := &dummyHandler{}
	wsServer := newDummyWSServer()
	router := setRoutesForTest(dh, wsServer)

	// Define test cases for each endpoint.
	testCases := []struct {
		name     string
		endpoint string
		expected string
	}{
		{"Register", "/register", "register"},
		{"Signin", "/signin", "signin"},
		{"Logout", "/logout", "logout"},
		{"Captcha", "/captcha", "captcha"},
		{"Verified", "/verified", "verified"},
		{"SessionActive", "/sessionActive", "sessionActive"},
		{"RequestPasswordReset", "/request-password-reset", "requestPasswordReset"},
		{"ResetPassword", "/reset-password", "resetPassword"},
		{"LLMHandler", "/llmConvo", "llmConvo"},
		{"AllUsers", "/allUsers", "allUsers"},
		{"Followers", "/followers", "followers"},
		{"Following", "/following", "following"},
		{"CurrentUser", "/currentUser", "currentUser"},
		{"UserData", "/userData", "userData"},
		{"ChangeStatus", "/changeStatus", "userStatus"},
		{"Follow", "/follow", "follow"},
		{"CancelFollowRequest", "/cancelFollowRequest", "cancelFollowRequest"},
		{"Unfollow", "/unfollow", "unfollow"},
		{"ResponseFollowRequest", "/responseFollowRequest", "responseFollowRequest"},
		{"AllGroups", "/allGroups", "allGroups"},
		{"UserGroups", "/userGroups", "userGroups"},
		{"OtherUserGroups", "/otherUserGroups", "otherUserGroups"},
		{"GroupInfo", "/groupInfo", "groupInfo"},
		{"GroupMembers", "/groupMembers", "groupMembers"},
		{"GroupRequests", "/groupRequests", "groupRequests"},
		{"CancelGroupRequests", "/cancelGroupRequests", "cancelGroupRequests"},
		{"NewGroup", "/newGroup", "newGroup"},
		{"NewGroupInvite", "/newGroupInvite", "newGroupInvite"},
		{"NewGroupRequest", "/newGroupRequest", "newGroupRequest"},
		{"ResponseGroupRequest", "/responseGroupRequest", "responseGroupRequest"},
		{"ResponseInviteRequest", "/responseInviteRequest", "responseInviteRequest"},
		{"Notifications", "/notifications", "notifications"},
		{"Messages", "/messages", "messages"},
		{"UnreadMessages", "/unreadMessages", "unreadMessages"},
		{"MessageRead", "/messageRead", "messageRead"},
		{"NewMessage", "/newMessage", "newMessage"},
		{"ChatList", "/chatList", "chatList"},
		{"ResponseChatRequest", "/responseChatRequest", "responseChatRequest"},
		{"ConversationsMsg", "/conversationsMsg", "conversationsMsg"},
		{"UpdateNickname", "/updateNickname", "updateNickname"},
		{"UpdateAvatar", "/updateAvatar", "updateAvatar"},
		{"WS", "/ws", "ws"},
		{"UploadFiles", "/api/upload", "uploadFiles"},
		{"ListFiles", "/api/files", "listFiles"},
		{"DeleteFile", "/api/files/", "deleteFile"},
	}

	// Log header for the ASCII table.
	t.Log("\n+----------------------+----------------+----------------+")
	t.Log("| Endpoint             | Expected Output| Actual Output  |")
	t.Log("+----------------------+----------------+----------------+")

	// Run each test case.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.endpoint, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			actual := rec.Body.String()
			t.Logf("| %-20s | %-14s | %-14s |", tc.endpoint, tc.expected, actual)
			if actual != tc.expected {
				t.Errorf("For endpoint %s: expected %q, got %q", tc.endpoint, tc.expected, actual)
			}
		})
	}

	// Log footer for the table.
	t.Log("+----------------------+----------------+----------------+")
}
