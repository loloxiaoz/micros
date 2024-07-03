package main

import (
	"micros/internal/common"
	"micros/internal/config"
	"micros/internal/server"
	"os"
)

func main() {
	common.PrintVersion()
	//flags
	initDir := os.Getenv("GOPATH") + "/src/micros/configs/conf.ini"
	yamlDir := os.Getenv("GOPATH") + "/src/micros/configs/conf.yaml"

	//config
	conf := config.New(initDir, yamlDir)

	//server
	s := server.New(conf)
	s.Run()
}
