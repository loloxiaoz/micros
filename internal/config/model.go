package config

//DB 数据库配置
type DB struct {
	DBType      string `json:"dbType"`
	Database  	string `json:"database"`
	Host        string `json:"host"`
	Port 		string `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	MaxIdleConn int    `json:"maxIdleConn"`
	MaxOpenConn int    `json:"maxOpenConn"`
	MaxIdleTime int    `json:"maxIdleTime"`
}

//Log 日志配置
type Log struct {
	Path  string `json:"path"`
	Level string `json:"level"`
	RotateSize string `json:"rotateSize"`
	RotateHourly string `json:"rotateHourly"`
	Rotate string `json:"rotate"`
	Retention int `json:"retention"`
}

//Opt 选项配置
type Opt struct {
	APIDoc string `json:"apiDoc"`
	Profile string `json:"profile"`
	Monitor string 	`json:"monitor"`
	Trace string `json:"trace"`
	Stat string `json:"stat"`
}
