package logger

import (
	"bytes"
	"fmt"
	"micros/internal/common"
	"micros/internal/config"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

//Log 日志对象
var Log *logrus.Logger

//MyFormatter 日志格式
type MyFormatter struct {}

// Format 格式化
func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error){
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format(common.TimeFormat)
	var newLog string
	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n",
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else{
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

// Init 创建日志对象
func Init(c *config.Log) {
		// 实例化
		Log = logrus.New()
		Log.SetLevel(logrus.DebugLevel)

		// 设置 rotatelogs
		logWriter, _ := rotatelogs.New(
			c.Path +".%Y%m%d.log",
			rotatelogs.WithLinkName(c.Path),
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
		Log.AddHook(lfshook.NewHook(writeMap, &logrus.JSONFormatter{}))
		return
}