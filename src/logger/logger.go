package logger

import (
	"log"
	"os"
)

const (
	calldepth = 3
)

var L Logger = newLogger()

func newLogger() Logger {
	logFile, _ := os.Create("../debug.log")
	return &DefaultLogger{log.New(logFile, "", log.LstdFlags|log.Lshortfile)}
}

type Logger interface {
	Sql(v ...interface{})

	Http(v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

func SetLogger(logger Logger) {
	L = logger
}

func Debug(v ...interface{}) {
	L.Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	L.Debugf(format, v...)
}

func Info(v ...interface{}) {
	L.Info(v...)
}
func Infof(format string, v ...interface{}) {
	L.Infof(format, v...)
}

func Warn(v ...interface{}) {
	L.Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	L.Warnf(format, v...)
}

func Error(v ...interface{}) {
	L.Error(v...)
}
func Errorf(format string, v ...interface{}) {
	L.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	L.Fatal(v...)
}
func Fatalf(format string, v ...interface{}) {
	L.Fatalf(format, v...)
}

func Panic(v ...interface{}) {
	L.Panic(v...)
}
func Panicf(format string, v ...interface{}) {
	L.Panicf(format, v...)
}
