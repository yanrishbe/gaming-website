// Package main is used for running the server
package main

import (
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/yanrishbe/gaming-website/server"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	api := server.New()
	api.InitRouter()
	log.Fatal(http.ListenAndServe(":8080", api.Router))
}
