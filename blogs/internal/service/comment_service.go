package service

import (
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"
	"blogs/internal/repository"
	appErrors "blogs/pkg/errors"
)

type commentService struct {
	commentRepo repository.CommentRepository
	articleRepo repository.ArticleRepository
}

func NewCommentService(commentRepo repository.CommentRepository, articleRepo repository.ArticleRepository) CommentService {
	return &commentService{commentRepo: commentRepo, articleRepo: articleRepo}
}

func (s *commentService) AddComment(comment *entity.Comment) error {
	if err := s.commentRepo.AddComment(comment); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "发表评论失败", err)
	}
	return nil
}

func (s *commentService) GetComments(articleID int, status int) ([]*response.CommentInfo, error) {
	return s.commentRepo.GetComments(articleID, status)
}

func (s *commentService) HideComment(commentID int, userID int, status int) error {
	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return appErrors.Wrap(appErrors.CodeNotFound, "评论不存在", err)
	}

	article, err := s.articleRepo.GetByID(comment.ArticleID)
	if err != nil {
		return appErrors.Wrap(appErrors.CodeNotFound, "文章不存在", err)
	}

	if article.UserID != userID {
		return appErrors.New(appErrors.CodeForbidden, "无权操作此评论")
	}

	if err := s.commentRepo.HideComment(commentID, status); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "更新评论状态失败", err)
	}
	return nil
}

func (s *commentService) DeleteComment(commentID int, userID int) error {
	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return appErrors.Wrap(appErrors.CodeNotFound, "评论不存在", err)
	}
	if comment.UserID != userID {
		return appErrors.New(appErrors.CodeForbidden, "无权操作此评论")
	}
	if err := s.commentRepo.DeleteComment(commentID); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "删除评论失败", err)
	}
	return nil
}

func (s *commentService) ListUserCommentsGet(userID int, page, pageSize int, status int) ([]response.CommentInfo, int64, error) {
	return s.commentRepo.ListUserCommentsGet(userID, page, pageSize, status)
}

func (s *commentService) ListUserCommentsTake(userID int, page, pageSize int) ([]response.CommentInfo, int64, error) {
	return s.commentRepo.ListUserCommentsTake(userID, page, pageSize)
}
