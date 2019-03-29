package entities

import "errors"

type RequestPoints struct {
	Points int `json:"points"`
}

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
	Error   error  `json:"error"`
}

var Users = make(map[int]*User)
var UsersCounter = 0

func IsValid(user *User) bool {
	if user.Name == "" {
		user.Error = errors.New("wrong input, the name is not defined")
		return false
	} else if user.Balance < 300 {
		user.Error = errors.New("wrong input, not enough balance to register a user")
		return false
	}
	return true
}

//fixme
//what to do, cannot write the error in Error field (no if-construction)
func SaveUser(user *User, usersCounter *int) error {
	*usersCounter += 1
	user.Id = *usersCounter
	user.Balance -= 300
	Users[user.Id] = user
	return nil
}

//fixme
//what to do, cannot write the error in Error field (no if-construction)
func DeleteUser(id int) error {
	delete(Users, id)
	return nil
}

func UserTake(id, points int) error {
	if Users[id].Balance < points {
		Users[id].Error = errors.New("not enough balance to execute the request")
		return errors.New("error taking client's points")
	}
	Users[id].Balance -= points
	return nil
}

//fixme
//what to do, cannot write the error in Error field (no if-construction)
func UserFund(id, points int) error {
	Users[id].Balance += points
	return nil
}

