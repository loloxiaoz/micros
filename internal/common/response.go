package common

import "time"

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
func Error(err error) interface{} {
	//ret := response{Stat: 0, Code: err.Code, Message: e.Message, Data: e.Info}
	//	return ret
	return ""
}
