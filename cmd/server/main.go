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
	dir := os.Getenv("GOPATH") + "/src/micros/configs/conf.ini"

	//config
	conf := config.New(dir)

	//server
	s := server.NewServer(conf)
	s.Run()
}
