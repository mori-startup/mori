package db

import (
	"database/sql"

	"mori/pkg/models"
)

type NotifRepository struct {
	DB *sql.DB
}

// Save inserts a new notification into the notifications table.
func (repo *NotifRepository) Save(notification models.Notification) error {
	query := `
		INSERT INTO notifications (notif_id, user_id, type, content, sender) 
		VALUES ($1, $2, $3, $4, $5);
	`
	_, err := repo.DB.Exec(query, notification.ID, notification.TargetID, notification.Type, notification.Content, notification.Sender)
	return err
}

// Delete removes a notification by ID.
func (repo *NotifRepository) Delete(notificationId string) error {
	query := `
		DELETE FROM notifications 
		WHERE notif_id = $1;
	`
	_, err := repo.DB.Exec(query, notificationId)
	return err
}

// DeleteByType removes notifications of a specific type and content for a user.
func (repo *NotifRepository) DeleteByType(notif models.Notification) error {
	query := `
		DELETE FROM notifications 
		WHERE user_id = $1 AND type = $2 AND content = $3;
	`
	_, err := repo.DB.Exec(query, notif.TargetID, notif.Type, notif.Content)
	return err
}

// GetGroupRequests retrieves group request notifications for a specific group ID.
func (repo *NotifRepository) GetGroupRequests(groupId string) ([]models.Notification, error) {
	query := `
		SELECT content, notif_id, type, sender, user_id 
		FROM notifications 
		WHERE user_id = $1 AND type = 'GROUP_REQUEST';
	`
	rows, err := repo.DB.Query(query, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		if err := rows.Scan(&notif.Content, &notif.ID, &notif.Type, &notif.Sender, &notif.TargetID); err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}
	return notifications, rows.Err()
}

// GetUserFromRequest retrieves the user ID from a request notification.
func (repo *NotifRepository) GetUserFromRequest(notificationId string) (string, error) {
	query := `
		SELECT content 
		FROM notifications 
		WHERE notif_id = $1;
	`
	var userId string
	err := repo.DB.QueryRow(query, notificationId).Scan(&userId)
	return userId, err
}

// CheckIfExists checks if a notification exists for a user with the same content and type.
func (repo *NotifRepository) CheckIfExists(notif models.Notification) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM notifications 
		WHERE user_id = $1 AND content = $2 AND type = $3;
	`
	var count int
	err := repo.DB.QueryRow(query, notif.TargetID, notif.Content, notif.Type).Scan(&count)
	return count > 0, err
}

// GetGroupId retrieves the group ID from a notification.
func (repo *NotifRepository) GetGroupId(notificationId string) (string, error) {
	query := `
		SELECT content 
		FROM notifications 
		WHERE notif_id = $1;
	`
	var groupId string
	err := repo.DB.QueryRow(query, notificationId).Scan(&groupId)
	return groupId, err
}

// GetAll retrieves all notifications for a user.
func (repo *NotifRepository) GetAll(userId string) ([]models.Notification, error) {
	query := `
		SELECT content, notif_id, type, sender, user_id 
		FROM notifications 
		WHERE user_id = $1 
		   OR (SELECT administrator FROM groups WHERE group_id = notifications.user_id) = $1;
	`
	rows, err := repo.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		if err := rows.Scan(&notif.Content, &notif.ID, &notif.Type, &notif.Sender, &notif.TargetID); err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}
	return notifications, rows.Err()
}

// GetChatNotifById retrieves a chat notification by its ID.
func (repo *NotifRepository) GetChatNotifById(notificationId string) (models.Notification, error) {
	query := `
		SELECT content, user_id, sender 
		FROM notifications 
		WHERE notif_id = $1;
	`
	var notif models.Notification
	err := repo.DB.QueryRow(query, notificationId).Scan(&notif.Content, &notif.TargetID, &notif.Sender)
	return notif, err
}

// CheckIfChatRequestExists checks if a chat request notification exists between two users.
func (repo *NotifRepository) CheckIfChatRequestExists(senderId, receiverId string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM notifications 
		WHERE user_id = $1 AND sender = $2 AND type = 'CHAT_REQUEST';
	`
	var count int
	err := repo.DB.QueryRow(query, receiverId, senderId).Scan(&count)
	return count > 0, err
}

// GetContentFromChatRequest retrieves the content of a chat request between two users.
func (repo *NotifRepository) GetContentFromChatRequest(senderId, receiverId string) (string, error) {
	query := `
		SELECT content 
		FROM notifications 
		WHERE user_id = $1 AND sender = $2 AND type = 'CHAT_REQUEST';
	`
	var content string
	err := repo.DB.QueryRow(query, receiverId, senderId).Scan(&content)
	return content, err
}
