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
	UsersMap     map[int]entities.User
	UsersCounter int
}

func canRegister(user entities.User) bool {
	if user.Name == "" || user.Balance < 300 {
		return false
	}
	return true
}

// SaveUser registers a new user
func (db *DB) SaveUser(usr entities.User) (entities.User, error) {
	us := usr
	if !canRegister(us) {
		return us, errors.New("user's data is not valid")
	}
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.UsersCounter++
	us.ID = db.UsersCounter
	us.Balance -= 300
	db.UsersMap[db.UsersCounter] = us
	return us, nil
}

func (db *DB) GetUser(id int) (entities.User, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	us, doesExist := db.UsersMap[id]
	if !doesExist {
		return us, errors.New("the id cannot match any user")
	}
	return us, nil
}

// DeleteUser removes a user from the UsersMap
func (db *DB) DeleteUser(id int) error {
	_, err := db.GetUser(id)
	if err != nil {
		return err
	}
	db.mutex.Lock()
	defer db.mutex.Unlock()
	delete(db.UsersMap, id)
	return nil
}

// UserTake takes the requested amount of points from a user's balance
func (db *DB) UserTake(id, points int) (entities.User, error) {
	us, err := db.GetUser(id)
	if err != nil {
		return us, err
	}
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if us.Balance < points {
		return us, errors.New("not enough balance to execute the request")
	}
	us.Balance -= points
	db.UsersMap[id] = us
	return us, nil
}

// UserFund adds the requested amount of points to user's balance
func (db *DB) UserFund(id, points int) (entities.User, error) {
	us, err := db.GetUser(id)
	if err != nil {
		return us, err
	}
	db.mutex.Lock()
	defer db.mutex.Unlock()
	us.Balance += points
	db.UsersMap[id] = us
	return us, nil
}

// New is used to create an instance of DB struct and initialize it
func New() *DB {
	return &DB{
		mutex:    &sync.RWMutex{},
		UsersMap: make(map[int]entities.User),
	}
}

// CountUsers returns the amount of elements in the db
func (db *DB) CountUsers() int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return len(db.UsersMap)
}
