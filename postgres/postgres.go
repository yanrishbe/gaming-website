package postgres

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"
	"github.com/yanrishbe/gaming-website/entity"
)

type DB struct {
	db *sql.DB
}

func setConnStr() (string, bool) {
	conn, ok := os.LookupEnv("CONN")
	return conn, ok
}

func New() (DB, error) {
	connStr, ok := setConnStr()
	if !ok {
		return DB{}, entity.DBErr(errors.New("empty connection string"))
	}
	db, err := sql.Open("postgres", connStr)
	//db, err := sql.Open("postgres", "user=postgres dbname=gaming_website password=docker2147 host=localhost port=5432 sslmode=disable")
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
		return entity.DBErr(err)
	}
	return nil
}

func (db DB) RegUser(u entity.User) (int, error) {
	err := u.IsValid()
	if err != nil {
		return 0, err
	}
	err = db.db.QueryRow(`
		INSERT INTO users (name, balance)
		VALUES ($1, $2 - 300)
 		RETURNING id`, u.Name, u.Balance).Scan(&u.ID)
	if err != nil {
		return 0, entity.DBErr(err)
	}
	return u.ID, nil
}

func (db DB) GetUser(id int) (entity.User, error) {
	if id <= 0 {
		return entity.User{}, entity.InvIDErr(errors.New("expected id > 0"))
	}
	u := entity.User{}
	err := db.db.QueryRow(`
		SELECT id, name, balance 
		FROM users WHERE id = $1`, id).Scan(&u.ID, &u.Name, &u.Balance)
	if err == sql.ErrNoRows {
		return u, entity.UserNotFoundErr(err)
	} else if err != nil {
		return u, entity.DBErr(err)
	}
	return u, nil
}
