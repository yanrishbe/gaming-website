package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func regNew(w http.ResponseWriter, r *http.Request) {
	//todo json writer
}

func New() {
	a := mux.NewRouter()
	a.HandleFunc("/user", regNew)
}
