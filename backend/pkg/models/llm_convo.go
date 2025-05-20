package models

type Conversation struct {
	UserID          string              `json:"user_id"`
	ConversationID  string              `json:"conversation_id"`
	Convo           []map[string]string `json:"convo"`
	NewConversation bool                `json:"new_conversation"`
}

type LLMConvoRepository interface {
	SaveConvo(Conversation) error
	GetAllConvo(Conversation) ([]Conversation, error)
	GetConvo(Conversation) (Conversation, error)
	GetLastConvoID() (string, error)
	GetLastConvo() (Conversation, error)
	DeleteConvo(Conversation) error
	// get all for specific chat
}
