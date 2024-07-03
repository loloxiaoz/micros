package logger

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

type Entry struct {
	caller bool
	Log    *logrus.Entry
}

func (e Entry) Debug(args ...interface{}) {
	if !e.caller {
		e.Log.Debug(args...)
		return
	}

	_, file, line, _ := runtime.Caller(1)
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Debug(args...)
}

func (e Entry) Debugf(format string, args ...interface{}) {
	if !e.caller {
		e.Log.Debugf(format, args...)
		return
	}

	_, file, line, _ := runtime.Caller(1)
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Debugf(format, args...)
}

func (e Entry) Info(args ...interface{}) {
	if !e.caller {
		e.Log.Info(args...)
		return
	}

	_, file, line, _ := runtime.Caller(1)
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Info(args...)
}

func (e Entry) Infof(format string, args ...interface{}) {
	if !e.caller {
		e.Log.Infof(format, args...)
		return
	}

	_, file, line, _ := runtime.Caller(1)
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Infof(format, args...)
}

func (e Entry) Error(args ...interface{}) {
	if !e.caller {
		e.Log.Error(args...)
		return
	}

	_, file, line, _ := runtime.Caller(1)
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Error(args...)
}

func (e Entry) Errorf(format string, args ...interface{}) {
	if !e.caller {
		e.Log.Errorf(format, args...)
		return
	}

	_, file, line, _ := runtime.Caller(1)
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Errorf(format, args...)
}

func (e Entry) Warn(args ...interface{}) {
	if !e.caller {
		e.Log.Warn(args...)
		return
	}

	_, file, line, _ := runtime.Caller(1)
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Warn(args...)
}

func (e Entry) Warnf(format string, args ...interface{}) {
	if !e.caller {
		e.Log.Warnf(format, args...)
		return
	}

	_, file, line, _ := runtime.Caller(1)
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Warnf(format, args...)
}

func (e Entry) Fatal(args ...interface{}) {
	if !e.caller {
		e.Log.Fatal(args...)
		return
	}

	_, file, line, _ := runtime.Caller(1)
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Fatal(args...)
}

func (e Entry) Fatalf(format string, args ...interface{}) {
	if !e.caller {
		e.Log.Fatalf(format, args...)
		return
	}

	_, file, line, _ := runtime.Caller(1)
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Fatalf(format, args...)
}

func (e Entry) Panic(args ...interface{}) {
	if !e.caller {
		e.Log.Panic(args...)
		return
	}

	_, file, line, _ := runtime.Caller(1)
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Panic(args...)
}

func (e Entry) Panicf(format string, args ...interface{}) {
	if !e.caller {
		e.Log.Panicf(format, args...)
		return
	}

	_, file, line, _ := runtime.Caller(1)
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Panicf(format, args...)
}

func (e Entry) WithField(key string, value interface{}) *Entry {
	return &Entry{
		Log:    e.Log.WithField(key, value),
		caller: e.caller,
	}
}

func (e Entry) WithFields(fields Fields) *Entry {
	return &Entry{
		Log:    e.Log.WithFields(logrus.Fields(fields)),
		caller: e.caller,
	}
}

func (e Entry) WithError(err error) *Entry {
	return &Entry{
		Log:    e.Log.WithField("error", err),
		caller: e.caller,
	}
}

// SetLevel 动态调整日志等级
func (e *Entry) SetLevel(level string) error {
	le, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	e.Log.Logger.SetLevel(le)
	return nil
}

// GetLevel 获取当前日志等级
func (e *Entry) GetLevel() string {
	level := e.Log.Logger.Level

	text, err := level.MarshalText()
	if err != nil {
		return "unknow"
	}

	return string(text)
}
