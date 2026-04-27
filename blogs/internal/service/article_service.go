package service

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"
	"blogs/internal/repository"
	"blogs/pkg/logger"
	"blogs/pkg/utils"
	"context"
	"os"
	"strings"

	"go.uber.org/zap"
)

type articleService struct {
	articleRepo  repository.ArticleRepository
	categoryRepo repository.CategoryRepository
	tagRepo      repository.TagRepository
	commentRepo  repository.CommentRepository
	redis        repository.RedisRepository
}

func NewArticleService(articleRepo repository.ArticleRepository, categoryRepo repository.CategoryRepository, tagRepo repository.TagRepository, commentRepo repository.CommentRepository, redis repository.RedisRepository) ArticleService {
	return &articleService{articleRepo: articleRepo, categoryRepo: categoryRepo, tagRepo: tagRepo, commentRepo: commentRepo, redis: redis}
}

func (s *articleService) ListWithContent(req request.ArticleListQuery) ([]response.ArticleInfo, int64, error) {
	return s.articleRepo.List(req)
}

func (s *articleService) Create(article *entity.Article, tagIDs []int, tagNames []string) error {
	return s.articleRepo.Create(article, tagIDs, tagNames)
}

func (s *articleService) Update(article *entity.Article, tagIDs []int, tagNames []string) error {
	oldArticle, err := s.articleRepo.GetByID(article.ID)
	if err == nil {
		if oldArticle.CoverImage != "" && oldArticle.CoverImage != article.CoverImage {
			oldPath := strings.TrimPrefix(oldArticle.CoverImage, "/")
			err = os.Remove(oldPath)
			if err != nil {
				logger.Warn("删除封面失败", zap.Error(err))
			}
		}

		oldImages := utils.ExtractImages(oldArticle.Content)
		newImages := utils.ExtractImages(article.Content)

		// 找出在旧内容中但不在新内容中的图片
		newImagesMap := make(map[string]bool)
		for _, img := range newImages {
			newImagesMap[img] = true
		}

		for _, oldImg := range oldImages {
			if !newImagesMap[oldImg] {
				path := strings.TrimPrefix(oldImg, "/")
				if err := os.Remove(path); err != nil {
					logger.Warn("清理冗余内容图片失败", zap.String("url", oldImg), zap.Error(err))
				}
			}
		}
	}

	return s.articleRepo.Update(article, tagIDs, tagNames)
}

func (s *articleService) Delete(id int) error {
	// 获取文章信息以获取封面路径
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除封面文件
	if article.CoverImage != "" {
		oldPath := strings.TrimPrefix(article.CoverImage, "/")
		if err = os.Remove(oldPath); err != nil {
			logger.Error("删除封面文件失败", zap.Error(err))
		}
	}

	return s.articleRepo.Delete(id)
}

func (s *articleService) GetByID(id int) (*response.ArticleInfo, error) {
	return s.articleRepo.GetByID(id)
}

func (s *articleService) IncrementViewCount(id int) error {
	err := s.redis.IncrView(context.Background(), id)
	return err
}

func (s *articleService) LikeArticle(articleID int, userID int) error {
	err := s.articleRepo.LikeArticle(articleID, userID)
	if err != nil {
		return err
	}
	err = s.redis.IncrLike(context.Background(), articleID)
	return err
}

func (s *articleService) UnlikeArticle(articleID int, userID int) error {
	err := s.articleRepo.UnlikeArticle(articleID, userID)
	if err != nil {
		return err
	}
	err = s.redis.DecrLike(context.Background(), articleID)
	return err
}

func (s *articleService) IsLiked(articleID int, userID int) (bool, error) {
	return s.articleRepo.IsLiked(articleID, userID)
}

func (s *articleService) GetStats(articleID int) (int, int, error) {
	view, like, err := s.redis.GetStats(context.Background(), articleID)
	if err != nil {
		logger.Error("缓存获取点赞量，浏览量失败", zap.Error(err))
		return 0, 0, err
	}
	return int(view), int(like), nil
}
