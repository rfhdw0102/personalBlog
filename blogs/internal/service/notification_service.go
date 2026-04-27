package service

import (
	"blogs/internal/model/dto/response"
	"blogs/internal/repository"
	appErrors "blogs/pkg/errors"
)

type notificationService struct {
	notificationRepo repository.NotificationRepository
}

func NewNotificationService(notificationRepo repository.NotificationRepository) NotificationService {
	return &notificationService{notificationRepo: notificationRepo}
}

func (s *notificationService) ListNotifications(recipientID int, page, pageSize int) ([]response.NotificationInfo, int64, error) {
	return s.notificationRepo.ListByRecipientID(recipientID, page, pageSize)
}

func (s *notificationService) MarkAsRead(notificationID int, recipientID int) error {
	if err := s.notificationRepo.MarkAsRead(notificationID, recipientID); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "标记消息已读失败", err)
	}
	return nil
}

func (s *notificationService) MarkAllAsRead(recipientID int) error {
	if err := s.notificationRepo.MarkAllAsRead(recipientID); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "全部标记已读失败", err)
	}
	return nil
}

func (s *notificationService) GetUnreadCount(recipientID int) (int64, error) {
	return s.notificationRepo.GetUnreadCount(recipientID)
}

func (s *notificationService) DeleteNotification(notificationID int, recipientID int) error {
	if err := s.notificationRepo.Delete(notificationID, recipientID); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "删除通知失败", err)
	}
	return nil
}
