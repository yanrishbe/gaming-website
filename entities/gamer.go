//Package entities implements methods for storing and writing user's data and handling user's commands
package entities

import "errors"

//RequestPoints represents a struct to send take and fund requests to the gaming website
type RequestPoints struct {
	Points int `json:"points"`
}

//User struct represents a struct necessary for storing and changing user's data
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
	Error   string `json:"error"`
}

//Users collects all user's accounts in a map, which key is a user's id
var Users = make(map[int]*User)

//UsersCounter is necessary for incrementing user's id
var UsersCounter = 0

//IsValid checks whether it is possible to register a new user or not
func IsValid(user *User) bool {
	if user.Name == "" || user.Balance < 300 {
		return false
	}
	return true
}

//SaveUser registers a new user
func SaveUser(user *User, usersCounter *int) error {
	*usersCounter++
	user.ID = *usersCounter
	user.Balance -= 300
	Users[user.ID] = user
	return nil
}

//DeleteUser removes a user from Users map
func DeleteUser(id int) error {
	if Users[id].Error != "" {
		Users[id].Error = ""
	}
	delete(Users, id)
	return nil
}

//UserTake takes the requested amount of points from a user's balance
func UserTake(id, points int) error {
	if Users[id].Error != "" {
		Users[id].Error = ""
	}
	if Users[id].Balance < points {
		return errors.New("not enough balance to execute the request")
	}
	Users[id].Balance -= points
	return nil
}

//UserFund adds the requested amount of points to user's balance
func UserFund(id, points int) error {
	//на входе можем получить старую ошибку из userTake
	if Users[id].Error != "" {
		Users[id].Error = ""
	}
	Users[id].Balance += points
	return nil
}
