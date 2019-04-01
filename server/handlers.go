//Package server
package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/yanrishbe/gaming-website/db"
	"github.com/yanrishbe/gaming-website/entities"

	"github.com/gorilla/mux"
)

//RequestPoints represents a struct to send take and fund requests to the gaming website
type RequestPoints struct {
	Points int `json:"points"`
}

//UserResponse struct is a struct used for sending an answer to a client
type UserResponse struct {
	entities.User `json:"user"`
	Error         string `json:"error"`
}

//
type API struct {
	Router *mux.Router
	DB     *db.DB
}

func isValid(user entities.User) bool {
	if user.Name == "" || user.Balance < 300 {
		return false
	}
	return true

}

func (a *API) registerNewUser(w http.ResponseWriter, r *http.Request) {
	var user UserResponse

	if errDecode := json.NewDecoder(r.Body).Decode(&user.User); errDecode != nil {
		user.Error = errDecode.Error()
		JSONResponse(w, http.StatusUnprocessableEntity, user, user.Error)
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()

	if !isValid(user.User) {
		user.Error = errors.New("user's data is not valid").Error()
		JSONResponse(w, http.StatusUnprocessableEntity, user, user.Error)
		return
	}

	if errSave := a.DB.SaveUser(&user.User); errSave != nil {
		user.Error = errSave.Error()
		JSONResponse(w, http.StatusInternalServerError, user, user.Error)
		return
	}

	JSONResponse(w, http.StatusCreated, user, "successfully created a client")
}

func (a *API) getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	var userResponse = new(UserResponse) //equivalent to &UserResponse{}

	if errParams != nil {
		userResponse.Error = errParams.Error()
		JSONResponse(w, http.StatusBadRequest, *userResponse, userResponse.Error) //deceptive request routing
		return
	}

	user, doesExist := a.DB.UsersMap[id]

	if !doesExist {
		userResponse.Error = errors.New("the id cannot match any user").Error()
		JSONResponse(w, http.StatusBadRequest, *userResponse, userResponse.Error) //deceptive request routing?????
		return
	}
	userResponse.User = *user

	JSONResponse(w, http.StatusOK, *userResponse, "successfully sent info about the user")
}

func (a *API) deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	var userResponse = new(UserResponse)

	if errParams != nil {
		userResponse.Error = errParams.Error()
		JSONResponse(w, http.StatusBadRequest, *userResponse, userResponse.Error) //deceptive request routing
		return
	}

	user, doesExist := a.DB.UsersMap[id]

	if !doesExist {
		userResponse.Error = errors.New("the id cannot match any user").Error()
		JSONResponse(w, http.StatusBadRequest, *userResponse, userResponse.Error) //deceptive request routing?????
		return
	}

	if errDelete := a.DB.DeleteUser(user.ID); errDelete != nil {
		userResponse.Error = errDelete.Error()
		JSONResponse(w, http.StatusInternalServerError, *userResponse, userResponse.Error)
		return
	}

	JSONResponseNoUser(w, http.StatusNoContent, "successfully deleted the user")
}

func (a *API) takeUserPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	var userResponse = new(UserResponse)

	if errParams != nil {
		userResponse.Error = errParams.Error()
		JSONResponse(w, http.StatusBadRequest, *userResponse, userResponse.Error) //deceptive request routing
		return
	}

	user, doesExist := a.DB.UsersMap[id]

	if !doesExist {
		userResponse.Error = errors.New("the id cannot match any user").Error()
		JSONResponse(w, http.StatusBadRequest, *userResponse, userResponse.Error) //deceptive request routing?????
		return
	}

	var points RequestPoints

	if errDecode := json.NewDecoder(r.Body).Decode(&points); errDecode != nil {
		userResponse.Error = errDecode.Error()
		JSONResponse(w, http.StatusUnprocessableEntity, *userResponse, userResponse.Error)
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()

	if errTake := a.DB.UserTake(user.ID, points.Points); errTake != nil {
		userResponse.Error = errTake.Error()
		JSONResponse(w, http.StatusUnprocessableEntity, *userResponse, userResponse.Error)
		return
	}

	userResponse.User = *a.DB.UsersMap[id]
	JSONResponse(w, http.StatusOK, *userResponse, "successfully took the points from the client")
}

func (a *API) fundUserPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	var userResponse = new(UserResponse)

	if errParams != nil {
		userResponse.Error = errParams.Error()
		JSONResponse(w, http.StatusBadRequest, *userResponse, userResponse.Error) //deceptive request routing
		return
	}

	user, doesExist := a.DB.UsersMap[id]

	if !doesExist {
		userResponse.Error = errors.New("the id cannot match any user").Error()
		JSONResponse(w, http.StatusBadRequest, *userResponse, userResponse.Error) //deceptive request routing?????
		return
	}

	var points RequestPoints

	if errDecode := json.NewDecoder(r.Body).Decode(&points); errDecode != nil {
		userResponse.Error = errDecode.Error()
		JSONResponse(w, http.StatusUnprocessableEntity, *userResponse, userResponse.Error)
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()

	if errFund := a.DB.UserFund(user.ID, points.Points); errFund != nil {
		userResponse.Error = errFund.Error()
		JSONResponse(w, http.StatusUnprocessableEntity, *userResponse, userResponse.Error)
		return
	}

	userResponse.User = *a.DB.UsersMap[id]
	JSONResponse(w, http.StatusOK, *userResponse, "the client successfully funded the points")
}

//InitRouter registers handlers and returns a pointer to the router
func (a *API) InitRouter() {
	//a.Router = mux.NewRouter()
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

func (a *API) New() {
	a.Router = mux.NewRouter()
	a.DB = new(db.DB)
	a.DB.UsersMap = make(map[int]*entities.User)
	a.DB.UsersCounter = 0
	//a.DB = &db.DB{UsersMap: *new(map[int]*entities.User), UsersCounter: 0}
}
