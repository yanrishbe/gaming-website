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
	api, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	// do not call log.Fatal, panic() or os.Exit()  outside test or main!!!!
	// this is bad!!!
	// all your functions except tests and functions in *main.go* MUST NOT panic, Fatal, or os.Exit()!!!
	// otherwise it's really simple to mess things up by calling your function and getting unexpected Fatal/panic

	//you should rewrite your code like this:
	//err := server.New()
	//if err != nil {
	//	log.Fatalf("server init: %v", err)
	//}

	//api.InitRouter()
	// you can do InitRouter() inside New(), so you don't have to call it in main
	log.Fatal(http.ListenAndServe(":8080", api.Router))
}
