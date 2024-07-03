package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"micros/internal/common"
	"net/http"
	"runtime/debug"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logFilePath = "./logs/access"
)

	func logHandler()  gin.HandlerFunc  {
		// 实例化
		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)

		// 设置 rotatelogs
		logWriter, _ := rotatelogs.New(
			logFilePath+".%Y%m%d.log",
			rotatelogs.WithLinkName(logFilePath),
			rotatelogs.WithMaxAge(30*24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)

		writeMap := lfshook.WriterMap{
			logrus.InfoLevel:  logWriter,
			logrus.FatalLevel: logWriter,
			logrus.DebugLevel: logWriter,
			logrus.WarnLevel:  logWriter,
			logrus.ErrorLevel: logWriter,
			logrus.PanicLevel: logWriter,
		}
		logger.AddHook(lfshook.NewHook(writeMap, &nested.Formatter{
			//HideKeys:        true,
			TimestampFormat: common.TimeFormat,
			FieldsOrder:     []string{"time", "client_ip", "status_code", "latency_time", "req_method", "req_uri"},
		}))

		return func(ctx *gin.Context) {
			data, err := ctx.GetRawData()
			if err != nil {
				logger.Error("请求参数解析失败：", err)
				common.Error(common.ErrorInvalidArgument, err)
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
					logger.WithFields(logrus.Fields{
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
				logger.WithFields(logrus.Fields{
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