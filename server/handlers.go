// Package server is used to handle all client's requests
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/yanrishbe/gaming-website/db"
	"github.com/yanrishbe/gaming-website/entities"

	"github.com/gorilla/mux"
)

// RequestPoints represents a struct to send "take" and "fund" requests to the gaming website
type RequestPoints struct {
	Points int `json:"points"`
}

// UserResponse struct is a struct used for sending an answer to a client
type UserResponse struct {
	ID            int `json:"id"`
	entities.User `json:"user"`
	Error         string `json:"error"`
}

// API struct is used to initialize a router and a database
type API struct {
	Router *mux.Router
	DB     *db.DB
	Logrus *logrus.Logger
}

func (a *API) registerNewUser(w http.ResponseWriter, r *http.Request) {
	var user UserResponse
	if errDecode := json.NewDecoder(r.Body).Decode(&user.User); errDecode != nil {
		user.Error = errDecode.Error()
		a.JSONResponse(w, http.StatusUnprocessableEntity, user, user.Error)
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()

	id, errSave := a.DB.SaveUser(&user.User)
	if errSave != nil {
		user.Error = errSave.Error()
		if match := strings.EqualFold(user.Error, "user's data is not valid"); match {
			a.JSONResponse(w, http.StatusUnprocessableEntity, user, user.Error)
			return
		}
		a.JSONResponse(w, http.StatusInternalServerError, user, user.Error)
		return
	}
	user.ID = id
	a.JSONResponse(w, http.StatusCreated, user, "successfully created a client")
}

func (a *API) getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	var userResponse = new(UserResponse)

	if errParams != nil {
		userResponse.Error = errParams.Error()
		a.JSONResponse(w, http.StatusBadRequest, *userResponse, userResponse.Error)
		return
	}

	user, errGet := a.DB.GetUser(id)

	if errGet != nil {
		userResponse.Error = errGet.Error()
		a.JSONResponse(w, http.StatusNotFound, *userResponse, userResponse.Error)
		return
	}
	userResponse.User = *user
	userResponse.ID = id
	a.JSONResponse(w, http.StatusOK, *userResponse, "successfully sent info about the user")
}

func (a *API) deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	var userResponse = new(UserResponse)

	if errParams != nil {
		userResponse.Error = errParams.Error()
		a.JSONResponse(w, http.StatusBadRequest, *userResponse, userResponse.Error)
		return
	}

	if errDelete := a.DB.DeleteUser(id); errDelete != nil {
		userResponse.Error = errDelete.Error()
		if match := strings.EqualFold(userResponse.Error, "the id cannot match any user"); match {
			a.JSONResponse(w, http.StatusNotFound, *userResponse, userResponse.Error)
			return
		}
		a.JSONResponse(w, http.StatusInternalServerError, *userResponse, userResponse.Error)
		return
	}

	a.ResponseNoUser(w, http.StatusOK, "successfully deleted the user")
}

func (a *API) takeUserPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	var userResponse = new(UserResponse)

	if errParams != nil {
		userResponse.Error = errParams.Error()
		a.JSONResponse(w, http.StatusBadRequest, *userResponse, userResponse.Error)
		return
	}

	var points RequestPoints

	if errDecode := json.NewDecoder(r.Body).Decode(&points); errDecode != nil {
		userResponse.Error = errDecode.Error()
		a.JSONResponse(w, http.StatusUnprocessableEntity, *userResponse, userResponse.Error)
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()

	if errTake := a.DB.UserTake(id, points.Points); errTake != nil {
		userResponse.Error = errTake.Error()
		if match := strings.EqualFold(userResponse.Error, "the id cannot match any user"); match {
			a.JSONResponse(w, http.StatusNotFound, *userResponse, userResponse.Error)
			return
		}
		a.JSONResponse(w, http.StatusUnprocessableEntity, *userResponse, userResponse.Error)
		return
	}
	user, errGet := a.DB.GetUser(id)
	if errGet != nil { //fixme
		panic(errGet)
	}
	userResponse.User = *user
	userResponse.ID = id
	a.JSONResponse(w, http.StatusOK, *userResponse, "successfully took the points from the client")
}

func (a *API) fundUserPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	var userResponse = new(UserResponse)

	if errParams != nil {
		userResponse.Error = errParams.Error()
		a.JSONResponse(w, http.StatusBadRequest, *userResponse, userResponse.Error)
		return
	}

	var points RequestPoints

	if errDecode := json.NewDecoder(r.Body).Decode(&points); errDecode != nil {
		userResponse.Error = errDecode.Error()
		a.JSONResponse(w, http.StatusUnprocessableEntity, *userResponse, userResponse.Error)
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()

	if errFund := a.DB.UserFund(id, points.Points); errFund != nil {
		userResponse.Error = errFund.Error()
		if match := strings.EqualFold(userResponse.Error, "the id cannot match any user"); match {
			a.JSONResponse(w, http.StatusNotFound, *userResponse, userResponse.Error)
			return
		}
		a.JSONResponse(w, http.StatusUnprocessableEntity, *userResponse, userResponse.Error)
		return
	}
	user, _ := a.DB.GetUser(id)
	userResponse.User = *user
	userResponse.ID = id
	a.JSONResponse(w, http.StatusOK, *userResponse, "the client successfully funded the points")
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

// Run the app on it's router
func (a *API) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

// New initializes an instance of API struct
func New(log *logrus.Logger) *API {
	return &API{
		Router: mux.NewRouter(), DB: db.New(), Logrus: log,
	}
}
