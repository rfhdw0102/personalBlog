package job

import (
	"blogs/internal/repository"
	"blogs/pkg/logger"
	"context"
	"go.uber.org/zap"
	"strconv"
)

type ArticleSync struct {
	redis       repository.RedisRepository
	articleRepo repository.ArticleRepository
}

func NewArticleSync(redisRepo repository.RedisRepository, articleRepo repository.ArticleRepository) *ArticleSync {
	return &ArticleSync{
		redis:       redisRepo,
		articleRepo: articleRepo,
	}
}

func (s *ArticleSync) SyncChangedToMySQL(ctx context.Context) {
	ids, err := s.redis.GetDiff(ctx)
	if err != nil {
		logger.Error("获取缓存差异失败", zap.Error(err))
	}
	if len(ids) == 0 {
		return
	}
	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Warn("类型转换失败", zap.Error(err))
			return
		}
		view, like, err := s.redis.GetStats(ctx, id)
		if err != nil {
			logger.Error("缓存获取浏览量、点赞量失败", zap.Error(err))
			return
		}
		err = s.articleRepo.UpdateViewAndLike(id, int(view), int(like))
		if err != nil {
			logger.Error("更新失败浏览量、点赞量:", zap.Int("id:", id))
			return
		}
	}
	if err = s.redis.DelKey(ctx, "article:changed"); err != nil {
		logger.Error("清空变更集合失败", zap.Error(err))
	}
}
