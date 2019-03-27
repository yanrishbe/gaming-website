package entities

import "errors"

type Request struct {
	Points int `json:"points"`
}

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
	Error   error  `json:"error"`
}

var Users = make(map[int]*User)
var usersCounter = 0

func IsValid(user *User) {
	if user.Name == "" {
		user.Error = errors.New("wrong input, the name is not defined")
		return
	} else if user.Balance < 300 {
		user.Error = errors.New("wrong input, not enough balance to register a user")
		return
	}
}

func SaveUser(user *User) {
	user.Id = usersCounter + 1
	user.Balance -= 300
	Users[user.Id] = user
}

func DeleteUser(id int) {
	delete(Users, id)
}

func UserTake(id, points int) {
	if Users[id].Balance < points {
		Users[id].Error = errors.New("not enough balance to execute the request")
		return
	}
	Users[id].Balance -= points
}

func UserFund(id, points int){
	Users[id].Balance += points
}

