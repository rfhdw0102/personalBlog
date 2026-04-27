package response

import (
	"blogs/internal/model/entity"
	"time"
)

type NotificationInfo struct {
	ID                   int                     `json:"id"`
	RecipientID          int                     `json:"recipient_id"`
	ActorID              int                     `json:"actor_id"`
	ActorName            string                  `json:"actor_name"`
	ActorAvatar          string                  `json:"actor_avatar"`
	ArticleID            int                     `json:"article_id"`
	ArticleTitle         string                  `json:"article_title"`
	CommentID            int                     `json:"comment_id"`
	CommentContent       string                  `json:"comment_content"`
	ParentCommentContent string                  `json:"parent_comment_content"`
	Type                 entity.NotificationType `json:"type"`
	IsRead               bool                    `json:"is_read"`
	CreatedAt            time.Time               `json:"created_at"`
}
