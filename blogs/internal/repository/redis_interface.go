package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	SetKey(ctx context.Context, key string, value string, num time.Duration) error
	GetKey(ctx context.Context, key string) (string, error)
	DelKey(ctx context.Context, key string) error

	XAdd(ctx context.Context, stream string, values map[string]interface{}) (string, error)
	XGroupCreateMkStream(ctx context.Context, stream, group, start string) error
	XReadGroup(ctx context.Context, args *redis.XReadGroupArgs) ([]redis.XStream, error)
	XAck(ctx context.Context, stream, group string, ids ...string) error

	IncrView(ctx context.Context, articleID int) error
	IncrLike(ctx context.Context, articleID int) error
	DecrLike(ctx context.Context, articleID int) error
	GetStats(ctx context.Context, articleID int) (view int64, like int64, err error)
	GetDiff(ctx context.Context) ([]string, error)
}
