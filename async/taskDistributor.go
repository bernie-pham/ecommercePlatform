package async

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendMail(
		ctx context.Context,
		payload *EmailDeliveryPayload,
		opt ...asynq.Option) error
	DistributeTaskNotification(
		ctx context.Context,
		payload *NotificationPayload,
		opt ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
