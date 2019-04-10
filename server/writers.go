package server

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// JSONResponse encodes user's data for a client
func (a *API) JSONResponse(w http.ResponseWriter, code int, user UserResponse, message string) {
	//a.Logrus.Debug(message)
	a.Logrus.WithFields(log.Fields{
		"code":    code,
		"user":    user,
		"message": message,
	}).Debug()
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
		//log.Println(message)
		a.Logrus.WithFields(log.Fields{
			"code":    code,
			"user":    user,
			"message": message,
		}).Debug(code, user, message)
		return
	}
}

// ResponseNoUser encodes data for a client without returning a User struct entity
func (a *API) ResponseNoUser(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if _, errWrite := w.Write([]byte(message)); errWrite != nil {
		//log.Println(message)
		a.Logrus.WithFields(log.Fields{
			"code":    code,
			"user":    "no user",
			"message": message,
		}).Debug()
	}
}
