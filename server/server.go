package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/yanrishbe/gaming-website/entities"
	"log"
	"net/http"
	"strconv"
)

func JsonWrite(){

}
func UserAvailable(isEmpty bool, user entities.User, w http.ResponseWriter){
	if !isEmpty {
		user.Error = errors.New("the id cannot match any user")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
			log.Println("error encoding data for a client")
			return
		}
		return
	}
}

func registerNewUser(w http.ResponseWriter, r *http.Request) { //TODO: sth strange with sending status codes & errors
	var user entities.User

	if user.Error = json.NewDecoder(r.Body).Decode(&user); user.Error != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
			log.Println("error decoding client's data: ", user.Error)
			return
		}
		return
	}

	entities.IsValid(&user)
	entities.SaveUser(&user)

	if user.Error != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
			panic(errAnswer)
		}
		log.Println("error: ", user.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
		log.Println("error encoding data for a client")
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			panic(errClose)
		}
	}()
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])

	if errParams != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error converting string to int: ", errParams)
		return
	}

	user, isEmpty := entities.Users[id]

	if !isEmpty {
		user.Error = errors.New("the id cannot match any user")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
			log.Println("error encoding data for a client")
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
		log.Println("error encoding data for a client")
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])

	if errParams != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error converting string to int: ", errParams)
		return
	}

	user, isEmpty := entities.Users[id]

	if !isEmpty {
		user.Error = errors.New("the id cannot match any user")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
			log.Println("error encoding data for a client")
			return
		}
		return
	}

	entities.DeleteUser(user.Id)
	w.WriteHeader(http.StatusNoContent)
}

func takeUserPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])

	if errParams != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error converting string to int: ", errParams)
		return
	}

	user, isEmpty := entities.Users[id]

	if !isEmpty {
		user.Error = errors.New("the id cannot match any user")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
			log.Println("error encoding data for a client")
			return
		}
		return
	}

	var points entities.Request

	if user.Error = json.NewDecoder(r.Body).Decode(&points); user.Error != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
			log.Println("error decoding client's data: ", user.Error)
			return
		}
		return
	}
	entities.UserTake(user.Id, points.Points)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
		log.Println("error encoding data for a client")
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			panic(errClose)
		}
	}()
}

func fundUserPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])

	if errParams != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error converting string to int: ", errParams)
		return
	}

	user, isEmpty := entities.Users[id]

	if !isEmpty {
		user.Error = errors.New("the id cannot match any user")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
			log.Println("error encoding data for a client")
			return
		}
		return
	}

	var points entities.Request
	if user.Error = json.NewDecoder(r.Body).Decode(&points); user.Error != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
			log.Println("error decoding client's data: ", user.Error)
			return
		}
		return
	}
	entities.UserFund(user.Id, points.Points)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if errAnswer := json.NewEncoder(w).Encode(user); errAnswer != nil {
		log.Println("error encoding data for a client")
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			panic(errClose)
		}
	}()
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
