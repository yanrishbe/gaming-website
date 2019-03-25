package entities

type User struct {
	Id int `json:"id, omitempty"`
	Name string `json:"name, omitempty"`
	Balance int `json:"balance"`
	Points int `json:"points"`
}
 var Users []User
var UsersCounter = 1