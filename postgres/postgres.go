package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/yanrishbe/gaming-website/entity"
)

type DB struct {
	db *sql.DB
}

func New() (DB, error) {
	connStr, ok := os.LookupEnv("CONN")
	if !ok {
		return DB{}, entity.DBErr(errors.New("empty connection string"))
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return DB{}, entity.DBErr(err)
	}
	err = db.Ping()
	gm := DB{db: db}
	if err != nil {
		return DB{}, entity.DBErr(err)
	}
	err = gm.CreateTables()
	if err != nil {
		return DB{}, entity.DBErr(err)
	}
	gm.db.SetMaxOpenConns(5)
	return gm, nil
}

func (db DB) CreateTables() error {
	_, err := db.db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		balance INT NOT NULL CHECK(balance>=0))`)
	if err != nil {
		return entity.DBErr(fmt.Errorf("table 'users' failed: %v", err))
	}
	_, err = db.db.Exec(`
		CREATE TABLE IF NOT EXISTS tournaments (
		id serial PRIMARY KEY,
		name TEXT NOT NULL,
		deposit INT NOT NULL CHECK(deposit>=0),
		prize INT NOT NULL CHECK(prize>=0) DEFAULT 0,
		finished BOOLEAN NOT NULL DEFAULT FALSE,
		winner_id INT NOT NULL)`)
	if err != nil {
		return entity.DBErr(fmt.Errorf("table 'tournaments' failed: %v", err))
	}
	_, err = db.db.Exec(`
		CREATE TABLE IF NOT EXISTS tournament_req (
		tournament_id INT NOT NULL REFERENCES tournaments (id) ON DELETE RESTRICT,
		user_id INT NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
		PRIMARY KEY (tournament_id, user_id) )`)
	if err != nil {
		return entity.DBErr(fmt.Errorf("table 'tournament_req' failed: %v", err))
	}
	return nil
}

func (db DB) CreateUser(u entity.User) (entity.User, error) {
	err := db.db.QueryRow(`
		INSERT INTO users (name, balance)
		VALUES ($1, $2)
 		RETURNING id`, u.Name, u.Balance).Scan(&u.ID)
	if err != nil {
		return u, entity.DBErr(err)
	}
	return u, nil
}

func (db DB) GetUser(id int) (entity.User, error) {
	if id <= 0 {
		return entity.User{}, entity.InvIDErr(errors.New("expected id > 0"))
	}
	u := entity.User{}
	err := db.db.QueryRow(`
		SELECT id, name, balance 
		FROM users 
		WHERE id = $1`, id).Scan(&u.ID, &u.Name, &u.Balance)
	if err == sql.ErrNoRows {
		return u, entity.UserNotFoundErr(err)
	} else if err != nil {
		return u, entity.DBErr(err)
	}
	return u, nil
}

func (db DB) DelUser(id int) error {
	u, err := db.GetUser(id)
	if err != nil {
		return err
	}
	_, err = db.db.Exec(`
		DELETE FROM users 
		WHERE id = $1`, u.ID)
	if err != nil {
		return entity.DBErr(err)
	}
	return nil
}

func (db DB) TakePoints(id, points int) (entity.User, error) {
	u, err := db.GetUser(id)
	if err != nil {
		return u, err
	}
	err = db.db.QueryRow(`
		UPDATE users
		SET balance = balance - $1
		WHERE id = $2
		RETURNING balance`, points, u.ID).Scan(&u.Balance)
	if err != nil {
		return u, entity.DBErr(err)
	}
	return u, nil
}

func (db DB) FundPoints(id, points int) (entity.User, error) {
	u, err := db.GetUser(id)
	if err != nil {
		return u, err
	}
	err = db.db.QueryRow(`
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2
		RETURNING balance`, points, u.ID).Scan(&u.Balance)
	if err != nil {
		return u, entity.DBErr(err)
	}
	return u, nil
}
