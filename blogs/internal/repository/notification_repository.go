package repository

import (
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"

	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *entity.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) ListByRecipientID(recipientID int, page, pageSize int) ([]response.NotificationInfo, int64, error) {
	var total int64
	db := r.db.Model(&entity.Notification{}).Where("recipient_id = ?", recipientID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var notifications []entity.Notification
	err := db.Preload("Actor").
		Preload("Article").
		Preload("Comment").
		Preload("Comment.ParentComment").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&notifications).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	res := make([]response.NotificationInfo, len(notifications))
	for i, n := range notifications {
		res[i] = response.NotificationInfo{
			ID:           n.ID,
			RecipientID:  n.RecipientID,
			ActorID:      n.ActorID,
			ActorName:    n.Actor.Username,
			ActorAvatar:  n.Actor.Avatar,
			ArticleID:    n.ArticleID,
			ArticleTitle: n.Article.Title,
			CommentID:    n.CommentID,
			Type:         n.Type,
			IsRead:       n.IsRead,
			CreatedAt:    n.CreatedAt,
		}

		// 安全获取评论内容
		if n.Comment != nil {
			res[i].CommentContent = n.Comment.Content
			// 如果是回复，获取父评论内容
			if n.Type == "reply" && n.Comment.ParentComment != nil {
				res[i].ParentCommentContent = n.Comment.ParentComment.Content
			}
		}
	}

	return res, total, nil
}

func (r *notificationRepository) MarkAsRead(notificationID int, recipientID int) error {
	return r.db.Model(&entity.Notification{}).Where("id = ? AND recipient_id = ?", notificationID, recipientID).Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(recipientID int) error {
	return r.db.Model(&entity.Notification{}).Where("recipient_id = ? AND is_read = ?", recipientID, false).Update("is_read", true).Error
}

func (r *notificationRepository) GetUnreadCount(recipientID int) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Notification{}).Where("recipient_id = ? AND is_read = ?", recipientID, false).Count(&count).Error
	return count, err
}

func (r *notificationRepository) Delete(notificationID int, recipientID int) error {
	return r.db.Where("id = ? AND recipient_id = ?", notificationID, recipientID).Delete(&entity.Notification{}).Error
}
