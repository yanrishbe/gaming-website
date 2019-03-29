package entities

import "errors"

type RequestPoints struct {
	Points int `json:"points"`
}

type ResponseDelete struct {
	Error error
}

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
	Error   string `json:"error"`
}

var Users = make(map[int]*User)
var UsersCounter = 0

func IsValid(user *User) bool {
	if user.Name == "" || user.Balance < 300 {
		return false
	}
	return true
}

func SaveUser(user *User, usersCounter *int) error {
	*usersCounter += 1
	user.Id = *usersCounter
	user.Balance -= 300
	Users[user.Id] = user
	return nil
}

func DeleteUser(id int) error {
	if Users[id].Error != "" {
		Users[id].Error = ""
	}
	delete(Users, id)
	return nil
}

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

//на входе можем получить старую ошибку из userTake
func UserFund(id, points int) error {
	if Users[id].Error != "" {
		Users[id].Error = ""
	}
	Users[id].Balance += points
	return nil
}
