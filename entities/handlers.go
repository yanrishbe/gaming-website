package entities

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func registerNewUser(w http.ResponseWriter, r *http.Request) {
	var user User
	//user error for decoding error
	if user.Error = json.NewDecoder(r.Body).Decode(&user); user.Error != nil {
		JSONResponse(w, http.StatusUnprocessableEntity, user, "error decoding client's data")
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			panic(errClose)
		}
	}()

	if IsValid(&user) == false {
		JSONResponse(w, http.StatusBadRequest, user, "user's data is not valid")
		return
	}

	if errSave := SaveUser(&user, &UsersCounter); errSave != nil {
		JSONResponse(w, http.StatusInternalServerError, user, "error saving a user")
		return
	}

	JSONResponse(w, http.StatusCreated, user, "successfully created a client")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])

	if errParams != nil {
		JSONResponseNoUser(w, http.StatusInternalServerError, "error converting string to int")
		return
	}

	user, doesExist := Users[id]

	//fixme
	if doesExist == false {
		//send a mock user with mistake or just answer?????
		JSONResponseNoUser(w, http.StatusBadRequest, "the id cannot match any user")
		return
	}

	JSONResponse(w, http.StatusOK, *user, "successfully sent info about the user")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])

	if errParams != nil {
		JSONResponseNoUser(w, http.StatusInternalServerError, "error converting string to int")
		return
	}

	user, doesExist := Users[id]

	//fixme
	if doesExist == false {
		//send a mock user with mistake or just answer?????
		JSONResponseNoUser(w, http.StatusBadRequest, "the id cannot match any user")
		return
	}

	//fixme
	if errDelete := DeleteUser(user.Id); errDelete != nil {
		//cannot send a user with error yet????
		JSONResponseNoUser(w, http.StatusInternalServerError, "error deleting a user")
		return
	}

	JSONResponseNoUser(w, http.StatusNoContent, "successfully deleted the user")
}

func takeUserPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])

	if errParams != nil {
		JSONResponseNoUser(w, http.StatusInternalServerError, "error converting string to int")
		return
	}

	user, doesExist := Users[id]

	//fixme
	if doesExist == false {
		//send a mock user with mistake or just answer?????
		JSONResponseNoUser(w, http.StatusBadRequest, "the id cannot match any user")
		return
	}

	var points RequestPoints

	if user.Error = json.NewDecoder(r.Body).Decode(&points); user.Error != nil {
		JSONResponse(w, http.StatusUnprocessableEntity, *user, "error decoding client's data")
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			panic(errClose)
		}
	}()

	if errTake := UserTake(user.Id, points.Points); errTake != nil {
		JSONResponse(w, http.StatusBadRequest, *user, errTake.Error())
	}

	JSONResponse(w, http.StatusOK, *user, "successfully took the points from the client")
}

func fundUserPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])

	if errParams != nil {
		JSONResponseNoUser(w, http.StatusInternalServerError, "error converting string to int")
		return
	}

	user, doesExist := Users[id]

	//fixme
	if doesExist == false {
		//send a mock user with mistake or just answer?????
		JSONResponseNoUser(w, http.StatusBadRequest, "the id cannot match any user")
		return
	}

	var points RequestPoints

	if user.Error = json.NewDecoder(r.Body).Decode(&points); user.Error != nil {
		JSONResponse(w, http.StatusUnprocessableEntity, *user, "error decoding client's data")
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			panic(errClose)
		}
	}()

	if errFund := UserFund(user.Id, points.Points); errFund != nil {
		JSONResponse(w, http.StatusBadRequest, *user, errFund.Error())
	}

	JSONResponse(w, http.StatusOK, *user, "the client successfully funded the points")
}

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/user", registerNewUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", getUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", deleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/user/{id}/take", takeUserPoints).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}/fund", fundUserPoints).Methods(http.MethodPost)
	return router
}
