package routes

import (
	"ccgo/asynq"
	"ccgo/controllers"
	"ccgo/dao/asynqClient"
	"ccgo/dao/mysql"
	"ccgo/logger"
	"encoding/json"
	"fmt"
	"time"
)

func (s *CCgo) ImportCrm(inputSting string, r *string) error {
	//start
	start := time.Now()
	data := RequestFromPhp(inputSting)
	traceID := data["traceID"].(string)

	//TODO
	// 1、mysql数据处理
	tmp, err := mysql.ImportCrm(data, traceID, start)
	res, _ := json.Marshal(tmp)
	if err != nil {
		logger.CCgoActionLogger(traceID, fmt.Sprintf("ImportCrm Exec fail:%v", err))
		ResponseToPhp(controllers.CodeFail, string(res), traceID, r, start)
		return nil
	}

	// 呼叫名单
	if tmp.Success > 0 {
		// 2、任务投递
		dataTmp := data["data"].(map[string]any)
		defaultData := dataTmp["default"].(map[string]any)
		taskData := map[string]any{
			"number":      tmp.PhoneData,
			"customer_id": defaultData["customer_id"].(string),
		}
		asynq.CrmImportDeliveryTaskAdd(asynqClient.AsynqClient, traceID, taskData)
	}

	//end
	ResponseToPhp(controllers.CodeSuccess, string(res), traceID, r, start)
	return nil
}

func (s *CCgo) DeleteCrm(inputSting string, r *string) error {
	//start
	start := time.Now()
	data := RequestFromPhp(inputSting)
	traceID := data["traceID"].(string)

	//TODO
	// 1、mysql数据处理
	tmp, err := mysql.DeleteCrm(data, traceID)
	res, _ := json.Marshal(tmp)
	if err != nil {
		logger.CCgoActionLogger(traceID, fmt.Sprintf("DeleteCrm Exec fail:%v", err))
		ResponseToPhp(controllers.CodeFail, string(res), traceID, r, start)
		return nil
	}

	// 数据回收
	dataTmp := data["data"].(map[string]any)
	if tmp.Total > 0 && dataTmp["recycle"] == true {
		defaultData := dataTmp["default"].(map[string]any)
		table := "test.tbl_crm_recycle_" + defaultData["customer_id"].(string)
		columns := defaultData["columns"].(string)
		// 2、任务投递
		for i := 1; i <= tmp.GoNum; i++ {
			tasKData := map[string]string{
				"table": table,
				"field": columns,
				"data":  tmp.OutFile[i],
			}
			asynq.CrmRecycleDeliveryTaskAdd(asynqClient.AsynqClient, traceID, tasKData)
		}
	}

	//end
	ResponseToPhp(controllers.CodeSuccess, string(res), traceID, r, start)
	return nil
}
