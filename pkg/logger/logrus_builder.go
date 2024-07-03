package logger

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)
func filter(msg string, r ...string) string {
	replacer := strings.NewReplacer("\t", "", "\r", "", "\n", "")
	if len(r) > 0 {
		replacer = strings.NewReplacer("\t", r[0], "\r", r[0], "\n", r[0])
	}
	return replacer.Replace(msg)
}

type logrusBuilder struct {
	LogrusI *logrus.Logger
}

// NewLogrusBuilder 新建日志builder
func NewLogrusBuilder(logger *logrus.Logger) *logrusBuilder {
	return &logrusBuilder{LogrusI: logger}
}

//LoggerX 实现LoggerX
func (l *logrusBuilder) LoggerX(ctx context.Context, lvl string, tag string, args interface{}, v ...interface{}) {
	if tag == "" {
		tag = "NoTagError"
	}

	tag = filter(tag)
	_, message := l.Build(ctx, args, v...)

	field := l.LogrusI.WithFields(logrus.Fields{
		"tag": tag,
	})
	switch lvl {
	case "DEBUG":
		field.Debug(message)
	case "TRACE":
		field.Trace(message)
	case "INFO":
		field.Info(message)
	case "WARNING":
		field.Warn(message)
	case "ERROR":
		field.Error(message)
	case "FATAL":
		field.Panic(message)
	}
}

// Build interface Builder function implemented
func (l *logrusBuilder) Build(ctx context.Context, args interface{}, v ...interface{}) (position string, message string) {

	switch t := args.(type) {
	case *StackErr:
		message = t.Info
	case error:
		message = t.Error()
	case string:
		if len(v) > 0 {
			message = fmt.Sprintf(t, v...)
		} else {
			message = t
		}
	default:
		message = fmt.Sprint(t)
	}
	message = filter(message)
	return
}