package service

import (
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"
)

type CommentService interface {
	AddComment(comment *entity.Comment) error
	GetComments(articleID int, status int) ([]*response.CommentInfo, error)
	HideComment(commentID int, userID int, status int) error
	DeleteComment(commentID int, userID int) error
	ListUserCommentsGet(userID int, page, pageSize int, status int) ([]response.CommentInfo, int64, error)
	ListUserCommentsTake(userID int, page, pageSize int) ([]response.CommentInfo, int64, error)
}
