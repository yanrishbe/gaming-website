package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/yanrishbe/gaming-website/entity"
)

type DB struct {
	db *sql.DB
}

func SetConnStr() string {
	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		user = "postgres"
	}
	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		password = "docker2147"
	}
	dbname, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		dbname = "gaming_website"
	}
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = "localhost"
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "5432"
	}
	sslmode, ok := os.LookupEnv("SSLMODE")
	if !ok {
		sslmode = "disable"
	}

	return fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v sslmode=%v", user, password, dbname, host, port, sslmode)
}

func New() (DB, error) {
	connStr := SetConnStr()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return DB{}, err
	}
	err = db.Ping()
	gm := DB{db: db}
	if err != nil {
		return DB{}, err
	}
	err = gm.CreateTables()
	if err != nil {
		return DB{}, err
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
		return err
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
		return 0, err
	}
	return u.ID, nil
}

func (db DB) GetUser(id int) (entity.User, error) {
	if id <= 0 {
		return entity.User{}, errors.New("expected id greater than 0")
	}
	u := entity.User{}
	err := db.db.QueryRow(`
		SELECT id, name, balance 
		FROM users WHERE id = $1`, id).Scan(&u.ID, &u.Name, &u.Balance)
	if err == sql.ErrNoRows {
		return u, err
	}
	return u, err
}
