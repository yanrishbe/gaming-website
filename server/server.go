package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yanrishbe/gaming-website/entity"

	"github.com/yanrishbe/gaming-website/game"

	"github.com/gorilla/mux"
)

type ReqPoints struct {
	Points int `json:"points"`
}

type API struct {
	r *mux.Router
	c game.Controller
}

func New(c game.Controller) (*mux.Router, error) {
	a := API{
		r: mux.NewRouter(),
		c: c,
	}
	a.r.HandleFunc("/user", a.regUser).Methods(http.MethodPost)
	a.r.HandleFunc("/user/{id}", a.getUser).Methods(http.MethodGet)
	a.r.HandleFunc("/user/{id}", a.delUser).Methods(http.MethodDelete)
	a.r.HandleFunc("/user/{id}/take", a.takePoints).Methods(http.MethodPost)
	a.r.HandleFunc("/user/{id}/fund", a.fundPoints).Methods(http.MethodPost)
	return a.r, nil
}

func readID(r *http.Request) (int, error) {
	strID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strID)
	if err != nil {
		return 0, entity.InvIDErr(err)
	}
	return id, nil
}

func (a API) regUser(w http.ResponseWriter, r *http.Request) {
	u := entity.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		errResp(w, entity.DecodeErr(err))
		return
	}
	u, err = a.c.RegUser(u)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, u)
}

func (a API) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := readID(r)
	if err != nil {
		errResp(w, err)
		return
	}
	u, err := a.c.GetUser(id)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, u)
}

func (a API) delUser(w http.ResponseWriter, r *http.Request) {
	id, err := readID(r)
	if err != nil {
		errResp(w, err)
		return
	}
	err = a.c.DelUser(id)
	if err != nil {
		errResp(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) takePoints(w http.ResponseWriter, r *http.Request) {
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
	u, err := a.c.TakePoints(id, points.Points)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, u)
}

func (a *API) fundPoints(w http.ResponseWriter, r *http.Request) {
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
	u, err := a.c.FundPoints(id, points.Points)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, u)
}
