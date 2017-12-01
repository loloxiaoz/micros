package main

import (
	"github.com/micros/server"
)

func main() {
	prjName := "example"
	configPath := "./config.yaml"
	server.NewServer(prjName, configPath)
}
