package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/yanrishbe/gaming-website/entities"
	"log"
	"net/http"
	"strconv"
)

func registerNewUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User

	if errDecode := json.NewDecoder(r.Body).Decode(&user); errDecode != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if errAnswer := json.NewEncoder(w).Encode(errDecode); errAnswer != nil {
			panic(errAnswer)
		}
		log.Println("Error decoding client's data: ", errDecode)
		return
	}

	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			panic(errClose)
		}
	}()

	if user.Name == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		_, errName := w.Write([]byte("Wrong input, please write your name"))
		if errName != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error sending message: ", errName)
			return
		}
	} else if *user.Balance < 300 || user.Balance == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		_, errBalance := w.Write([]byte("Not enough balance to register a user"))
		if errBalance != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error sending message: ", errBalance)
			return
		}
	}

	user = entities.CreateUser(user)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if errEncode := json.NewEncoder(w).Encode(user); errEncode != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if errAnswer := json.NewEncoder(w).Encode(errEncode); errAnswer != nil {
			panic(errAnswer)
		}
		log.Println("Error encoding data for a client: ", errEncode)
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	if errParams != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error converting string to int: ", errParams)
		return
	}

	for _, user := range entities.Users {
		if user.Id == id {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if errEncode := json.NewEncoder(w).Encode(user); errEncode != nil {
				w.WriteHeader(http.StatusInternalServerError)
				if errAnswer := json.NewEncoder(w).Encode(errEncode); errAnswer != nil {
					panic(errAnswer)
				}
				log.Println("Error encoding data for a client: ", errEncode)
				return
			}
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			_, errId := w.Write([]byte("There is no users matching the id"))
			if errId != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error sending message: ", errId)
				return
			}
		}
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	var doesExist = false
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	if errParams != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error converting string to int: ", errParams)
		return
	}

	for _, user := range entities.Users {
		if user.Id == id {
			entities.RemoveUser(user.Id)
			doesExist = true
			w.WriteHeader(http.StatusNoContent)
			break
		}
	}

	if doesExist == false {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		_, errId := w.Write([]byte("There is no users matching the id"))
		if errId != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error sending message: ", errId)
			return
		}
	}
}

func takeUserPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, errParams := strconv.Atoi(params["id"])
	if errParams != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error converting string to int: ", errParams)
		return
	}
	for _, user := range entities.Users {
		if user.Id == id {
			var pointsMap map[string]interface{}
			if errDecode := json.NewDecoder(r.Body).Decode(&pointsMap); errDecode != nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				if errAnswer := json.NewEncoder(w).Encode(errDecode); errAnswer != nil {
					panic(errAnswer)
				}
				log.Println("Error decoding client's data: ", errDecode)
				return
			}

			points,keyCheck := pointsMap["points"].(int)//TODO:
			if !keyCheck {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				if errAnswer := json.NewEncoder(w).Encode("Wrong input"); errAnswer != nil {
					panic(errAnswer)
				}
				log.Println("Wrong client's input")
				return
			}
			user, errTake := entities.UserTake(id, points)
			if errTake != nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				if errAnswer := json.NewEncoder(w).Encode(errTake); errAnswer != nil {
					panic(errAnswer)
				}
				log.Println("Wrong client's input: ", errTake)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if errEncode := json.NewEncoder(w).Encode(user); errEncode != nil {
				w.WriteHeader(http.StatusInternalServerError)
				if errAnswer := json.NewEncoder(w).Encode(errEncode); errAnswer != nil {
					panic(errAnswer)
				}
				log.Println("Error encoding data for a client: ", errEncode)
				return
			}
		}
	}
	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			panic(errClose)
		}
	}()
}

func fundUserPoints(w http.ResponseWriter, r *http.Request) {

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
