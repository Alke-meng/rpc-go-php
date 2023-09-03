package asynq

import (
	"ccgo/dao/mysql"
	"ccgo/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	TypeCrmImportDelivery  = "crmImport:deliver"
	TypeCrmRecycleDelivery = "crmRecycle:deliver"
)

type CrmImportDeliveryPayload struct {
	TraceID string
	Data    map[string]any
}

type CrmRecycleDeliveryPayload struct {
	TraceID string
	Data    map[string]string
}

func NewCrmImportDeliveryTask(traceID string, data map[string]any) (*asynq.Task, error) {
	payload, err := json.Marshal(CrmImportDeliveryPayload{TraceID: traceID, Data: data})
	if err != nil {
		logger.CCgoTaskLogger(traceID, fmt.Sprintf("NewCrmImportDeliveryTask error:%v", err))
		return nil, err
	}
	return asynq.NewTask(TypeCrmImportDelivery, payload), nil
}

func HandleCrmImportDeliveryTask(ctx context.Context, t *asynq.Task) error {
	start := time.Now()
	var p CrmImportDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		logger.CCgoTaskLogger(p.TraceID, fmt.Sprintf("HandleCrmImportDeliveryTask failed:%v: %w", err, asynq.SkipRetry))
		return fmt.Errorf("HandleCrmImportDeliveryTask failed: %v: %w", err, asynq.SkipRetry)
	}
	//逻辑处理start...
	res := mysql.CallerListADD(p.Data, p.TraceID)
	if res != nil {
		logger.CCgoTaskLogger(p.TraceID, fmt.Sprintf("HandleCrmImportDeliveryTask action fail:%v:", res))
		return errors.New(fmt.Sprintf("HandleCrmImportDeliveryTask action fail:%v:", res))
	}

	logger.CCgoTaskCostLogger(p.TraceID, fmt.Sprintf("CallerListADD task suceess"), time.Now().Sub((start)))
	return nil
}

func NewCrmRecycleDeliveryTask(traceID string, data map[string]string) (*asynq.Task, error) {
	payload, err := json.Marshal(CrmRecycleDeliveryPayload{TraceID: traceID, Data: data})
	if err != nil {
		logger.CCgoTaskLogger(traceID, fmt.Sprintf("NewCrmRecycleDeliveryTask error:%v", err))
		return nil, err
	}
	return asynq.NewTask(TypeCrmRecycleDelivery, payload), nil
}

func HandleCrmRecycleDeliveryTask(ctx context.Context, t *asynq.Task) error {
	start := time.Now()
	var p CrmRecycleDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		logger.CCgoTaskLogger(p.TraceID, fmt.Sprintf("HandleCrmRecycleDeliveryTask failed:%v: %w", err, asynq.SkipRetry))
		return fmt.Errorf("HandleCrmRecycleDeliveryTask failed: %v: %w", err, asynq.SkipRetry)
	}
	//逻辑处理start...
	res := mysql.CrmRecycle(p.Data, p.TraceID)
	if res != nil {
		logger.CCgoTaskLogger(p.TraceID, fmt.Sprintf("HandleCrmRecycleDeliveryTask action fail:%v:", res))
		return errors.New(fmt.Sprintf("HandleCrmRecycleDeliveryTask action fail:%v:", res))
	}

	err := os.Remove(viper.GetString("tmp_file_path") + "/" + p.Data["data"])
	if err != nil {
		logger.CCgoTaskLogger(p.TraceID, fmt.Sprintf("file %s remove fail:%v:", p.Data, err))
	}

	logger.CCgoTaskCostLogger(p.TraceID, fmt.Sprintf("crm recycle task suceess, file is %s", p.Data["data"]), time.Now().Sub((start)))
	return nil
}
