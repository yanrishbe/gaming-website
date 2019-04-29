package server

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/yanrishbe/gaming-website/entity"
)

func jsonResp(w http.ResponseWriter, u entity.User) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(u)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user": u,
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
		logrus.Debug("encoding error")
	}
	return
}
