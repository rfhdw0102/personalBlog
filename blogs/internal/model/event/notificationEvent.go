package event

import "time"

type NotificationType string

const (
	NotificationTypeLike    NotificationType = "like"
	NotificationTypeComment NotificationType = "comment"
	NotificationTypeReply   NotificationType = "reply"
)

type NotificationEvent struct {
	Type        NotificationType `json:"type"`
	RecipientID int              `json:"recipient_id"`
	ActorID     int              `json:"actor_id"`
	ArticleID   int              `json:"article_id"`
	CommentID   int              `json:"comment_id,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}
