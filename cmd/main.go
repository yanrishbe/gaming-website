// Package main is used for running the server
package main

import (
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/yanrishbe/gaming-website/server"
)

func main() {
	api := server.New()
	api.Logrus.SetFormatter(&logrus.JSONFormatter{})
	api.InitRouter()
	log.Fatal(http.ListenAndServe(":8080", api.Router))
}
