package main

import (
	"micros/internal/common"
	"micros/internal/server"
)

func main() {
	common.PrintVersion()
	s := server.NewServer("hero")
	s.Run()
}
