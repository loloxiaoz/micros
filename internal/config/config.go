package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"github.com/gookit/config/v2/yaml"
)

func init() {
	// 支持ENV变量解析
	config.WithOptions(config.ParseEnv)

	// 添加yaml驱动
	config.AddDriver(yaml.Driver)
	config.AddDriver(ini.Driver)

	// 加载配置，可以同时传入多个文件
	err := config.LoadFiles("configs/conf.ini", "configs/conf.yaml")
	if err != nil {
		panic(err)
	}
}