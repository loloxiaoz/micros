package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"micros/internal/common"
	"micros/internal/config"
	"micros/internal/db"
	"micros/internal/log"
	"micros/internal/server"
	"micros/internal/service"
)

func main() {
	common.PrintVersion()

	//flags
	initDir := flag.String("ci", os.Getenv("GOPATH")+"/src/micros/configs/conf.ini", "config file path, ini")
	yamlDir := flag.String("cy", os.Getenv("GOPATH")+"/src/micros/configs/conf.yaml", "config file path, yaml")
	flag.Parse()

	//config
	conf, err := config.New(*initDir, *yamlDir)
	if err != nil {
		fmt.Printf("config init fail, error is %s", err.Error())
		os.Exit(1)
	}

	//logger
	if err := log.Init(&conf.Log); err != nil {
		fmt.Printf("logger init fail, error is %s", err.Error())
		os.Exit(1)
	}

	// DB
	initDB(&conf.DB)

	//server
	startServer(conf)
}

func initDB(c *config.DB) {

	db.Init(c)
	service.Init()
}

func startServer(conf *config.Conf) {
	s := server.New(conf)
	go func() {
		if err := s.Run(); err != nil && err != http.ErrServerClosed {
			log.Logger().Errorf("server listen err:%s", err)
		}
	}()

	//shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Logger().Fatal("server shutdown error")
	}
	log.Logger().Info("server exit")
}