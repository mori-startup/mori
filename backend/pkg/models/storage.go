package models

import "database/sql"

// Repositories contains all the repo structs
type Repositories struct {
	DB           *sql.DB
	UserRepo     UserRepository
	SessionRepo  SessionRepository
	GroupRepo    GroupRepository
	NotifRepo    NotifRepository
	MsgRepo      MsgRepository
	LLMConvoRepo LLMConvoRepository
}
