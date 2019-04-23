package server

import (
	"encoding/json"
	"net/http"

	"github.com/yanrishbe/gaming-website/entity"

	"github.com/sirupsen/logrus"
)

func jsonResp(w http.ResponseWriter, user entity.User) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user": user,
		}).Debug("encoding error")
		return
	}
}

// keep things consistent.
// If you have both function doing the same job I assume the should
// both be in the same place, and be named in a same way
// so i'm moving it here
func errResp(w http.ResponseWriter, err error) {
	resp := entity.HandlerErr(err)
	w.WriteHeader(resp.Code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// be consistent in what err != nil syntax you're using
	// here you write in 1 line, in previous function that does almost the same you did it in muliple lines.
	// so how it should be written? :)
	// in previous example you log err, here you didn't. So should you log error or not? ))
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}
