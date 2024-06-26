package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"github.com/gookit/config/v2/yaml"
)


var conf Conf

//Conf 全部配置
type Conf struct {
	DB `json:"DB"`
	Log `json:"Log"`
	Opt `json:"Opt"`
	Project string `json:"project"`
}

//Init 初始化
func (conf *Conf) Init(path ...string) {
	// 支持ENV变量解析
	config.WithOptions(config.ParseEnv)
	config.WithOptions(func(opt *config.Options) {
		opt.DecoderConfig.TagName = "json"
	})

	// 添加yaml驱动
	config.AddDriver(yaml.Driver)
	config.AddDriver(ini.Driver)

	// 加载配置，可以同时传入多个文件
	err := config.LoadFiles(path...)
	if err != nil {
		panic(err)
	}
	config.Decode(&conf)
}

//IsAPIDoc 是否开启api文档
func (conf *Conf) IsAPIDoc() bool {
	return conf.Opt.APIDoc == "true"
}

//IsProfile 是否开启profile
func (conf *Conf) IsProfile() bool {
	return conf.Opt.Profile == "true"
}

//IsMonitor 是否开启监控
func (conf *Conf) IsMonitor() bool {
	return conf.Opt.Monitor == "true"
}

//IsTrace 是否开启trace
func (conf *Conf) IsTrace() bool {
	return conf.Opt.Trace == "true"
}

//IsStat 是否开启统计
func (conf *Conf) IsStat() bool {
	return conf.Opt.Stat == "true"
}