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

	//flags
	initDir := flag.String("ci", os.Getenv("GOPATH") + "/src/micros/configs/conf.ini", "config file path, ini")
	yamlDir := flag.String("cy", os.Getenv("GOPATH") + "/src/micros/configs/conf.yaml", "config file path, yaml")
	flag.Parse()

	fmt.Printf("init config file path is %s\n", *initDir)
	fmt.Printf("yaml config file path is %s\n", *yamlDir)

	//config
	conf, err := config.New(*initDir, *yamlDir)
	if err !=nil {
		fmt.Printf("config init fail, error is %s", err.Error())
		os.Exit(1)
	}

	//logger
	if err := logger.Init(&conf.Log); err != nil {
		fmt.Printf("logger init fail, error is %s", err.Error())
		os.Exit(1)
	}

	//server
	ctx, cancel := context.WithCancel(context.Background())
	s := server.New(conf)
	s.Run(ctx)
	cancel()
}
