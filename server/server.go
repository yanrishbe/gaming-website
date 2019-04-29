package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yanrishbe/gaming-website/entity"

	"github.com/yanrishbe/gaming-website/game"

	"github.com/gorilla/mux"
)

type API struct {
	r *mux.Router
	c game.Controller
}

func (a API) initRouter() {
	a.r.HandleFunc("/user", a.register).Methods(http.MethodPost)
	a.r.HandleFunc("/user/{id}", a.get).Methods(http.MethodGet)
}

func New(c game.Controller) (API, error) {
	a := API{
		r: mux.NewRouter(),
		c: c,
	}
	a.initRouter()
	return a, nil
}

func (a API) GetRouter() *mux.Router {
	return a.r
}

func readID(r *http.Request) (int, error) {
	strID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strID)
	if err != nil {
		return 0, entity.InvIDErr(err)
	}
	return id, nil
}

func (a API) register(w http.ResponseWriter, r *http.Request) {
	u := entity.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		errResp(w, entity.DecodeErr(err))
		return
	}
	u, err = a.c.Register(u)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, u)
}

func (a API) get(w http.ResponseWriter, r *http.Request) {
	id, err := readID(r)
	if err != nil {
		errResp(w, err)
		return
	}
	u, err := a.c.Get(id)
	if err != nil {
		errResp(w, err)
		return
	}
	jsonResp(w, u)
}

func (a API) delete(w http.ResponseWriter, r *http.Request) {
	id, err := readID(r)
	if err != nil {
		errResp(w, err)
		return
	}
	err = a.c.Delete(id)
	if err != nil {
		errResp(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
