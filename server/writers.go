package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSONResponse encodes user's data for a client
func JSONResponse(w http.ResponseWriter, code int, user UserResponse, message string) {
	log.Println(message)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
		log.Println(errAnswer)
		return
	}
}

// JSONResponseNoUser encodes data for a client without  returning a User struct entity
func ResponseNoUser(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if _, errWrite := w.Write([]byte(message)); errWrite != nil {
		log.Println(errWrite)
	}

}
