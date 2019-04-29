package server

import (
	"encoding/json"
	"net/http"

	"github.com/yanrishbe/gaming-website/entity"

	"github.com/yanrishbe/gaming-website/game"

	"github.com/gorilla/mux"
)

type API struct {
	r *mux.Router
	c game.Controller
}

func (a API) initRouter() {
	a.r.HandleFunc("/user", a.regUser).Methods(http.MethodPost)
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

func (a API) regUser(w http.ResponseWriter, r *http.Request) {
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
