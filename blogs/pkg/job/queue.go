package job

import (
	"blogs/internal/model/entity"
	"blogs/internal/model/event"
	"blogs/internal/repository"
	"blogs/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

const (
	StreamKey    = "notifications_stream"
	GroupName    = "notification_group"
	ConsumerName = "notification_consumer"
)

type NotificationConsumer struct {
	redisRepo        repository.RedisRepository
	notificationRepo repository.NotificationRepository
}

func NewNotificationConsumer(redisRepo repository.RedisRepository, notificationRepo repository.NotificationRepository) *NotificationConsumer {
	return &NotificationConsumer{
		redisRepo:        redisRepo,
		notificationRepo: notificationRepo,
	}
}

// Start 开启消费者
func (c *NotificationConsumer) Start(ctx context.Context) {
	// 创建消费者组, 从最开始消费
	err := c.redisRepo.XGroupCreateMkStream(ctx, StreamKey, GroupName, "0")
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		logger.Error("创建消费者组失败", zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 从未ACK开始读取消息
			streams, err := c.redisRepo.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    GroupName,
				Consumer: ConsumerName,
				Streams:  []string{StreamKey, ">"},
				Count:    10,
				Block:    time.Second * 5,
			})

			if err != nil {
				if !errors.Is(err, redis.Nil) {
					logger.Error("从streams流中读取失败", zap.Error(err))
				}
				continue
			}

			for _, stream := range streams {
				for _, msg := range stream.Messages {
					c.processMessage(ctx, msg)
				}
			}
		}
	}
}

// processMessage 处理消息
func (c *NotificationConsumer) processMessage(ctx context.Context, msg redis.XMessage) {
	eventJSON, ok := msg.Values["event"].(string)
	if !ok {
		logger.Error("Invalid message format", zap.Any("msg", msg))
		return
	}

	var ev event.NotificationEvent
	if err := json.Unmarshal([]byte(eventJSON), &ev); err != nil {
		logger.Error("反序列化消息失败", zap.Error(err))
		return
	}

	if ev.RecipientID == ev.ActorID {
		err := c.redisRepo.XAck(ctx, StreamKey, GroupName, msg.ID)
		if err != nil {
			logger.Error("ACK 消息失败:", zap.String("msgID:", msg.ID))
		}
		return
	}

	notification := &entity.Notification{
		RecipientID: ev.RecipientID,
		ActorID:     ev.ActorID,
		ArticleID:   ev.ArticleID,
		CommentID:   ev.CommentID,
		Type:        entity.NotificationType(ev.Type),
		IsRead:      false,
	}

	if err := c.notificationRepo.Create(notification); err != nil {
		logger.Error("写入mysql消息失败", zap.Error(err))
		return
	}

	err := c.redisRepo.XAck(ctx, StreamKey, GroupName, msg.ID)
	if err != nil {
		logger.Error("ACK 消息失败:", zap.String("msgID:", msg.ID))
		return
	}
}
