package repository

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"
	"blogs/internal/model/event"
	"blogs/pkg/logger"
	"context"
	"encoding/json"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type articleRepository struct {
	db        *gorm.DB
	redisRepo RedisRepository
}

func NewArticleRepository(db *gorm.DB, redisRepo RedisRepository) ArticleRepository {
	return &articleRepository{db: db, redisRepo: redisRepo}
}

func (r *articleRepository) List(req request.ArticleListQuery) ([]response.ArticleInfo, int64, error) {
	db := r.applyArticleFilters(r.db.Model(&entity.Article{}), req)

	total, err := r.countArticles(db)
	if err != nil {
		return nil, 0, err
	}

	articles, err := r.findArticles(db, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	return r.toDTOList(articles), total, nil
}

func (r *articleRepository) Create(article *entity.Article, tagIDs []int, tagNames []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&article).Error; err != nil {
			return err
		}
		return r.syncArticleTags(tx, article.ID, tagIDs, tagNames, false)
	})
}

func (r *articleRepository) Update(article *entity.Article, tagIDs []int, tagNames []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(article).Error; err != nil {
			return err
		}
		return r.syncArticleTags(tx, article.ID, tagIDs, tagNames, true)
	})
}

func (r *articleRepository) processTags(tx *gorm.DB, tagIDs []int, tagNames []string) ([]int, error) {
	tagIDMap := make(map[int]bool)
	for _, id := range tagIDs {
		tagIDMap[id] = true
	}

	validNames := r.uniqueTagNames(tagNames)

	if len(validNames) > 0 {
		var existingTags []entity.Tag
		if err := tx.Where("name IN ?", validNames).Find(&existingTags).Error; err != nil {
			return nil, err
		}

		existingNameMap := r.buildExistingTagNameMap(existingTags, tagIDMap)
		for _, tag := range existingTags {
			tagIDMap[tag.ID] = true
		}

		if err := r.createMissingTags(tx, validNames, existingNameMap, tagIDMap); err != nil {
			return nil, err
		}
	}

	return r.collectTagIDs(tagIDMap), nil
}

func (r *articleRepository) Delete(id int) error {
	// 使用模型实例进行删除，以触发 GORM 的 BeforeDelete 钩子实现级联删除
	return r.db.Delete(&entity.Article{Base: entity.Base{ID: id}}).Error
}

func (r *articleRepository) GetByID(id int) (*response.ArticleInfo, error) {
	var article entity.Article
	err := r.preloadArticleRelations(r.db.Model(&entity.Article{})).
		Where("articles.id = ?", id).
		First(&article).Error
	if err != nil {
		return nil, err
	}

	return r.toDTO(&article), nil
}

func (r *articleRepository) LikeArticle(articleID int, userID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		like := entity.Like{
			ArticleID: articleID,
			UserID:    userID,
		}
		if err := r.createLike(tx, like); err != nil {
			return err
		}

		return r.publishLikeNotification(tx, articleID, userID)
	})
}

func (r *articleRepository) UnlikeArticle(articleID int, userID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("article_id = ? AND user_id = ?", articleID, userID).Delete(&entity.Like{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}
		return nil
	})
}

func (r *articleRepository) UpdateViewAndLike(articleID int, view int, like int) error {
	return r.db.Model(&entity.Article{}).
		Where("id = ?", articleID).
		Updates(map[string]interface{}{
			"view_count": view,
			"like_count": like,
		}).Error
}

func (r *articleRepository) IsLiked(articleID int, userID int) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Like{}).Where("article_id = ? AND user_id = ?", articleID, userID).Count(&count).Error
	return count > 0, err
}

func (r *articleRepository) GetAdjacent(articleID int, sort int) (*response.AdjacentArticles, error) {
	result := &response.AdjacentArticles{}

	if sort == 1 {
		var current entity.Article
		if err := r.db.Select("view_count").First(&current, articleID).Error; err != nil {
			return nil, err
		}

		var prev entity.Article
		err := r.db.Select("id", "title").
			Where("(view_count < ? OR (view_count = ? AND id < ?)) AND status = ?", current.ViewCount, current.ViewCount, articleID, "published").
			Order("view_count DESC, id DESC").
			First(&prev).Error
		if err == nil {
			result.Prev = &response.AdjacentArticle{ID: prev.ID, Title: prev.Title}
		}

		var next entity.Article
		err = r.db.Select("id", "title").
			Where("(view_count > ? OR (view_count = ? AND id > ?)) AND status = ?", current.ViewCount, current.ViewCount, articleID, "published").
			Order("view_count ASC, id ASC").
			First(&next).Error
		if err == nil {
			result.Next = &response.AdjacentArticle{ID: next.ID, Title: next.Title}
		}
	} else {
		var prev entity.Article
		err := r.db.Select("id", "title").
			Where("id < ? AND status = ?", articleID, "published").
			Order("id DESC").
			First(&prev).Error
		if err == nil {
			result.Prev = &response.AdjacentArticle{ID: prev.ID, Title: prev.Title}
		}

		var next entity.Article
		err = r.db.Select("id", "title").
			Where("id > ? AND status = ?", articleID, "published").
			Order("id ASC").
			First(&next).Error
		if err == nil {
			result.Next = &response.AdjacentArticle{ID: next.ID, Title: next.Title}
		}
	}

	return result, nil
}

func (r *articleRepository) toDTOList(articles []entity.Article) []response.ArticleInfo {
	res := make([]response.ArticleInfo, len(articles))
	for i, a := range articles {
		res[i] = *r.toDTO(&a)
	}
	return res
}

func (r *articleRepository) toDTO(a *entity.Article) *response.ArticleInfo {
	return &response.ArticleInfo{
		ID:           a.ID,
		Title:        a.Title,
		Content:      a.Content,
		Summary:      a.Summary,
		UserID:       a.UserID,
		Username:     a.User.Username,
		Avatar:       a.User.Avatar,
		CoverImage:   a.CoverImage,
		CategoryID:   a.CategoryID,
		CategoryName: a.Category.Name,
		Tags:         a.Tags,
		Status:       a.Status,
		ViewCount:    a.ViewCount,
		LikeCount:    a.LikeCount,
		UpdatedAt:    a.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (r *articleRepository) applyArticleFilters(db *gorm.DB, req request.ArticleListQuery) *gorm.DB {
	if req.Query != "" {
		db = db.Where("articles.title LIKE ?", "%"+req.Query+"%")
	}
	if req.Status != "" {
		db = db.Where("articles.status = ?", req.Status)
	}
	if req.CategoryID != 0 {
		db = db.Where("articles.category_id = ?", req.CategoryID)
	}
	if req.TagID != 0 {
		db = db.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", req.TagID)
	}
	if req.Sort == 0 {
		return db.Order("articles.created_at DESC")
	}
	return db.Order("articles.view_count DESC")
}

func (r *articleRepository) countArticles(db *gorm.DB) (int64, error) {
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *articleRepository) findArticles(db *gorm.DB, page, pageSize int) ([]entity.Article, error) {
	var articles []entity.Article
	err := r.preloadArticleRelations(db).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles).Error
	return articles, err
}

func (r *articleRepository) preloadArticleRelations(db *gorm.DB) *gorm.DB {
	return db.Preload("User").Preload("Category").Preload("Tags")
}

func (r *articleRepository) syncArticleTags(tx *gorm.DB, articleID int, tagIDs []int, tagNames []string, clearOld bool) error {
	if clearOld {
		if err := tx.Table("article_tags").Where("article_id = ?", articleID).Delete(nil).Error; err != nil {
			return err
		}
	}

	finalTagIDs, err := r.processTags(tx, tagIDs, tagNames)
	if err != nil {
		return err
	}
	if len(finalTagIDs) == 0 {
		return nil
	}
	articleTags := r.buildArticleTagRelations(articleID, finalTagIDs)
	return tx.Table("article_tags").Create(&articleTags).Error
}

func (r *articleRepository) buildArticleTagRelations(articleID int, tagIDs []int) []map[string]interface{} {
	articleTags := make([]map[string]interface{}, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		articleTags = append(articleTags, map[string]interface{}{
			"article_id": articleID,
			"tag_id":     tagID,
		})
	}
	return articleTags
}

func (r *articleRepository) uniqueTagNames(tagNames []string) []string {
	validNames := make([]string, 0, len(tagNames))
	seenNames := make(map[string]bool)
	for _, name := range tagNames {
		name = strings.TrimSpace(name)
		if name == "" || seenNames[name] {
			continue
		}
		validNames = append(validNames, name)
		seenNames[name] = true
	}
	return validNames
}

func (r *articleRepository) buildExistingTagNameMap(tags []entity.Tag, tagIDMap map[int]bool) map[string]int {
	existingNameMap := make(map[string]int, len(tags))
	for _, tag := range tags {
		tagIDMap[tag.ID] = true
		existingNameMap[tag.Name] = tag.ID
	}
	return existingNameMap
}

func (r *articleRepository) createMissingTags(tx *gorm.DB, tagNames []string, existingNameMap map[string]int, tagIDMap map[int]bool) error {
	for _, name := range tagNames {
		if _, exists := existingNameMap[name]; exists {
			continue
		}

		newTag := entity.Tag{Name: name}
		if err := tx.Create(&newTag).Error; err != nil {
			return err
		}
		tagIDMap[newTag.ID] = true
		existingNameMap[name] = newTag.ID
	}
	return nil
}

func (r *articleRepository) collectTagIDs(tagIDMap map[int]bool) []int {
	finalIDs := make([]int, 0, len(tagIDMap))
	for id := range tagIDMap {
		finalIDs = append(finalIDs, id)
	}
	return finalIDs
}

func (r *articleRepository) createLike(tx *gorm.DB, like entity.Like) error {
	err := tx.Create(&like).Error
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "Duplicate") {
		return nil
	}
	return err
}

func (r *articleRepository) publishLikeNotification(tx *gorm.DB, articleID, userID int) error {
	recipientID, err := r.findArticleOwnerID(tx, articleID)
	if err != nil {
		logger.Warn("获取文章作者失败，跳过点赞通知", zap.Error(err))
		return nil
	}

	return r.publishNotification(event.NotificationEvent{
		Type:        event.NotificationTypeLike,
		RecipientID: recipientID,
		ActorID:     userID,
		ArticleID:   articleID,
		CreatedAt:   time.Now(),
	})
}

func (r *articleRepository) findArticleOwnerID(tx *gorm.DB, articleID int) (int, error) {
	var article entity.Article
	if err := tx.Select("user_id").First(&article, articleID).Error; err != nil {
		return 0, err
	}
	return article.UserID, nil
}

func (r *articleRepository) publishNotification(ev event.NotificationEvent) error {
	evJSON, err := json.Marshal(ev)
	if err != nil {
		logger.Error("json序列化失败", zap.Error(err))
		return err
	}

	_, err = r.redisRepo.XAdd(context.Background(), "notifications_stream", map[string]interface{}{"event": string(evJSON)})
	if err != nil {
		logger.Warn("通知发送失败", zap.Error(err))
	}
	return nil
}
