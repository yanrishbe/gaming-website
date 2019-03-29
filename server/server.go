package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/yanrishbe/gaming-website/entities"
	"log"
	"net/http"
	"strconv"
)

func JSONResponse(w http.ResponseWriter, code int, user entities.User, message string) {
	JSONResponseNoUser(w, code, message)
	if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
		log.Println("error encoding data for a client")
		return
	}
}

//fixme
func JSONResponseNoUser(w http.ResponseWriter, code int, message string) {
	//no body only status code???
	log.Println(message)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
}

func registerNewUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User
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

	if entities.IsValid(&user) == false {
		JSONResponse(w, http.StatusBadRequest, user, "user's data is not valid")
		return
	}

	if errSave := entities.SaveUser(&user, &entities.UsersCounter); errSave != nil {
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

	user, doesExist := entities.Users[id]

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

	user, doesExist := entities.Users[id]

	//fixme
	if doesExist == false {
		//send a mock user with mistake or just answer?????
		JSONResponseNoUser(w, http.StatusBadRequest, "the id cannot match any user")
		return
	}

	//fixme
	if errDelete := entities.DeleteUser(user.Id); errDelete != nil {
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

	user, doesExist := entities.Users[id]

	//fixme
	if doesExist == false {
		//send a mock user with mistake or just answer?????
		JSONResponseNoUser(w, http.StatusBadRequest, "the id cannot match any user")
		return
	}

	var points entities.RequestPoints

	if user.Error = json.NewDecoder(r.Body).Decode(&points); user.Error != nil {
		JSONResponse(w, http.StatusUnprocessableEntity, *user, "error decoding client's data")
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			panic(errClose)
		}
	}()

	if errTake := entities.UserTake(user.Id, points.Points); errTake != nil {
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

	user, doesExist := entities.Users[id]

	//fixme
	if doesExist == false {
		//send a mock user with mistake or just answer?????
		JSONResponseNoUser(w, http.StatusBadRequest, "the id cannot match any user")
		return
	}

	var points entities.RequestPoints

	if user.Error = json.NewDecoder(r.Body).Decode(&points); user.Error != nil {
		JSONResponse(w, http.StatusUnprocessableEntity, *user, "error decoding client's data")
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			panic(errClose)
		}
	}()

	if errFund := entities.UserFund(user.Id, points.Points); errFund != nil {
		JSONResponse(w, http.StatusBadRequest, *user, errFund.Error())
	}

	JSONResponse(w, http.StatusOK, *user, "the client successfully funded the points")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/user", registerNewUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", getUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", deleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/user/{id}/take", takeUserPoints).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}/fund", fundUserPoints).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))
}
