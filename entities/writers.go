package entities

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, code int, user User, message string) {
	JSONResponseNoUser(w, code, message)
	if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
		log.Println("error encoding data for a client")
		return
	}
}

//fixme
func JSONResponseNoUser(w http.ResponseWriter, code int, message string) {
	//no body only status code???
	log.Println(message)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
}
