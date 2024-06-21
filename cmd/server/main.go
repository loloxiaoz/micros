package main

import (
	"micros/internal/common"
	"micros/internal/server"
)

func main() {
	common.PrintVersion()
	server.NewServer("hero")
}
