package service

import (
	"blogs/internal/model/dto/response"
)

type NotificationService interface {
	ListNotifications(recipientID int, page, pageSize int) ([]response.NotificationInfo, int64, error)
	MarkAsRead(notificationID int, recipientID int) error
	MarkAllAsRead(recipientID int) error
	GetUnreadCount(recipientID int) (int64, error)
	DeleteNotification(notificationID int, recipientID int) error
}
