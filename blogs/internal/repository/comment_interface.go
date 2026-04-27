package repository

import (
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"
)

type CommentRepository interface {
	AddComment(comment *entity.Comment) error
	GetComments(articleID int, status int) ([]*response.CommentInfo, error)
	HideComment(commentID int, status int) error
	DeleteComment(commentID int) error
	GetCommentByID(commentID int) (*entity.Comment, error)
	ListUserCommentsGet(userID int, page, pageSize int, status int) ([]response.CommentInfo, int64, error)
	ListUserCommentsTake(userID int, page, pageSize int) ([]response.CommentInfo, int64, error)
}
