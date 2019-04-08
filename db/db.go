// Package db implements main options of the game through a connection to database
package db

import (
	"errors"
	"sync"

	"github.com/yanrishbe/gaming-website/entities"
)

// DB struct stores users' data in UsersMap
type DB struct {
	mutex        *sync.Mutex
	UsersMap     map[int]*entities.User
	UsersCounter int
}

func canRegister(user entities.User) bool {
	if user.Name == "" || user.Balance < 300 {
		return false
	}
	return true
}

// SaveUser registers a new user
func (db *DB) SaveUser(user *entities.User) error {
	if !canRegister(*user) {
		return errors.New("user's data is not valid")
	}
	db.mutex.Lock()
	db.UsersCounter++
	user.ID = db.UsersCounter
	user.Balance -= 300
	db.UsersMap[user.ID] = user
	db.mutex.Unlock()
	return nil
}

// DeleteUser removes a user from the UsersMap
func (db *DB) DeleteUser(id int) error {
	_, doesExist := db.UsersMap[id]
	if !doesExist {
		return errors.New("the id cannot match any user")
	}
	db.mutex.Lock()
	delete(db.UsersMap, id)
	db.mutex.Unlock()
	return nil
}

// UserTake takes the requested amount of points from a user's balance
func (db *DB) UserTake(id, points int) error {
	_, doesExist := db.UsersMap[id]
	if !doesExist {
		return errors.New("the id cannot match any user")
	}
	if db.UsersMap[id].Balance < points {
		return errors.New("not enough balance to execute the request")
	}
	db.mutex.Lock()
	db.UsersMap[id].Balance -= points
	db.mutex.Unlock()
	return nil
}

// UserFund adds the requested amount of points to user's balance
func (db *DB) UserFund(id, points int) error {
	_, doesExist := db.UsersMap[id]
	if !doesExist {
		return errors.New("the id cannot match any user")
	}
	db.mutex.Lock()
	db.UsersMap[id].Balance += points
	db.mutex.Unlock()
	return nil
}

// New is used to create an instance of DB struct and initialize it
func New() *DB {
	return &DB{
		mutex:    &sync.Mutex{},
		UsersMap: make(map[int]*entities.User),
	}
}
