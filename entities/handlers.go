package entities

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func registerNewUser(w http.ResponseWriter, r *http.Request) {
	var user User

	if errDecode := json.NewDecoder(r.Body).Decode(&user); errDecode != nil {
		user.Error = errDecode.Error()
		JSONResponse(w, http.StatusUnprocessableEntity, user, user.Error)
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()

	if !IsValid(&user) {
		user.Error = errors.New("user's data is not valid").Error()
		JSONResponse(w, http.StatusBadRequest, user, user.Error)
		return
	}

	if errSave := SaveUser(&user, &UsersCounter); errSave != nil {
		user.Error = errSave.Error()
		JSONResponse(w, http.StatusInternalServerError, user, user.Error)
		return
	}

	JSONResponse(w, http.StatusCreated, user, "successfully created a client")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	var userError = new(User)
	//fmt.Println(user)
	if errParams != nil {
		userError.Error = errParams.Error()
		JSONResponse(w, http.StatusBadRequest, *userError, userError.Error)
		return
	}

	user, doesExist := Users[id]
	//если не создаю юзера то далее будет паника т.к. в юзер записывается nil  а потом в nil я пытаюсь записать ошибку
	//fmt.Println(user)
	if !doesExist {
		userError.Error = errors.New("the id cannot match any user").Error()
		JSONResponse(w, http.StatusBadRequest, *userError, userError.Error)
		return
	}

	JSONResponse(w, http.StatusOK, *user, "successfully sent info about the user")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	var userError = new(User)

	if errParams != nil {
		userError.Error = errParams.Error()
		JSONResponse(w, http.StatusBadRequest, *userError, userError.Error)
		return
	}

	user, doesExist := Users[id]

	if !doesExist {
		userError.Error = errors.New("the id cannot match any user").Error()
		JSONResponse(w, http.StatusBadRequest, *userError, userError.Error)
		return
	}

	if errDelete := DeleteUser(user.ID); errDelete != nil {
		user.Error = errDelete.Error()
		JSONResponse(w, http.StatusInternalServerError, *user, user.Error)
		return
	}

	JSONResponseNoUser(w, http.StatusNoContent, "successfully deleted the user")
}

func takeUserPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	var userError = new(User)

	if errParams != nil {
		userError.Error = errParams.Error()
		JSONResponse(w, http.StatusBadRequest, *userError, userError.Error)
		return
	}

	user, doesExist := Users[id]

	if !doesExist {
		userError.Error = errors.New("the id cannot match any user").Error()
		JSONResponse(w, http.StatusBadRequest, *userError, userError.Error)
		return
	}

	var points RequestPoints

	if errDecode := json.NewDecoder(r.Body).Decode(&points); errDecode != nil {
		user.Error = errDecode.Error()
		JSONResponse(w, http.StatusUnprocessableEntity, *user, user.Error)
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()

	if errTake := UserTake(user.ID, points.Points); errTake != nil {
		user.Error = errTake.Error()
		JSONResponse(w, http.StatusBadRequest, *user, user.Error)
		return
	}

	JSONResponse(w, http.StatusOK, *user, "successfully took the points from the client")
}

func fundUserPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	var userError = new(User)

	if errParams != nil {
		userError.Error = errParams.Error()
		JSONResponse(w, http.StatusBadRequest, *userError, userError.Error)
		return
	}

	user, doesExist := Users[id]

	if !doesExist {
		userError.Error = errors.New("the id cannot match any user").Error()
		JSONResponse(w, http.StatusBadRequest, *userError, userError.Error)
		return
	}

	var points RequestPoints

	if errDecode := json.NewDecoder(r.Body).Decode(&points); errDecode != nil {
		user.Error = errDecode.Error()
		JSONResponse(w, http.StatusUnprocessableEntity, *user, "error decoding client's data")
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()

	if errFund := UserFund(user.ID, points.Points); errFund != nil {
		user.Error = errFund.Error()
		JSONResponse(w, http.StatusBadRequest, *user, user.Error)
		return
	}

	JSONResponse(w, http.StatusOK, *user, "the client successfully funded the points")
}

//InitRouter registers handlers and returns a pointer to the router
func InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/user", registerNewUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", getUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", deleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/user/{id}/take", takeUserPoints).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}/fund", fundUserPoints).Methods(http.MethodPost)
	return router
}
