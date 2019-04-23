// Package server is used to handle all client's requests
package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yanrishbe/gaming-website/db"
	"github.com/yanrishbe/gaming-website/entity"

	"github.com/gorilla/mux"
)

// ReqPoints represents a struct to send "take" and "fund" requests to the gaming website
type ReqPoints struct {
	Points int `json:"points"`
}

// API struct is used to initialize a router and a database
type API struct {
	Router *mux.Router
	DB     db.DB //previously *db.DB
}

// nice helper func :)
func readID(r *http.Request) (int, error) {
	strID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strID)
	if err != nil {
		return 0, entity.InvIDErr(err)
	}
	return id, nil
}

func (a *API) registerNewUser(w http.ResponseWriter, r *http.Request) {
	u := entity.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		errResp(w, entity.DecodeErr(err))
		return
	}
	u, err = a.DB.SaveUser(u)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, u)
}

func (a *API) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := readID(r)
	if err != nil {
		errResp(w, err)
		return
	}
	u, err := a.DB.GetUser(id)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, u)
}

func (a *API) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := readID(r)
	if err != nil {
		errResp(w, err)
		return
	}
	err = a.DB.DeleteUser(id)
	if err != nil {
		errResp(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) takeUserPoints(w http.ResponseWriter, r *http.Request) {
	id, err := readID(r)
	if err != nil {
		errResp(w, err)
		return
	}
	points := ReqPoints{}
	err = json.NewDecoder(r.Body).Decode(&points)
	if err != nil {
		errResp(w, entity.DecodeErr(err))
		return
	}
	u, err := a.DB.UserTake(id, points.Points)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, u)
}

func (a *API) fundUserPoints(w http.ResponseWriter, r *http.Request) {
	id, err := readID(r)
	if err != nil {
		errResp(w, err)
		return
	}
	points := ReqPoints{}
	err = json.NewDecoder(r.Body).Decode(&points)
	if err != nil {
		errResp(w, entity.DecodeErr(err))
		return
	}
	u, err := a.DB.UserFund(id, points.Points)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, u)
}

// initRouter registers handlers
func (a *API) InitRouter() {
	a.Router.HandleFunc("/user", a.registerNewUser).Methods(http.MethodPost)
	a.Router.HandleFunc("/user/{id}", a.getUser).Methods(http.MethodGet)
	a.Router.HandleFunc("/user/{id}", a.deleteUser).Methods(http.MethodDelete)
	a.Router.HandleFunc("/user/{id}/take", a.takeUserPoints).Methods(http.MethodPost)
	a.Router.HandleFunc("/user/{id}/fund", a.fundUserPoints).Methods(http.MethodPost)
}

// Handler's code is very good. You just have minor issues with initialization and object dependency management

// New initializes an instance of API struct
func New() (*API, error) {
	db, err := db.New() // You should not create database here, you should create it in main and pass it here in params to New()
	// Then it will be much clearer who is responsible for closing database.
	if err != nil {
		return nil, err
	}
	a := &API{
		Router: mux.NewRouter(), DB: db,
	}
	a.InitRouter()
	return a, nil
}
