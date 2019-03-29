package main

import (
	"github.com/yanrishbe/gaming-website/entities"
	"log"
	"net/http"
)

func main() {
	router := entities.InitRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
