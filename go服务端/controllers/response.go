package controllers

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseInfo(code ResCode, data interface{}) (rd *ResponseData) {
	rd = &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: data,
	}
	return
}
