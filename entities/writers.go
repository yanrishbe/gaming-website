package entities

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, code int, user User, message string) {
	log.Println(message)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
		log.Println(errAnswer)
		return
	}
}

func JSONResponseNoUser(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if _, errWrite := w.Write([]byte(message)); errWrite != nil {
		log.Println(errWrite)
	}

}
