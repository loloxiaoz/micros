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
func Error(bizError string, err error) interface{} {
	bizCode := ErrCodeMap[bizError]
	return response{Stat: 0, Code: bizCode, Message: bizError, Data: err.Error()}
}
