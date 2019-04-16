package server

import (
	"encoding/json"
	"net/http"

	"github.com/yanrishbe/gaming-website/entity"

	"github.com/sirupsen/logrus"
)

// JSONResponse encodes user's data for a client
func (a *API) JSONResponse(w http.ResponseWriter, user entity.User) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user": user,
		}).Debug("encoding error")
		return
	}
}
