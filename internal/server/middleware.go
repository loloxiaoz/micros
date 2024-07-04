package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"

	"micros/internal/common"
	"micros/internal/log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func logHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := ctx.GetRawData()
		if err != nil {
			log.Logger().Error("请求参数解析失败：", err)
			ctx.JSON(http.StatusBadRequest, common.Error(common.BindJSONError, err))
			return
		}
		//Body数据只能读取一次，需要需要写回去
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

		startTime := time.Now()
		defer func() {
			//状态码
			statusCode := ctx.Writer.Status()
			//结束时间
			endTime := time.Now()
			// 执行时间
			latencyTime := endTime.Sub(startTime)

			if err := recover(); err != nil {
				log.Logger().WithFields(logrus.Fields{
					"client_ip":    ctx.ClientIP(),
					"status_code":  statusCode,
					"latency_time": latencyTime.Milliseconds(),
					"req_method":   ctx.Request.Method,
					"req_uri":      ctx.Request.RequestURI,
				}).Error(fmt.Sprintf("panic: %s", err), "\n", string(data), "\n", string(debug.Stack()))

				ctx.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}

			// 日志格式
			log.Logger().WithFields(logrus.Fields{
				"client_ip":    ctx.ClientIP(),
				"status_code":  statusCode,
				"latency_time": latencyTime.Milliseconds(),
				"req_method":   ctx.Request.Method,
				"req_uri":      ctx.Request.RequestURI,
			}).Info(string(data))
		}()

		//处理请求
		ctx.Next()
	}
}
