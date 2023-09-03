package asynqClient

import (
	"ccgo/settings"
	"github.com/hibiken/asynq"
)

var AsynqClient *asynq.Client

func Init(cfg *settings.RedisQueueConfig) {
	AsynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

func Close() {
	_ = AsynqClient.Close()
}
