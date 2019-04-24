// Package db implements main options of the game through a connection to database
package db

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"github.com/yanrishbe/gaming-website/entity"
)

// DB manages users' data using Postgres
type DB struct {
	db *sql.DB
}

// regarding DB and other long-running functions - it's always advisable to pass context.Context as a first parameter to
// the long-running function,  and do all child calls using this context, for example
// db.PingContext(ctx, ...), db.QueryContext(ctx, ...), db.ExecContext(ctx, ...)
// because for example DB call can hang forever and your handler can hang. It's also a nice way to
// distinguish external calls from internal calls. When you see ctx.Context - you know that this function may block
// for really long time.
// It's not a problem, and you don't have to do that. Just a note for your information:)

func New(connStr string) (DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return DB{}, entity.DBErr(err)
	}
	err = db.Ping()
	if err != nil {
		return DB{}, entity.DBErr(err)
	}
	gm := DB{db: db}
	err = gm.CreateTables()
	if err != nil {
		return DB{}, entity.DBErr(err)
	}
	gm.db.SetMaxOpenConns(5)
	return gm, nil
}

func (gm DB) CreateTables() error {
	_, err := gm.db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		balance INT NOT NULL CHECK(balance >= 0))`)
	if err != nil {
		return entity.DBErr(err)
	}
	return nil
}

func (gm DB) Close() error {
	return gm.db.Close()
}

func (gm DB) SaveUser(u entity.User) (entity.User, error) {
	err := u.CanRegister()
	if err != nil {
		return u, err
	}
	err = gm.db.QueryRow(`
			INSERT INTO users (name, balance) 
			VALUES ($1, $2 - 300) 
			RETURNING id, balance`, u.Name, u.Balance).Scan(&u.ID, &u.Balance)
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
	if err == sql.ErrNoRows {
		return u, entity.UserNotFoundErr(err)
	}
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

func (gm DB) DeleteUser(id int) error {
	u, err := gm.GetUser(id)
	if err != nil {
		return err
	}
	_, err = gm.db.Exec("DELETE FROM users WHERE id = $1", u.ID)
	if err != nil {
		return entity.DBErr(err)
	}
	return nil
}

func (gm DB) CountUsers() (int, error) {
	row := gm.db.QueryRow("SELECT COUNT(id) FROM users")
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, entity.DBErr(err)
	}
	return count, nil
}
