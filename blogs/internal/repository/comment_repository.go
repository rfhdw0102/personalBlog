package repository

import (
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"
	"blogs/internal/model/event"
	"blogs/pkg/logger"
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type commentRepository struct {
	db        *gorm.DB
	redisRepo RedisRepository
}

func NewCommentRepository(db *gorm.DB, redisRepo RedisRepository) CommentRepository {
	return &commentRepository{db: db, redisRepo: redisRepo}
}

func (r *commentRepository) AddComment(comment *entity.Comment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			return err
		}

		ev, err := r.buildNotificationEvent(tx, comment)
		if err != nil {
			logger.Warn("构建通知失败，跳过", zap.Error(err))
			return nil
		}

		return r.publishNotification(ev)
	})
}

func (r *commentRepository) buildNotificationEvent(tx *gorm.DB, comment *entity.Comment) (event.NotificationEvent, error) {
	if comment.ParentID > 0 {
		return r.buildReplyNotification(tx, comment)
	}
	return r.buildCommentNotification(tx, comment)
}

func (r *commentRepository) buildReplyNotification(tx *gorm.DB, comment *entity.Comment) (event.NotificationEvent, error) {
	var parent entity.Comment
	if err := tx.Select("user_id").First(&parent, comment.ParentID).Error; err != nil {
		return event.NotificationEvent{}, err
	}

	return event.NotificationEvent{
		Type:        event.NotificationTypeReply,
		RecipientID: parent.UserID,
		ActorID:     comment.UserID,
		ArticleID:   comment.ArticleID,
		CommentID:   comment.ID,
		CreatedAt:   time.Now(),
	}, nil
}

func (r *commentRepository) buildCommentNotification(tx *gorm.DB, comment *entity.Comment) (event.NotificationEvent, error) {
	var article entity.Article
	if err := tx.Select("user_id").First(&article, comment.ArticleID).Error; err != nil {
		return event.NotificationEvent{}, err
	}

	return event.NotificationEvent{
		Type:        event.NotificationTypeComment,
		RecipientID: article.UserID,
		ActorID:     comment.UserID,
		ArticleID:   comment.ArticleID,
		CommentID:   comment.ID,
		CreatedAt:   time.Now(),
	}, nil
}

func (r *commentRepository) publishNotification(ev event.NotificationEvent) error {
	evJSON, err := json.Marshal(ev)
	if err != nil {
		logger.Error("json序列化失败", zap.Error(err))
		return err
	}

	_, err = r.redisRepo.XAdd(context.Background(), "notifications_stream", map[string]interface{}{
		"event": string(evJSON),
	})
	if err != nil {
		logger.Warn("通知发送失败", zap.Error(err))
	}
	return nil
}

func (r *commentRepository) GetComments(articleID int, status int) ([]*response.CommentInfo, error) {
	flatList, err := r.findComments(articleID, status)
	if err != nil {
		return nil, err
	}
	return r.buildCommentTree(flatList), nil
}

func (r *commentRepository) findComments(articleID int, status int) ([]response.CommentInfo, error) {
	var list []response.CommentInfo

	err := r.db.Model(&entity.Comment{}).
		Joins("LEFT JOIN users ON comments.user_id = users.id").
		Select("comments.id, comments.article_id, comments.user_id, comments.content, comments.parent_id, comments.status, comments.created_at, users.username, users.avatar").
		Where("comments.article_id = ? AND comments.status = ?", articleID, status).
		Order("comments.created_at asc").
		Find(&list).Error

	return list, err
}

func (r *commentRepository) buildCommentTree(flatList []response.CommentInfo) []*response.CommentInfo {
	nodeMap := make(map[int]*response.CommentInfo, len(flatList))

	for i := range flatList {
		nodeMap[flatList[i].ID] = &flatList[i]
	}

	var roots []*response.CommentInfo

	for i := range flatList {
		node := &flatList[i]

		if node.ParentID == 0 {
			roots = append(roots, node)
			continue
		}

		if parent, ok := nodeMap[node.ParentID]; ok {
			parent.Children = append(parent.Children, node)
		} else {
			// 父评论缺失 → 降级
			roots = append(roots, node)
		}
	}

	return roots
}

func (r *commentRepository) HideComment(commentID int, status int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		ids, err := r.getAllDescendantIDs(tx, commentID)
		if err != nil {
			return err
		}
		ids = append(ids, commentID)

		return tx.Model(&entity.Comment{}).
			Where("id IN ?", ids).
			Update("status", status).Error
	})
}

func (r *commentRepository) DeleteComment(commentID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		ids, err := r.getAllDescendantIDs(tx, commentID)
		if err != nil {
			return err
		}
		ids = append(ids, commentID)

		return tx.Where("id IN ?", ids).Delete(&entity.Comment{}).Error
	})
}

func (r *commentRepository) getAllDescendantIDs(tx *gorm.DB, parentID int) ([]int, error) {
	var ids []int

	// WITH：定义临时结果集（CTE，Common Table Expression），查询中可复用
	//RECURSIVE：允许 CTE 引用自身，实现递归查询
	//comment_tree：给这个临时结果集起的名字
	//递归 CTE 专门用来处理树形结构数据（评论嵌套回复、部门层级、分类树等）
	sql := `
	WITH RECURSIVE comment_tree AS (
		SELECT id FROM comments WHERE id = ?
		UNION ALL
		SELECT c.id FROM comments c
		INNER JOIN comment_tree ct ON c.parent_id = ct.id
	)
	SELECT id FROM comment_tree WHERE id != ?
	`

	err := tx.Raw(sql, parentID, parentID).Pluck("id", &ids).Error
	return ids, err
}

func (r *commentRepository) GetCommentByID(commentID int) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.Where("id = ?", commentID).First(&comment).Error
	return &comment, err
}

func (r *commentRepository) ListUserCommentsGet(userID int, page, pageSize int, status int) ([]response.CommentInfo, int64, error) {
	db := r.db.Model(&entity.Comment{}).
		Joins("JOIN articles ON comments.article_id = articles.id").
		Where("articles.user_id = ? AND comments.status = ?", userID, status)

	return r.listCommentsWithBase(db, page, pageSize)
}

func (r *commentRepository) ListUserCommentsTake(userID int, page, pageSize int) ([]response.CommentInfo, int64, error) {
	db := r.db.Model(&entity.Comment{}).
		Where("comments.user_id = ? AND comments.status = 1", userID)

	return r.listCommentsWithBase(db, page, pageSize)
}

func (r *commentRepository) listCommentsWithBase(db *gorm.DB, page, pageSize int) ([]response.CommentInfo, int64, error) {
	total, err := r.countComments(db)
	if err != nil {
		return nil, 0, err
	}

	comments, err := r.findCommentsWithPage(db, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return r.toDTOList(comments), total, nil
}

func (r *commentRepository) countComments(db *gorm.DB) (int64, error) {
	var total int64
	err := db.Count(&total).Error
	return total, err
}

func (r *commentRepository) findCommentsWithPage(db *gorm.DB, page, pageSize int) ([]entity.Comment, error) {
	var comments []entity.Comment

	err := db.Preload("User").
		Preload("Article").
		Preload("ParentComment.User").
		Order("created_at desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&comments).Error

	return comments, err
}

func (r *commentRepository) toDTOList(comments []entity.Comment) []response.CommentInfo {
	res := make([]response.CommentInfo, len(comments))
	for i := range comments {
		res[i] = r.toDTO(&comments[i])
	}
	return res
}

func (r *commentRepository) toDTO(c *entity.Comment) response.CommentInfo {
	dto := response.CommentInfo{
		ID:           c.ID,
		ArticleID:    c.ArticleID,
		UserID:       c.UserID,
		Username:     c.User.Username,
		Avatar:       c.User.Avatar,
		ArticleTitle: c.Article.Title,
		Content:      c.Content,
		ParentID:     c.ParentID,
		CreatedAt:    c.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if c.ParentComment != nil {
		dto.ParentUsername = c.ParentComment.User.Username
	}
	return dto
}
