package common

import (
	"time"
)

type response struct {
	Time    time.Time   `json:"time"`
	Code    int         `json:"code"`
	Stat    int         `json:"stat"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// Success 成功
func Success(v interface{}) interface{} {
	return response{Time: time.Now(), Stat: 1, Code: 0, Message: "ok", Data: v}
}

// Error 失败
func Error(bizError *BizError, err error) interface{} {
	return response{Time: time.Now(), Stat: 0, Code: bizError.Code, Message: bizError.Message, Data: err.Error()}
}
