// Package db implements main options of the game through a connection to database
package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/yanrishbe/gaming-website/entity"
)

// DB struct stores users' data in UsersMap
type DB struct {
	db *sql.DB
}

func New() (DB, error) {
	connStr := "user=postgres password=docker2147 dbname=gaming_website host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return DB{}, entity.DBErr(err)
	}
	err = db.Ping()
	if err != nil {
		return DB{}, entity.DBErr(err)
	}
	gm := DB{db: db}
	err = gm.createTables()
	if err != nil {
		return DB{}, entity.DBErr(err)
	}
	gm.db.SetMaxOpenConns(20)
	return gm, nil
}

func (gm DB) createTables() error {
	_, err := gm.db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		balance INT NOT NULL CHECK(balance >= 0)`)
	if err != nil {
		return entity.DBErr(err)
	}
	return nil
}

func (gm DB) Close() error {
	return gm.db.Close()
}

func (gm DB) Register(u entity.User) (entity.User, error) {
	err := u.CanRegister()
	if err != nil {
		return u, err
	}
	_, err = gm.db.Exec("INSERT INTO users (name, balance) VALUES ($1, $2 - 300)", u.Name, u.Balance)
	if err != nil {
		return u, entity.DBErr(err)
	}
	err = gm.db.QueryRow("SELECT id, balance FROM users WHERE name = $1", u.Name).Scan(&u.ID, &u.Balance)
	if err != nil {
		return u, entity.DBErr(err)
	}
	return u, nil
}

func (gm DB) GetUser(id int) (entity.User, error) {
	if id <= 0 {
		return entity.User{}, entity.InvIDErr(errors.New("expected id greater than 0"))
	}
	u := entity.User{}
	err := gm.db.QueryRow(`SELECT id, name, balance FROM users 
		WHERE id = $1`, id).Scan(&u.ID, &u.Name, &u.Balance)
	return u, err
}

func (gm DB) UserTake(id, points int) (entity.User, error) {
	u, err := gm.GetUser(id)
	if err != nil {
		return u, err
	}
	_, err = gm.db.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", points, u.ID)
	if err != nil {
		return u, entity.DBErr(err)
	}
	err = gm.db.QueryRow("SELECT id, name, balance FROM users WHERE id = $1", u.ID).Scan(&u.ID, &u.Name, &u.Balance)
	if err != nil {
		return u, entity.DBErr(err)
	}
	return u, nil
}

func (gm DB) UserFund(id, points int) (entity.User, error) {
	u, err := gm.GetUser(id)
	if err != nil {
		return u, err
	}
	_, err = gm.db.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", points, u.ID)
	if err != nil {
		return u, entity.DBErr(err)
	}
	err = gm.db.QueryRow("SELECT id, name, balance FROM users WHERE id = $1", u.ID).Scan(&u.ID, &u.Name, &u.Balance)
	if err != nil {
		return u, entity.DBErr(err)
	}
	return u, nil
}

//FIXME////////////////////////////////////////////////////////////
func (gm DB) Delete(id int) error {
	_, err := gm.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("couldn't update user's data: %v", err)
	}
	err = ud.tx.Commit()
	if err != nil {
		return fmt.Errorf("transaction error: %v", err)
	}
	return nil
}
