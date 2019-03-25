package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	. "github.com/yanrishbe/gaming-website/entities"
	"log"
	"net/http"
)

//POST /user
//Request
//{
//"name" :  name,
//"balance": 1000
//}
//Response:
//{
//"id": 1,
//"name" :  name,
//"balance": 700
//}

func registerNewUser(w http.ResponseWriter, r *http.Request){
var user User
 errDecode := json.NewDecoder(r.Body).Decode(&user)
 if errDecode != nil {
	 w.WriteHeader(http.StatusInternalServerError)
	 log.Println("Error decoding client's data: ", errDecode)
	 return
 }
 user.Id = UsersCounter;UsersCounter++
 user.Balance -= 300

}

func getUser(w http.ResponseWriter, r *http.Request){

}

func deleteUser(w http.ResponseWriter, r *http.Request){

}

func takeUserPoints(w http.ResponseWriter, r *http.Request){

}

func fundUserPoints(w http.ResponseWriter, r *http.Request){

}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/user", registerNewUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", getUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", deleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/user/{id}/take", takeUserPoints).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}/fund", fundUserPoints).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))
}
