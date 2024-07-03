package logger

import (
	"fmt"
	"os"
	"strings"
)

// Config 日志配置
type Config struct {
	//存储路径
	Path string
	//日志级别
	Level string
	//日志标签 多日志时使用
	Tag string
	//日志格式
	Format string
	//Console输出
	Console bool
}

// Level 日志级别
type Level int

const (
	//DEBUG 级别
	DEBUG Level = iota
	//INFO 级别
	INFO 
	//WARN 级别
	WARN 
	//ERROR 级别
	ERROR
	//TRACE 级别
	TRACE
	//CRITICAL 级别
	CRITICAL
)

var levels 	= map[string]Level {
	 "DEBUG": DEBUG,
	 "INFO": INFO,
	 "WARN": WARN,
	 "ERROR": ERROR,
	 "TRACE": TRACE,
	 "CRITICAl" : CRITICAL,
}

// InitWithConfig 初始化日志配置
func InitWithConfig(c *Config) {
	c.check()
}

func (c *Config) check() {

	//level init
	if _, ok := levels[c.Level]; !ok {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: log level is wrong,  %s\n", c.Level)
			os.Exit(1)
	}
	paths := strings.Split(c.Path, "/")
	if len(paths) > 1 {
		//create path
		dir := strings.Join(paths[0:len(paths)-1], "/")
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not create directory %s, err:%s\n", dir, err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: log directory invalid %s\n", c.Path)
		os.Exit(1)
	}
}
