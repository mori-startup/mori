package models

import "time"

type ChatMessage struct {
	ID         string    `json:"id"`
	SenderId   string    `json:"senderId"`
	ReceiverId string    `json:"receiverId"`
	Type       string    `json:"type"` // GROUP|PERSON
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	IsRead     bool      `json:"isRead"`
	Sender     User      `json:"sender"`
}

type ChatStats struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	UnreadMsgCount int    `json:"unreadMsgCount"`
}

type ConversationMsg struct {
	ID              string `json:"id"`              // user_id ou group_id
	Type            string `json:"type"`            // "PERSON" ou "GROUP"
	Name            string `json:"name"`            // nickname / group name
	Avatar          string `json:"avatar"`          // chemin vers l'image
	LastMessage     string `json:"lastMessage"`     // contenu du dernier message
	LastMessageTime string `json:"lastMessageTime"` // date/heure du dernier message
}

type MsgRepository interface {
	Save(ChatMessage) error
	// get all for specific chat
	// needs  RECEIVER and SENDER as input
	GetAll(ChatMessage) ([]ChatMessage, error)
	GetAllGroup(userId, groupId string) ([]ChatMessage, error)
	GetUnread(userId string) ([]ChatStats, error)
	GetUnreadGroup(userId string) ([]ChatStats, error)
	// mark as read
	MarkAsRead(ChatMessage) error
	MarkAsReadGroup(ChatMessage) error

	GetConversationsMsg(userID string) ([]ConversationMsg, error)

	SaveGroupMsg(ChatMessage) error

	// returns list of user id's that have chat history with provided user
	GetChatHistoryIds(userId string) (map[string]bool, error)
	// responds tru if both users have chat history
	HasHistory(senderId, receiverId string) (bool, error)
}
