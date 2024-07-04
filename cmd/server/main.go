package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"micros/internal/common"
	"micros/internal/config"
	"micros/internal/logger"
	"micros/internal/server"
)

func main() {
	common.PrintVersion()
	flag.Parse()
	//flags
	initDir := os.Getenv("GOPATH") + "/src/micros/configs/conf.ini"
	yamlDir := os.Getenv("GOPATH") + "/src/micros/configs/conf.yaml"

	//config
	conf, err := config.New(initDir, yamlDir)
	if err !=nil {
		fmt.Printf("config init fail, error is %s", err.Error())
		os.Exit(1)
	}

	//logger
	logger.Init(&conf.Log)

	//server
	ctx, cancel := context.WithCancel(context.Background())
	s := server.New(conf)
	logger.Log.Info("server starting")
	s.Run(ctx)
	cancel()
	logger.Log.Warn("server stoped!")
}
