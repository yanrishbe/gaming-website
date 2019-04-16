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

func errResp(w http.ResponseWriter, err error) {
	resp := entity.HandlerErr(err)
	w.WriteHeader(resp.Code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}

// API struct is used to initialize a router and a database
type API struct {
	Router *mux.Router
	DB     *db.DB
}

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
	a.JSONResponse(w, u)
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
	a.JSONResponse(w, u)
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
	a.JSONResponse(w, u)
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
	a.JSONResponse(w, u)
}

// InitRouter registers handlers
func (a *API) InitRouter() {
	a.Router.HandleFunc("/user", a.registerNewUser).Methods(http.MethodPost)
	a.Router.HandleFunc("/user/{id}", a.getUser).Methods(http.MethodGet)
	a.Router.HandleFunc("/user/{id}", a.deleteUser).Methods(http.MethodDelete)
	a.Router.HandleFunc("/user/{id}/take", a.takeUserPoints).Methods(http.MethodPost)
	a.Router.HandleFunc("/user/{id}/fund", a.fundUserPoints).Methods(http.MethodPost)
}

// New initializes an instance of API struct
func New() *API {
	return &API{
		Router: mux.NewRouter(), DB: db.New(),
	}
}
