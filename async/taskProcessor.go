package async

import (
	"context"

	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/elastic/go-elasticsearch/v8"

	"github.com/hibiken/asynq"
	mail "github.com/xhit/go-simple-mail/v2"
)

var (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProccessor interface {
	Start() error
	HandleEmailDeliveryTask(ctx context.Context, task *asynq.Task) error
	HandleSyncTagTask(ctx context.Context, task *asynq.Task) error
	HandleSyncProductTask(ctx context.Context, task *asynq.Task) error
	HandleTaskNotification(
		ctx context.Context,
		task *asynq.Task,
	) error
}

type RedisTaskProccessor struct {
	server        *asynq.Server
	store         db.Store
	mailClient    *mail.SMTPClient
	elasticClient elasticsearch.Client
	mailSender    string
}

func NewRedisTaskProccessor(redisOpt asynq.RedisClientOpt, store db.Store, mailCient *mail.SMTPClient, sender string, elasticClient elasticsearch.Client) TaskProccessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			Concurrency: 4,
		},
	)

	return &RedisTaskProccessor{
		server:        server,
		store:         store,
		mailClient:    mailCient,
		mailSender:    sender,
		elasticClient: elasticClient,
	}
}

func (proccessor *RedisTaskProccessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailDeliver, proccessor.HandleEmailDeliveryTask)
	mux.HandleFunc(TypeNotification, proccessor.HandleTaskNotification)
	mux.HandleFunc(TypeSyncProduct, proccessor.HandleSyncProductTask)
	mux.HandleFunc(TypeSyncTag, proccessor.HandleSyncTagTask)
	return proccessor.server.Start(mux)
}
