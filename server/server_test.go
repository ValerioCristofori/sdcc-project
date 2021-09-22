package main

import (
	"syscall"
	"testing"
)


func TestServer(t *testing.T)  {

	go serveData()

	syscall.Pause()
}


