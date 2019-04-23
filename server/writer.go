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

func errResp(w http.ResponseWriter, err error) {
	resp := entity.HandlerErr(err)
	w.WriteHeader(resp.Code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}
