package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/yanrishbe/gaming-website/game"
	"github.com/yanrishbe/gaming-website/postgres"
	"github.com/yanrishbe/gaming-website/server"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	db, err := postgres.New()
	if err != nil {
		logrus.Fatal(err)
	}
	api, err := server.New(game.New(db))
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Fatal(http.ListenAndServe(":8080", api.GetRouter()))
}
