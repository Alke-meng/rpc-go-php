package controllers

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeFail
	CodeInvalidParam
	CodeServerBusy
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:      "success",
	CodeFail:         "action fail",
	CodeInvalidParam: "invalid param error",
	CodeServerBusy:   "system error",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}

	return msg
}
