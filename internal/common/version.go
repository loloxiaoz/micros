package common

import (
	"fmt"
)

var (
	//Version 版本号
	Version string
	//BuildTime 编译时间
	BuildTime string
	//Branch 分支
	Branch string
	//CommitId 代码提交ID
	CommitId string
	//CommitDate 代码提交日期
	CommitDate string
	//GoVersion go版本
	GoVersion string
)

// PrintVersion 打印版本信息
func PrintVersion() {
	fmt.Printf("version: %s\n", Version)
	fmt.Printf("buildTime: %s\n", BuildTime)
	fmt.Printf("branch: %s\n", Branch)
	fmt.Printf("commitID: %s\n", CommitId)
	fmt.Printf("commitDate: %s\n", CommitDate)
	fmt.Printf("go version: %s\n", GoVersion)
}
