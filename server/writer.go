package server

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/yanrishbe/gaming-website/entity"
)

func jsonResp(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		var logField string
		switch i.(type) {
		case entity.User:
			logField = "user"
		case entity.Tournament:
			logField = "tournament"
		}
		logrus.WithFields(logrus.Fields{
			logField: i,
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
