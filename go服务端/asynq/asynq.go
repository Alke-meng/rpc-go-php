package asynq

import (
	"ccgo/settings"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

func SetUp(rqc *settings.RedisQueueConfig) (err error) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     rqc.Addr,
			Password: rqc.Password,
			DB:       rqc.DB,
		},
		asynq.Config{
			// 每个进程并发执行的worker数量
			Concurrency: viper.GetInt("asynq.concurrency"),
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"cc-crm-import":  6,
				"cc-crm-recycle": 3,
				"low":            1,
			},
			// See the godoc for other configuration options
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeCrmImportDelivery, HandleCrmImportDeliveryTask)
	mux.HandleFunc(TypeCrmRecycleDelivery, HandleCrmRecycleDeliveryTask)

	err = srv.Run(mux)

	return err
}
