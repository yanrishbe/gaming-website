// Package main is used for running the server
package main

import (
	"log"
	"net/http"

	"github.com/yanrishbe/gaming-website/db"

	"github.com/sirupsen/logrus"
	"github.com/yanrishbe/gaming-website/server"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	connStr := "user=postgres password=docker2147 dbname=gaming_website host=localhost port=5432 sslmode=disable"
	dbm, err := db.New(connStr)
	if err != nil {
		log.Fatal(err)
	}
	api, err := server.New(dbm)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", api.Router))
}
