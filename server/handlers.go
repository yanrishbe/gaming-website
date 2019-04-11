// Package server is used to handle all client's requests
package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/yanrishbe/gaming-website/db"
	"github.com/yanrishbe/gaming-website/entity"

	"github.com/gorilla/mux"
)

// ReqPoints represents a struct to send "take" and "fund" requests to the gaming website
type ReqPoints struct {
	Points int `json:"points"`
}

// UserResp struct is a struct used for sending an answer to a client
type UserResp struct {
	entity.User `json:"user"`
	Error       string `json:"error"`
}

// API struct is used to initialize a router and a database
type API struct {
	Router *mux.Router
	DB     *db.DB
	Logrus *logrus.Logger
}

func (a *API) registerNewUser(w http.ResponseWriter, r *http.Request) {
	us := UserResp{}
	err := json.NewDecoder(r.Body).Decode(&us.User)
	if err != nil {
		us.Error = err.Error()
		a.JSONResponse(w, http.StatusUnprocessableEntity, us, us.Error)
		return
	}
	us.User, err = a.DB.SaveUser(us.User)
	if err != nil {
		us.Error = err.Error()
		if us.Error == "user's data is not valid" {
			a.JSONResponse(w, http.StatusUnprocessableEntity, us, us.Error)
			return
		}
		a.JSONResponse(w, http.StatusInternalServerError, us, us.Error)
		return
	}
	a.JSONResponse(w, http.StatusCreated, us, "successfully created a client")
}

func (a *API) getUser(w http.ResponseWriter, r *http.Request) {
	ur := UserResp{}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ur.Error = err.Error()
		a.JSONResponse(w, http.StatusBadRequest, ur, ur.Error)
		return
	}
	us, err := a.DB.GetUser(id)

	if err != nil {
		ur.Error = err.Error()
		a.JSONResponse(w, http.StatusNotFound, ur, ur.Error)
		return
	}
	ur.User = us
	a.JSONResponse(w, http.StatusOK, ur, "successfully sent info about the user")
}

func (a *API) deleteUser(w http.ResponseWriter, r *http.Request) {
	ur := UserResp{}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ur.Error = err.Error()
		a.JSONResponse(w, http.StatusBadRequest, ur, ur.Error)
		return
	}

	err = a.DB.DeleteUser(id)
	if err != nil {
		ur.Error = err.Error()
		if ur.Error == "the id cannot match any user" {
			a.JSONResponse(w, http.StatusNotFound, ur, ur.Error)
			return
		}
		a.JSONResponse(w, http.StatusInternalServerError, ur, ur.Error)
		return
	}

	a.ResponseNoUser(w, http.StatusOK, "successfully deleted the user")
}

func (a *API) takeUserPoints(w http.ResponseWriter, r *http.Request) {
	ur := UserResp{}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ur.Error = err.Error()
		a.JSONResponse(w, http.StatusBadRequest, ur, ur.Error)
		return
	}

	points := ReqPoints{}

	err = json.NewDecoder(r.Body).Decode(&points)
	if err != nil {
		ur.Error = err.Error()
		a.JSONResponse(w, http.StatusUnprocessableEntity, ur, ur.Error)
		return
	}

	ur.User, err = a.DB.UserTake(id, points.Points)
	if err != nil {
		ur.Error = err.Error()
		if ur.Error == "the id cannot match any user" {
			a.JSONResponse(w, http.StatusNotFound, ur, ur.Error)
			return
		}
		a.JSONResponse(w, http.StatusUnprocessableEntity, ur, ur.Error)
		return
	}
	a.JSONResponse(w, http.StatusOK, ur, "successfully took the points from the client")
}

func (a *API) fundUserPoints(w http.ResponseWriter, r *http.Request) {
	ur := UserResp{}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ur.Error = err.Error()
		a.JSONResponse(w, http.StatusBadRequest, ur, ur.Error)
		return
	}
	points := ReqPoints{}

	err = json.NewDecoder(r.Body).Decode(&points)
	if err != nil {
		ur.Error = err.Error()
		a.JSONResponse(w, http.StatusUnprocessableEntity, ur, ur.Error)
		return
	}
	ur.User, err = a.DB.UserFund(id, points.Points)
	if err != nil {
		ur.Error = err.Error()
		if ur.Error == "the id cannot match any user" {
			a.JSONResponse(w, http.StatusNotFound, ur, ur.Error)
			return
		}
		a.JSONResponse(w, http.StatusUnprocessableEntity, ur, ur.Error)
		return
	}
	a.JSONResponse(w, http.StatusOK, ur, "the client successfully funded the points")
}

// InitRouter registers handlers
func (a *API) InitRouter() {
	a.Logrus.SetFormatter(&logrus.JSONFormatter{})
	a.Router.HandleFunc("/user", a.registerNewUser).Methods(http.MethodPost)
	a.Router.HandleFunc("/user/{id}", a.getUser).Methods(http.MethodGet)
	a.Router.HandleFunc("/user/{id}", a.deleteUser).Methods(http.MethodDelete)
	a.Router.HandleFunc("/user/{id}/take", a.takeUserPoints).Methods(http.MethodPost)
	a.Router.HandleFunc("/user/{id}/fund", a.fundUserPoints).Methods(http.MethodPost)
}

// New initializes an instance of API struct
func New() *API {
	return &API{
		Router: mux.NewRouter(), DB: db.New(), Logrus: logrus.New(),
	}
}
