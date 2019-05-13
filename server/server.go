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
	a.r.HandleFunc("/tournament", a.regTourn).Methods(http.MethodPost)
	a.r.HandleFunc("/tournament/{id}", a.getTourn).Methods(http.MethodGet)
	a.r.HandleFunc("/tournament/{id}/join", a.joinTourn).Methods(http.MethodPost)
	a.r.HandleFunc("/tournament/{id}/finish", a.finishTourn).Methods(http.MethodPost)
	a.r.HandleFunc("/tournament/{id}", a.delTourn).Methods(http.MethodDelete)
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
func (a API) regTourn(w http.ResponseWriter, r *http.Request) {
	t := entity.Tournament{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		errResp(w, entity.DecodeErr(err))
		return
	}
	t.Users = []entity.Winner{}
	t, err = a.c.RegTourn(t)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, t)
}

func (a API) getTourn(w http.ResponseWriter, r *http.Request) {
	id, err := readID(r)
	if err != nil {
		errResp(w, err)
		return
	}
	t, err := a.c.GetTourn(id)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, t)
}

func (a API) finishTourn(w http.ResponseWriter, r *http.Request) {
	id, err := readID(r)
	if err != nil {
		errResp(w, err)
		return
	}
	t, err := a.c.FinishTourn(id)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, t)
}

func (a API) delTourn(w http.ResponseWriter, r *http.Request) {
	id, err := readID(r)
	if err != nil {
		errResp(w, err)
		return
	}
	err = a.c.DelTourn(id)
	if err != nil {
		errResp(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a API) joinTourn(w http.ResponseWriter, r *http.Request) {
	u := entity.UserTourn{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		errResp(w, entity.DecodeErr(err))
		return
	}
	id, err := readID(r)
	if err != nil {
		errResp(w, err)
		return
	}
	t, err := a.c.JoinTourn(id, u.ID)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, t)
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
