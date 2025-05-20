package handlers

import (
	"log"
	"net/http"

	"mori/pkg/utils"
	ws "mori/pkg/wsServer"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:8080" // Restrict to frontend's origin
	},
}

func (handler *Handler) SocketHandler(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	// Access user ID
	userId, ok := r.Context().Value(utils.UserKey).(string)
	if !ok || userId == "" {
		log.Println("Invalid or missing user ID")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error during connection upgrade: %v", err)
		return
	}

	// Create new client and register
	client := ws.NewClient(conn, wsServer.Repos, userId)
	wsServer.RegisterNewClient(client)

	// Start reading and writing routines
	go client.Writer()
	go client.Reader(wsServer)
}
