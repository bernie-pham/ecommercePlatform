package async

import (
	"context"

	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
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
	// DistributeSyncProductDataTask(
	// 	ctx context.Context,
	// 	payload *SyncAllPayload,
	// 	opt ...asynq.Option,
	// ) error
	DistributeSyncAllTagDataTask(
		ctx context.Context,
		opt ...asynq.Option,
	) error
	DistributeSyncAllProductDataTask(
		ctx context.Context,
		opt ...asynq.Option,
	) error
	DistributeSyncNewDataTask(
		ctx context.Context,
		payload *SyncNewPayload,
		opt ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
	store  db.Store
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt, store db.Store) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
		store:  store,
	}
}
