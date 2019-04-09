package server

import (
	"encoding/json"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

// JSONResponse encodes user's data for a client
func JSONResponse(w http.ResponseWriter, code int, user UserResponse, message string) {
	//log.Println(message)
	log.WithFields(log.Fields{
		"code":    code,
		"user":    user,
		"message": message,
	}).Debug(code, user, message)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
		log.WithFields(log.Fields{
			"code":    code,
			"user":    user,
			"message": message,
		}).Debug(code, user, message)
		return
	}
}

// ResponseNoUser encodes data for a client without returning a User struct entity
func ResponseNoUser(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if _, errWrite := w.Write([]byte(message)); errWrite != nil {
		//log.Debug(errWrite)
		//log.SetFormatter(&log.JSONFormatter{})
		log.WithFields(log.Fields{
			"code":    code,
			"user":    "no user",
			"message": message,
		}).Debug()
	}

}
