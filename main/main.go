// Package main is used for running the server
package main

import (
	"github.com/yanrishbe/gaming-website/server"
)

func main() {
	api := server.New()
	api.InitRouter()
	api.Run(":8080")
}
