package db

import (
	"database/sql"
	"encoding/json"
	"mori/pkg/models"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type LLMConvoRepository struct {
	DB *sql.DB
}

// Save inserts a new message into the conversations table.
func (repo *LLMConvoRepository) SaveConvo(convo models.Conversation) error {

	conversations, err := ConvoToPQ(convo)
	if err != nil {
		return err
	}

	if convo.NewConversation {

		query := `
		INSERT INTO conversations (user_id, conversation_id, convo, new_conversation)
		VALUES ($1, $2, $3::jsonb[], $4);
	`
		_, err = repo.DB.Exec(query, convo.UserID, convo.ConversationID, pq.Array(conversations), convo.NewConversation)
	} else {

		query := `
		UPDATE conversations
		SET convo = $1::jsonb[]
		WHERE conversation_id = $2
	`
		_, err = repo.DB.Exec(query, pq.Array(conversations), convo.ConversationID)

		query = `
		UPDATE conversations
		SET new_conversation = false
		WHERE conversation_id = $1
	`
		_, err = repo.DB.Exec(query, convo.NewConversation)
	}

	return err
}

// Get all conversations for a specific chat
func (repo *LLMConvoRepository) GetAllConvo(convo models.Conversation) ([]models.Conversation, error) {

	// conversations, err := transformConvo(convo)
	var jsonItems []string

	query := `
		SELECT user_id, conversation_id, convo, new_conversation
		FROM conversations
		WHERE user_id = $1
	`
	rows, err := repo.DB.Query(query, convo.UserID)
	if err != nil {
		return nil, err
	}

	var convos []models.Conversation
	for rows.Next() {
		var convo models.Conversation
		if err := rows.Scan(&convo.UserID, &convo.ConversationID, pq.Array(&jsonItems), &convo.NewConversation); err != nil {
			return nil, err
		}
		convo.Convo, err = transformConvo(jsonItems)
		convos = append(convos, convo)
	}
	defer rows.Close()

	return convos, rows.Err()
}

func (repo *LLMConvoRepository) GetConvo(convo models.Conversation) (models.Conversation, error) {

	// conversations, err := transformConvo(convo)
	var jsonItems []string

	query := `
		SELECT user_id, conversation_id, convo, new_conversation
		FROM conversations
		WHERE conversation_id = $1
	`
	rows, err := repo.DB.Query(query, convo.ConversationID)
	if err != nil {
		return models.Conversation{}, err
	}

	var convo_result models.Conversation
	for rows.Next() {
		if err := rows.Scan(&convo_result.UserID, &convo_result.ConversationID, pq.Array(&jsonItems), &convo_result.NewConversation); err != nil {
			return models.Conversation{}, err
		}
		convo_result.Convo, err = transformConvo(jsonItems)
		if err != nil {
			return models.Conversation{}, err
		}
	}
	defer rows.Close()

	return convo_result, rows.Err()
}

// get the last conversation_id from the conversations table
func (repo *LLMConvoRepository) GetLastConvoID() (string, error) {
	query := `
		SELECT conversation_id
		FROM conversations
		ORDER BY conversation_id DESC
		LIMIT 1
	`
	var convoID string
	err := repo.DB.QueryRow(query).Scan(&convoID)
	return convoID, err
}

// get the last conversation from the conversations table
func (repo *LLMConvoRepository) GetLastConvo() (models.Conversation, error) {

	query := `
		SELECT conversation_id, user_id, convo, new_conversation
		FROM conversations
		ORDER BY conversation_id DESC
		LIMIT 1
	`
	var convo models.Conversation
	err := repo.DB.QueryRow(query).Scan(&convo.ConversationID, &convo.UserID, &convo.Convo, &convo.NewConversation)
	return convo, err
}

// delete a conversation from the conversations table
func (repo *LLMConvoRepository) DeleteConvo(convo models.Conversation) error {
	query := `
		DELETE FROM conversations
		WHERE conversation_id = $1 AND user_id = $2
	`
	_, err := repo.DB.Exec(query, convo.ConversationID, convo.UserID)

	return err
}

func ConvoToPQ(convo models.Conversation) ([]string, error) {
	var jsonItems []string
	for _, m := range convo.Convo {
		j, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		jsonItems = append(jsonItems, string(j))
	}
	return jsonItems, nil
}

func transformConvo(convo []string) ([]map[string]string, error) {
	var jsonItems []map[string]string
	for _, m := range convo {
		var j map[string]string
		err := json.Unmarshal([]byte(m), &j)
		if err != nil {
			return nil, err
		}
		jsonItems = append(jsonItems, j)
	}
	return jsonItems, nil
}
