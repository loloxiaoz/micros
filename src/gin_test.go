package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestDemo(t *testing.T) {
	ret, _ := http.Get("http://127.0.0.1:8080/")
	fmt.Println(ret)
}
