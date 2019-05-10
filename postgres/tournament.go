package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/yanrishbe/gaming-website/entity"
)

func (db DB) CreateTourn(t entity.Tournament) (entity.Tournament, error) {
	err := db.db.QueryRow(`
		INSERT INTO tournaments (name, deposit)
		VALUES ($1, $2)
 		RETURNING id`, t.Name, t.Deposit).Scan(&t.ID)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't create tournament: %v", err))
	}
	t.Status = entity.Active
	return t, nil
}

func (db DB) ValidGetTourn(id int) (bool, error) {
	if id <= 0 {
		return false, entity.InvIDErr(errors.New("expected id greater than 0"))
	}
	var finished bool
	err := db.db.QueryRow(`
		SELECT finished
		FROM tournaments
		WHERE id=$1`, id).Scan(&finished)
	if err == sql.ErrNoRows {
		return false, entity.ReqErr(fmt.Errorf("tournament doesn't exist: %v", err))
	} else if err != nil {
		return false, entity.DBErr(err)
	}
	return finished, nil
}

func (db DB) GetTourn(id int) (entity.Tournament, error) {
	var t entity.Tournament
	err := db.db.QueryRow(`
		SELECT id, name, deposit, prize 
		FROM tournaments 
		WHERE id = $1 AND finished = FALSE`,
		id).Scan(&t.ID, &t.Name, &t.Deposit, &t.Prize)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't get tournament: %v", err))
	}

	rows, err := db.db.Query(`
		SELECT users.id, users.name
		FROM tournament_req
		INNER JOIN users ON tournament_req.user_id = users.id
		WHERE tournament_req.tournament_id = $1`, id)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't get tournament data: %v", err))
	}
	defer rows.Close()

	var u entity.UserTourn
	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Name)
		if err != nil {
			return t, entity.DBErr(fmt.Errorf("can't get tournament data: %v", err))
		}
		t.Users = append(t.Users, u)
	}
	err = rows.Err()
	if err != nil {
		return t, entity.DBErr(err)
	}
	t.Status = entity.Active
	return t, nil
}

func (db DB) GetTournFinished(id int) (entity.TournFinished, error) {
	var t entity.TournFinished
	err := db.db.QueryRow(`
		SELECT id, name, prize , winner_id
		FROM tournaments 
		WHERE id = $1`,
		id).Scan(&t.ID, &t.Name, &t.Prize, &t.Winner)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't get tournament: %v", err))
	}

	rows, err := db.db.Query(`
		SELECT users.id, users.name
		FROM tournament_req
		INNER JOIN users ON tournament_req.user_id = users.id
		WHERE tournament_req.tournament_id = $1`, id)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't get tournament data: %v", err))
	}
	defer rows.Close()

	var w entity.Winner
	for rows.Next() {
		err := rows.Scan(&w.ID, &w.Name)
		if w.ID == t.Winner {
			w.Winner = true
		}
		if err != nil {
			return t, entity.DBErr(fmt.Errorf("can't get tournament data: %v", err))
		}
		t.Users = append(t.Users, w)
	}

	err = rows.Err()
	if err != nil {
		return t, entity.DBErr(err)
	}
	t.Status = entity.Finished
	return t, nil
}

func (db DB) JoinTourn(tID, uID int) (entity.Tournament, error) {
	var t entity.Tournament
	t.ID = tID

	tx, err := db.db.Begin()
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("transaction error: %v", err))
	}
	defer tx.Rollback()

	err = tx.QueryRow(`
		SELECT deposit 
		FROM tournaments 
		WHERE id = $1 AND finished = FALSE`, tID).Scan(&t.Deposit)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't get tournament deposit: %v", err))
	}

	_, err = tx.Exec(`
		UPDATE users
		SET balance = balance - $1
		WHERE id = $2`, t.Deposit, uID)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't update user's balance: %v", err))
	}

	_, err = tx.Exec(`
		INSERT INTO tournament_req (tournament_id, user_id)
		VALUES ($1, $2)`, tID, uID)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't register a user: %v", err))
	}

	err = tx.QueryRow(`
		UPDATE tournaments
		SET prize = prize + $1
		WHERE id = $2
		RETURNING name, prize`, t.Deposit, t.ID).Scan(&t.Name, &t.Prize)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't update the prize: %v", err))
	}

	err = tx.Commit()
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("transaction error: %v", err))
	}
	return t, nil
}

func (db DB) ValidJoin(tID, uID int) error {
	var finished bool
	err := db.db.QueryRow(`
		SELECT finished
		FROM tournaments
		WHERE id=$1`, tID).Scan(&finished)
	if err == sql.ErrNoRows {
		return entity.ReqErr(fmt.Errorf("tournament doesn't exist: %v", err))
	} else if err != nil {
		return entity.DBErr(err)
	}
	if finished {
		return entity.ReqErr(errors.New("the tournament is finished"))
	}

	var name string
	err = db.db.QueryRow(`
		SELECT name
		FROM users
		WHERE id = $1`, uID).Scan(&name)
	if err == sql.ErrNoRows {
		return entity.ReqErr(errors.New("the user doesn't exist"))
	} else if err != nil {
		return entity.DBErr(err)
	}

	var id int
	err = db.db.QueryRow(`
		SELECT tournament_id
		FROM tournament_req 
		WHERE user_id = $1 AND tournament_id = $2`, uID, tID).Scan(&id)
	if err == nil {
		return entity.RegErr(fmt.Errorf("user is already registered"))
	} else if id != 0 {
		return entity.RegErr(fmt.Errorf("user is already registered"))
	}

	rows, err := db.db.Query(`
		SELECT tournament_id
		FROM tournament_req
		WHERE user_id = $1`, uID)
	if err != nil {
		return entity.DBErr(fmt.Errorf("can't get tournament data: %v", err))
	}
	if rows.Next() {
		return entity.ReqErr(errors.New("trial to join more than one tournament"))
	}
	err = rows.Err()
	if err != nil {
		return entity.DBErr(err)
	}
	defer rows.Close()

	var deposit int
	err = db.db.QueryRow(`
		SELECT deposit 
		FROM tournaments 
		WHERE id = $1`, tID).Scan(&deposit)
	if err != nil {
		entity.DBErr(err)
	}
	var balance int
	err = db.db.QueryRow(`
		SELECT balance 
		FROM users 
		WHERE id = $1`, uID).Scan(&balance)
	if err != nil {
		entity.DBErr(err)
	}
	if balance < deposit {
		return entity.RegErr(errors.New("balance is lower than deposit"))
	}
	return nil
}

func (db DB) ValidFinish(tID int) (bool, error) {
	var finished bool

	err := db.db.QueryRow(`
		SELECT finished
		FROM tournaments
		WHERE id=$1`, tID).Scan(&finished)
	if err == sql.ErrNoRows {
		return false, entity.ReqErr(fmt.Errorf("tournament doesn't exist: %v", err))
	} else if err != nil {
		return false, entity.DBErr(err)
	}

	rows, err := db.db.Query(`
		SELECT user_id
		FROM tournament_req
		WHERE tournament_id = $1`, tID)
	if err != nil {
		return false, entity.DBErr(fmt.Errorf("can't get tournament data: %v", err))
	}
	if !rows.Next() {
		return false, entity.ReqErr(errors.New("can't finish, no users"))
	}
	err = rows.Err()
	if err != nil {
		return false, entity.DBErr(err)
	}
	defer rows.Close()
	return finished, nil
}

func (db DB) TournUsers(tID int) (func() int, error) {
	rows, err := db.db.Query(`
		SELECT user_id
		FROM tournament_req
		WHERE tournament_id = $1`, tID)
	if err != nil {
		return nil, entity.DBErr(fmt.Errorf("can't get data: %v", err))
	}
	defer rows.Close()
	var userID []int
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, entity.DBErr(fmt.Errorf("can't get data: %v", err))
		}
		userID = append(userID, id)
	}
	err = rows.Err()
	if err != nil {
		return nil, entity.DBErr(fmt.Errorf("rows error: %v", err))
	}
	return func() int {
		rand.Seed(time.Now().UTC().UnixNano())
		r := rand.Intn(len(userID))
		return userID[r]
	}, nil
}

func (db DB) FinishTourn(tID, uID int) error {
	tx, err := db.db.Begin()
	if err != nil {
		return entity.DBErr(fmt.Errorf("transaction error: %v", err))
	}
	defer tx.Rollback()
	var prize int
	err = tx.QueryRow(`
		UPDATE tournaments
		SET winner_id = $1, finished = $2
		WHERE id = $3
		RETURNING prize`, uID, true, tID).Scan(&prize)
	if err != nil {
		return entity.DBErr(err)
	}

	_, err = tx.Exec(`
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2`, prize, uID)
	if err != nil {
		return entity.DBErr(err)
	}

	err = tx.Commit()
	if err != nil {
		return entity.DBErr(fmt.Errorf("transaction error: %v", err))
	}
	return nil
}

func (db DB) DelTourn(id int) error {
	return nil
}
