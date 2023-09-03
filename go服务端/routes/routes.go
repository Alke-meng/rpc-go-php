package routes

import (
	"ccgo/controllers"
	"ccgo/logger"
	"ccgo/tool"
	"encoding/json"
	"fmt"
	"net/rpc"
	"time"
)

type CCgo struct{}

func (s *CCgo) Hi(name string, r *string) error {
	*r = fmt.Sprintf("Hello, %s!", name)
	return nil
}

func Setup() {
	_ = rpc.Register(new(CCgo))
}

func ResponseToPhp(code controllers.ResCode, data, traceID string, r *string, startTime time.Time) {
	jsonTmp := controllers.ResponseInfo(code, data)
	res, _ := json.Marshal(jsonTmp)
	*r = string(res)
	// 记录响应日志
	logger.CCgoResponseLogger(code, tool.GetFunName(2), string(res), traceID, time.Now().Sub((startTime)))
	return
}

func RequestFromPhp(inputSting string) (mapData map[string]interface{}) {
	// 记录请求日志
	logger.CCgoRequestLogger(tool.GetFunName(2), inputSting)
	jsonData := []byte(inputSting)
	json.Unmarshal(jsonData, &mapData)
	return
}
