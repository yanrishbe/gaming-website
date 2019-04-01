//Package db implements main options of the game through a connection to database
package db

import (
	"errors"

	"github.com/yanrishbe/gaming-website/entities"
)

//UsersMap collects all user's accounts in a map, which key is a user's id
type UsersMap map[int]*entities.User

//UsersCounter is necessary for incrementing user's id
var UsersCounter = 0

//SaveUser registers a new user
func (usersMap *UsersMap) SaveUser(user *entities.User, usersCounter *int) error {
	*usersCounter++
	user.ID = *usersCounter
	user.Balance -= 300
	(*usersMap)[user.ID] = user
	return nil
}

//DeleteUser removes a user from Users map
func (usersMap *UsersMap) DeleteUser(id int) error {
	delete(*usersMap, id)
	return nil
}

//UserTake takes the requested amount of points from a user's balance
func (usersMap *UsersMap) UserTake(id, points int) error {
	if (*usersMap)[id].Balance < points {
		return errors.New("not enough balance to execute the request")
	}
	(*usersMap)[id].Balance -= points
	return nil
}

//UserFund adds the requested amount of points to user's balance
func (usersMap *UsersMap) UserFund(id, points int) error {
	(*usersMap)[id].Balance += points
	return nil
}
