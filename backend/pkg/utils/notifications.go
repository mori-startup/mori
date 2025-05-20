package utils

import "mori/pkg/models"

// replace notification message content based on type
func DefineNotificationMsg(notif *models.Notification) {
	switch notif.Type {

	case "FOLLOW":
		notif.Content = " sent you a following request "
	case "GROUP_INVITE":
		notif.Content = " invited you to join group "
	case "GROUP_REQUEST":
		notif.Content = " has requested to join your group "
	case "CHAT_REQUEST":
		notif.Content = " wants to chat with you"
	}
}
