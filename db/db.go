// Package db implements main options of the game through a connection to database
package db

import (
	"errors"
	"sync"

	"github.com/yanrishbe/gaming-website/entities"
)

// DB struct stores users' data in UsersMap
type DB struct {
	mutex        *sync.RWMutex
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
func (db *DB) SaveUser(user *entities.User) (int, error) {
	if !canRegister(*user) {
		return 0, errors.New("user's data is not valid")
	}
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.UsersCounter++
	//user.ID = db.UsersCounter
	user.Balance -= 300
	db.UsersMap[db.UsersCounter] = user
	return db.UsersCounter, nil
}

func (db *DB) GetUser(id int) (*entities.User, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	user, doesExist := db.UsersMap[id]
	if !doesExist {
		return nil, errors.New("the id cannot match any user")
	}
	return user, nil
}

// DeleteUser removes a user from the UsersMap
func (db *DB) DeleteUser(id int) error {
	if _, errGet := db.GetUser(id); errGet != nil {
		return errGet
	}
	db.mutex.Lock()
	defer db.mutex.Unlock()
	delete(db.UsersMap, id)
	return nil
}

// UserTake takes the requested amount of points from a user's balance
func (db *DB) UserTake(id, points int) error {
	if _, errGet := db.GetUser(id); errGet != nil {
		return errGet
	}
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if db.UsersMap[id].Balance < points {
		return errors.New("not enough balance to execute the request")
	}
	db.UsersMap[id].Balance -= points
	return nil
}

// UserFund adds the requested amount of points to user's balance
func (db *DB) UserFund(id, points int) error {
	if _, errGet := db.GetUser(id); errGet != nil {
		return errGet
	}
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.UsersMap[id].Balance += points
	return nil
}

// New is used to create an instance of DB struct and initialize it
func New() *DB {
	return &DB{
		mutex:    &sync.RWMutex{},
		UsersMap: make(map[int]*entities.User),
	}
}

func (db *DB) CountUsers() int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return len(db.UsersMap)
}
