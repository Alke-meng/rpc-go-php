package asynq

import (
	"ccgo/logger"
	"fmt"
	"github.com/hibiken/asynq"
	"time"
)

func CrmImportDeliveryTaskAdd(ayq *asynq.Client, traceID string, data map[string]any) {
	task, err := NewCrmImportDeliveryTask(traceID, data)
	if err != nil {
		logger.CCgoTaskLogger(traceID, fmt.Sprintf("could not create task:%v", err))
	}
	// 投递任务(延时5秒、重试次数3)
	info, err := ayq.Enqueue(
		task,
		asynq.Queue("cc-crm-import"),
		asynq.ProcessIn(5*time.Second),
		asynq.MaxRetry(3),
	)

	if err != nil {
		logger.CCgoTaskLogger(traceID, fmt.Sprintf("could not enqueue task:%v", err))
	}

	logger.CCgoTaskLogger(traceID, fmt.Sprintf("enqueued task: id=%s queue=%s", info.ID, info.Queue))
}

func CrmRecycleDeliveryTaskAdd(ayq *asynq.Client, traceID string, data map[string]string) {
	task, err := NewCrmRecycleDeliveryTask(traceID, data)
	if err != nil {
		logger.CCgoTaskLogger(traceID, fmt.Sprintf("could not create task:%v", err))
	}
	// 投递任务(延时5秒、重试次数3)
	info, err := ayq.Enqueue(
		task,
		asynq.Queue("cc-crm-recycle"),
		asynq.ProcessIn(5*time.Second),
		asynq.MaxRetry(3),
	)

	if err != nil {
		logger.CCgoTaskLogger(traceID, fmt.Sprintf("could not enqueue task:%v", err))
	}

	logger.CCgoTaskLogger(traceID, fmt.Sprintf("enqueued task: id=%s queue=%s", info.ID, info.Queue))
}
