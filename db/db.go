//Package db implements main options of the game through a connection to database
package db

import (
	"errors"

	"github.com/yanrishbe/gaming-website/entities"
)

//DB struct stores users' data in UsersMap
type DB struct {
	UsersMap     map[int]*entities.User
	UsersCounter int
}

//SaveUser registers a new user
func (db *DB) SaveUser(user *entities.User) error {
	db.UsersCounter++
	user.ID = db.UsersCounter
	user.Balance -= 300
	db.UsersMap[user.ID] = user
	return nil
}

//DeleteUser removes a user from the UsersMap
func (db *DB) DeleteUser(id int) error {
	delete(db.UsersMap, id)
	return nil
}

//UserTake takes the requested amount of points from a user's balance
func (db *DB) UserTake(id, points int) error {
	if db.UsersMap[id].Balance < points {
		return errors.New("not enough balance to execute the request")
	}
	db.UsersMap[id].Balance -= points
	return nil
}

//UserFund adds the requested amount of points to user's balance
func (db *DB) UserFund(id, points int) error {
	db.UsersMap[id].Balance += points
	return nil
}

//New is used to create an instance of DB struct and initialize it
func New() *DB {
	db := new(DB)
	db.UsersMap = make(map[int]*entities.User)
	return db
}
