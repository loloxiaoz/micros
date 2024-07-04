package logger

import (
	"github.com/sirupsen/logrus"
)

const (
	//DefaultFileSize 默认日志文件大小
	DefaultFileSize = 100 * 1024 * 1024

	//OneMonth 一周
	OneMonth = 31 * 24 * 60 * 60

	//DefaultLogName 默认日志文件名
	DefaultLogName = "default.log"
)


const (
	logAge           = "maxAge"
	logName          = "name"
	logLevel         = "level"
	logCaller        = "caller"
	logWatcherEnable = "enable"
	logWatcherBySize = "watcherBySize"
)

type option interface {
	Get(key string) interface{}
}

type config struct {
	key   string
	value interface{}
}

func newCfg(key string, value interface{}) *config {
	return &config{
		key:   key,
		value: value,
	}
}

func (c config) Get(key string) interface{} {
	if key == c.key {
		return c.value
	}

	return nil
}

// 命名
func WithLogName(name string) option {
	return newCfg(logName, name)
}

func WithCaller(flag bool) option {
	return newCfg(logCaller, flag)
}

// 监控日志生命周期
func WithWatchEnable(enable bool) option {
	return newCfg(logWatcherEnable, enable)
}

func WithWatchLogsBySize(size int64) option {
	return newCfg(logWatcherBySize, size)
}

//调整日志级别
func WithLogLevel(level string) option {
	return newCfg(logLevel, level)
}

func WithMaxAge(age int64) option {
	return newCfg(logAge, age)
}

func findLevel(opts ...option) logrus.Level {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logLevel); nil != value {
			level, _ := logrus.ParseLevel(value.(string))
			return level
		}
	}

	return logrus.InfoLevel
}

/**
最大日期
*/
func findMaxAge(opts ...option) int64 {
	for _, opt := range opts {
		if nil == opt {
			continue
		}
		if value := opt.Get(logAge); nil != value {
			val := value.(int64)
			if val < OneMonth {
				val = OneMonth
			}
			return val
		}
	}

	return OneMonth
}

func findLogName(opts ...option) string {
	for _, opt := range opts {
		if nil == opt {
			continue
		}
		if value := opt.Get(logName); nil != value {
			return value.(string)
		}
	}

	return DefaultLogName
}

/**
监控日志状态及生命周期，默认开启
*/
func findWatcherEnable(opts ...option) bool {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logWatcherEnable); nil != value {
			return value.(bool)
		}
	}

	return true
}

func findWatchLogsBySize(opts ...option) int64 {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logWatcherBySize); nil != value {
			return value.(int64)
		}
	}

	return DefaultFileSize 
}

func findCaller(opts ...option) bool {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logCaller); nil != value {
			return value.(bool)
		}
	}

	return false
}
