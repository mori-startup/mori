package db

import (
	"database/sql"

	"mori/pkg/models"
)

type SessionRepository struct {
	DB *sql.DB
}

// Set inserts a new session into the database.
func (repo *SessionRepository) Set(session models.Session) error {
	query := `
		INSERT INTO sessions (session_id, user_id, expiration_time) 
		VALUES ($1, $2, $3);
	`
	_, err := repo.DB.Exec(query, session.ID, session.UserID, session.ExpirationTime)
	return err
}

// Get retrieves a session based on the session ID.
func (repo *SessionRepository) Get(sessionID string) (models.Session, error) {
	query := `
		SELECT user_id, expiration_time 
		FROM sessions 
		WHERE session_id = $1 
		LIMIT 1;
	`
	var session models.Session
	err := repo.DB.QueryRow(query, sessionID).Scan(&session.UserID, &session.ExpirationTime)
	if err != nil {
		return session, err
	}
	session.ID = sessionID
	return session, nil
}

// Delete removes a session from the database based on the user ID.
func (repo *SessionRepository) Delete(session models.Session) error {
	query := `
		DELETE FROM sessions 
		WHERE user_id = $1;
	`
	_, err := repo.DB.Exec(query, session.UserID)
	return err
}

// Update modifies an existing session based on the user ID.
func (repo *SessionRepository) Update(session models.Session) error {
	query := `
		UPDATE sessions 
		SET expiration_time = $1, session_id = $2 
		WHERE user_id = $3;
	`
	_, err := repo.DB.Exec(query, session.ExpirationTime, session.ID, session.UserID)
	return err
}

// GetByUser checks if a session exists based on the user ID and retrieves it.
func (repo *SessionRepository) GetByUser(userID string) (models.Session, error) {
	query := `
		SELECT session_id, expiration_time 
		FROM sessions 
		WHERE user_id = $1 
		LIMIT 1;
	`
	var session models.Session
	err := repo.DB.QueryRow(query, userID).Scan(&session.ID, &session.ExpirationTime)
	if err != nil {
		return session, err
	}
	session.UserID = userID
	return session, nil
}
