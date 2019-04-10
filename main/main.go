// Package main is used for running the server
package main

import (
	"github.com/yanrishbe/gaming-website/logger"
	"github.com/yanrishbe/gaming-website/server"
)

func main() {
	log := logger.New("debug")
	api := server.New(log)
	api.InitRouter()
	api.Run(":8080")
}
