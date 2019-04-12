// Package db implements main options of the game through a connection to database
package db

import (
	"errors"
	"sync"

	"github.com/yanrishbe/gaming-website/entity"
)

// DB struct stores users' data in UsersMap
type DB struct {
	mutex        *sync.RWMutex
	UsersMap     map[int]entity.User
	UsersCounter int
}

func canRegister(user entity.User) error {
	if user.Name == "" {
		return errors.New("user's name is empty")
	} else if user.Balance < 300 {
		return errors.New("user has got not enough points")
	}
	return nil
}

// SaveUser registers a new user
func (db *DB) SaveUser(usr entity.User) (entity.User, error) {
	us := usr
	err := canRegister(us)
	if err != nil {
		return us, entity.DBRegisterErr(err)
	}
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.UsersCounter++
	us.ID = db.UsersCounter
	us.Balance -= 300
	db.UsersMap[db.UsersCounter] = us
	return us, nil
}

func (db *DB) GetUser(id int) (entity.User, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	us, doesExist := db.UsersMap[id]
	if !doesExist {
		return us, entity.UserNotFoundError(errors.New("the id cannot match any user"))
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
func (db *DB) UserTake(id, points int) (entity.User, error) {
	us, err := db.GetUser(id)
	if err != nil {
		return us, err
	}
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if us.Balance < points {
		return us, entity.FewBal(errors.New("not enough balance to execute the request"))
	}
	us.Balance -= points
	db.UsersMap[id] = us
	return us, nil
}

// UserFund adds the requested amount of points to user's balance
func (db *DB) UserFund(id, points int) (entity.User, error) {
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
		UsersMap: make(map[int]entity.User),
	}
}

// CountUsers returns the amount of elements in the db
func (db *DB) CountUsers() int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return len(db.UsersMap)
}
