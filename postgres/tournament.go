package postgres

import (
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
		return t, entity.DBErr(err)
	}
	return t, nil
}

func (db DB) GetTourn(id int) (entity.Tournament, error) {
	if id <= 0 {
		return entity.Tournament{}, entity.InvIDErr(errors.New("expected id > 0"))
	}
	var t entity.Tournament
	err := db.db.QueryRow(`
		SELECT id, name, deposit, prize 
		FROM tournaments 
		WHERE id = $1 AND finished = FALSE`,
		id).Scan(&t.ID, &t.Name, &t.Deposit, &t.Prize)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't get data: %v", err))
	}
	rows, err := db.db.Query(`
		SELECT users.id, users.name
		FROM tournament_req
		INNER JOIN users ON tournament_req.user_id = users.id
		WHERE tournament_req.tournament_id = $1`, id)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't get data: %v", err))
	}
	defer rows.Close()
	var u entity.UserTourn
	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Name)
		if err != nil {
			return t, entity.DBErr(fmt.Errorf("can't get data: %v", err))
		}
		t.Users = append(t.Users, u)
	}
	err = rows.Err()
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("rows error: %v", err))
	}
	t.Status = entity.Active
	return t, nil
}

func (db DB) JoinTourn(tID, uID int) (entity.Tournament, error) {
	var t entity.Tournament
	t.ID = tID

	tx, err := db.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return t, entity.DBErr(err)
	}

	err = tx.QueryRow(`
		SELECT deposit 
		FROM tournaments 
		WHERE id = $1 AND finished = FALSE`, tID).Scan(&t.Deposit)
	if err != nil {
		return t, entity.DBErr(err)
	}

	_, err = tx.Exec(`
		UPDATE users
		SET balance = balance - $1
		WHERE id = $2`, t.Deposit, uID)
	if err != nil {
		return t, entity.DBErr(err)
	}

	_, err = tx.Exec(`
		INSERT INTO tournament_req (tournament_id, user_id)
		VALUES ($1, $2)`, tID, uID)
	if err != nil {
		return t, entity.DBErr(err)
	}

	err = tx.QueryRow(`
		UPDATE tournaments
		SET prize = prize + $1
		WHERE id = $2
		RETURNING name, prize`, t.Deposit, t.ID).Scan(&t.Name, &t.Prize)

	err = tx.Commit()
	if err != nil {
		return t, entity.DBErr(err)
	}
	return t, nil
}

func (db DB) ValidJoin(tID, uID int) error {
	var deposit int
	err := db.db.QueryRow(`
		SELECT deposit 
		FROM tournaments 
		WHERE id = $1 AND finished = FALSE`, tID).Scan(&deposit)
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
		return entity.RegErr(fmt.Errorf("balance is lower than deposit %v", err))
	}
	var id int
	err = db.db.QueryRow(`
		SELECT id
		FROM tournament_req 
		WHERE id = $1`, uID).Scan(&id)
	if err == nil {
		return entity.RegErr(fmt.Errorf("user is already registered"))
	}
	var finished bool
	err = db.db.QueryRow(`
		SELECT finished
		FROM tournaments
		WHERE id = $1`, tID).Scan(&finished)
	if finished {
		return entity.RegErr(fmt.Errorf("the tournament is finished"))
	}
	return nil
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
		r := rand.Intn(len(userID) + 1)
		return userID[r]
	}, nil
}

func (db DB) FinishTourn(tID, uID int) (entity.TournFinished, error) {
	var t entity.TournFinished

	tx, err := db.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return t, entity.DBErr(err)
	}

	err = tx.QueryRow(`
		INSERT INTO tournaments (winner_id, finished)
		VALUES ($1, $2)
		WHERE id=$3
 		RETURNING name, prize`, uID, true, tID).Scan(&t.Name, &t.Prize)
	if err != nil {
		return t, entity.DBErr(err)
	}

	t.ID = tID
	t.Winner = uID
	t.Status = entity.Finished

	_, err = tx.Exec(`
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2`, t.Prize, t.Winner)
	if err != nil {
		return t, entity.DBErr(err)
	}
	return t, nil
}

func (db DB) ValidFinish(tID int) error {
	var finished bool
	err := db.db.QueryRow(`
		SELECT finished
		FROM tournaments
		WHERE id=$1`, tID).Scan(&finished)
	if err != nil {
		return entity.DBErr(err)
	}
	if finished {
		return entity.FinishErr(errors.New("tournament is already finished"))
	}
	return nil
}

func (db DB) DelTourn(id int) error {
	return nil
}
