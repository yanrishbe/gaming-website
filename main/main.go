//Package main is used to run the server
package main

import (
	"github.com/yanrishbe/gaming-website/server"
)

func main() {
	api := server.API{}
	api.New()
	//fmt.Println(api)
	api.InitRouter()
	api.Run(":8080")
}
