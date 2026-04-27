package repository

import (
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"
)

type NotificationRepository interface {
	Create(notification *entity.Notification) error
	ListByRecipientID(recipientID int, page, pageSize int) ([]response.NotificationInfo, int64, error)
	MarkAsRead(notificationID int, recipientID int) error
	MarkAllAsRead(recipientID int) error
	GetUnreadCount(recipientID int) (int64, error)
	Delete(notificationID int, recipientID int) error
}
