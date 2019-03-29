//Package main runs the server
package main

import (
	"log"
	"net/http"

	"github.com/yanrishbe/gaming-website/entities"
)

func main() {
	router := entities.InitRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
