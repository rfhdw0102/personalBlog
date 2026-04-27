package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const changedKey = "article:changed"

type redisRepository struct {
	rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) RedisRepository {
	return &redisRepository{rdb: rdb}
}

// SetKey 存入私钥
func (r *redisRepository) SetKey(ctx context.Context, key string, value string, num time.Duration) error {
	return r.rdb.Set(ctx, key, value, num).Err()
}

// GetKey 获取私钥
func (r *redisRepository) GetKey(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

// DelKey 删除
func (r *redisRepository) DelKey(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

func (r *redisRepository) XAdd(ctx context.Context, stream string, values map[string]interface{}) (string, error) {
	return r.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: values,
		MaxLen: 1000,
		Approx: true,
	}).Result()
}

func (r *redisRepository) XGroupCreateMkStream(ctx context.Context, stream, group, start string) error {
	return r.rdb.XGroupCreateMkStream(ctx, stream, group, start).Err()
}

func (r *redisRepository) XReadGroup(ctx context.Context, args *redis.XReadGroupArgs) ([]redis.XStream, error) {
	return r.rdb.XReadGroup(ctx, args).Result()
}

func (r *redisRepository) XAck(ctx context.Context, stream, group string, ids ...string) error {
	return r.rdb.XAck(ctx, stream, group, ids...).Err()
}

func (r *redisRepository) IncrView(ctx context.Context, articleID int) error {
	key := fmt.Sprintf("article:%d:view", articleID)
	pipe := r.rdb.TxPipeline()

	_ = pipe.Incr(ctx, key)
	pipe.SAdd(ctx, changedKey, articleID)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *redisRepository) IncrLike(ctx context.Context, articleID int) error {
	key := fmt.Sprintf("article:%d:like", articleID)
	pipe := r.rdb.TxPipeline()

	_ = pipe.Incr(ctx, key)
	pipe.SAdd(ctx, changedKey, articleID)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *redisRepository) DecrLike(ctx context.Context, articleID int) error {
	key := fmt.Sprintf("article:%d:like", articleID)
	_, err := r.rdb.Decr(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisRepository) GetStats(ctx context.Context, articleID int) (int64, int64, error) {
	viewKey := fmt.Sprintf("article:%d:view", articleID)
	likeKey := fmt.Sprintf("article:%d:like", articleID)

	view, err := r.rdb.Get(ctx, viewKey).Int64()
	if errors.Is(err, redis.Nil) {
		view = 0
	} else if err != nil {
		return 0, 0, err
	}

	like, err := r.rdb.Get(ctx, likeKey).Int64()
	if errors.Is(err, redis.Nil) {
		view = 0
	} else if err != nil {
		return 0, 0, err
	}

	return view, like, nil
}

// GetDiff 获取变更列表
func (r *redisRepository) GetDiff(ctx context.Context) ([]string, error) {
	ids, err := r.rdb.SMembers(ctx, changedKey).Result()
	if err != nil {
		fmt.Println("获取变更列表失败:", err)
		return []string{}, err
	}
	return ids, nil
}
