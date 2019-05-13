package postgres

import (
	"database/sql"
	"errors"
	"fmt"

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

func (db DB) GetTourn(id int) (entity.Tournament, error) {
	if id <= 0 {
		return entity.Tournament{}, entity.InvIDErr(errors.New("expected id greater than 0"))
	}
	var t entity.Tournament
	var finished bool

	err := db.db.QueryRow(`
		SELECT finished
		FROM tournaments
		WHERE id=$1`, id).Scan(&finished)
	if err == sql.ErrNoRows {
		return entity.Tournament{}, entity.ReqErr(fmt.Errorf("tournament doesn't exist: %v", err))
	} else if err != nil {
		return entity.Tournament{}, entity.DBErr(err)
	}

	if finished {
		t.Status = entity.Finished

		err = db.db.QueryRow(`
		SELECT id, name, deposit, prize, winner_id
		FROM tournaments 
		WHERE id = $1`,
			id).Scan(&t.ID, &t.Name, &t.Deposit, &t.Prize, &t.Winner)
		if err != nil {
			return t, entity.DBErr(fmt.Errorf("can't get tournament: %v", err))
		}
	} else {
		t.Status = entity.Active
		err = db.db.QueryRow(`
		SELECT id, name, deposit, prize
		FROM tournaments 
		WHERE id = $1`,
			id).Scan(&t.ID, &t.Name, &t.Deposit, &t.Prize)
		if err != nil {
			return t, entity.DBErr(fmt.Errorf("can't get tournament: %v", err))
		}
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
		if err != nil {
			return t, entity.DBErr(fmt.Errorf("can't get tournament data: %v", err))
		}
		if w.ID == t.Winner {
			w.Winner = true
		}
		t.Users = append(t.Users, w)
	}
	if len(t.Users) == 0 {
		t.Users = []entity.Winner{}
	}
	err = rows.Err()
	if err != nil {
		return t, entity.DBErr(err)
	}
	return t, nil
}

func (db DB) JoinTourn(tID, uID int, check func(balance int, deposit int) error) (entity.Tournament, error) {
	var t entity.Tournament
	t.ID = tID

	tx, err := db.db.Begin()
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("transaction error: %v", err))
	}
	defer tx.Rollback()

	var finished bool
	err = tx.QueryRow(`
		SELECT finished
		FROM tournaments
		WHERE id=$1`, tID).Scan(&finished)
	if err == sql.ErrNoRows {
		return t, entity.ReqErr(fmt.Errorf("tournament doesn't exist: %v", err))
	} else if err != nil {
		return t, entity.DBErr(err)
	}
	if finished {
		return t, entity.ReqErr(errors.New("the tournament is finished"))
	}

	var name string
	err = tx.QueryRow(`
		SELECT name
		FROM users
		WHERE id = $1`, uID).Scan(&name)
	if err == sql.ErrNoRows {
		return t, entity.ReqErr(errors.New("the user doesn't exist"))
	} else if err != nil {
		return t, entity.DBErr(err)
	}

	var id int
	err = tx.QueryRow(`
		SELECT tournament_id
		FROM tournament_req 
		WHERE user_id = $1 AND tournament_id = $2`, uID, tID).Scan(&id)
	if err == nil {
		return t, entity.RegErr(fmt.Errorf("user is already registered"))
	} else if id != 0 {
		return t, entity.RegErr(fmt.Errorf("user is already registered"))
	}

	err = tx.QueryRow(`
		SELECT deposit 
		FROM tournaments 
		WHERE id = $1`, tID).Scan(&t.Deposit)
	if err != nil {
		return t, entity.DBErr(err)
	}
	var balance int
	err = tx.QueryRow(`
		SELECT balance 
		FROM users 
		WHERE id = $1`, uID).Scan(&balance)
	if err != nil {
		return t, entity.DBErr(err)
	}

	err = check(balance, t.Deposit)
	if err != nil {
		return t, err
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

func getTournUsers(tx *sql.Tx, tID int) ([]int, error) {
	rows, err := tx.Query(`
		SELECT user_id
		FROM tournament_req
		WHERE tournament_id = $1`, tID)
	if err != nil {
		return nil, entity.DBErr(fmt.Errorf("can't get data: %v", err))
	}
	defer rows.Close()
	var userIDs []int
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, entity.DBErr(fmt.Errorf("can't get data: %v", err))
		}
		userIDs = append(userIDs, id)
	}
	err = rows.Err()
	if err != nil {
		return nil, entity.DBErr(fmt.Errorf("rows error: %v", err))
	}
	return userIDs, nil
}

func (db DB) FinishTourn(tID int, chooseWinner func(ids []int) int) error {
	tx, err := db.db.Begin()
	if err != nil {
		return entity.DBErr(fmt.Errorf("transaction error: %v", err))
	}
	defer tx.Rollback()

	var finished bool
	err = tx.QueryRow(`
		SELECT finished
		FROM tournaments
		WHERE id=$1`, tID).Scan(&finished)
	if err == sql.ErrNoRows {
		return entity.ReqErr(fmt.Errorf("tournament doesn't exist: %v", err))
	} else if err != nil {
		return entity.DBErr(err)
	}

	rows, err := tx.Query(`
		SELECT user_id
		FROM tournament_req
		WHERE tournament_id = $1`, tID)
	if err != nil {
		return entity.DBErr(fmt.Errorf("can't get tournament data: %v", err))
	}
	if !rows.Next() {
		return entity.ReqErr(errors.New("can't finish, no users"))
	}
	err = rows.Err()
	if err != nil {
		return entity.DBErr(err)
	}
	rows.Close()
	//todo trying to start a new statement before reading all of the rows of the preceding statement
	users, err := getTournUsers(tx, tID)
	if err != nil {
		return err
	}
	var uID = chooseWinner(users)

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
	tx, err := db.db.Begin()
	if err != nil {
		return entity.DBErr(fmt.Errorf("transaction error: %v", err))
	}
	defer tx.Rollback()
	_, err = tx.Exec(`
		DELETE FROM tournament_req
		WHERE tournament_id = $1`, id)
	if err != nil {
		return entity.DBErr(fmt.Errorf("can't delete from tournament_req table: %v", err))
	}
	_, err = tx.Exec(`
		DELETE FROM tournaments
		WHERE id = $1`, id)
	if err != nil {
		return entity.DBErr(fmt.Errorf("can't delete from tournaments table: %v", err))
	}
	err = tx.Commit()
	if err != nil {
		return entity.DBErr(fmt.Errorf("transaction error: %v", err))
	}
	return nil
}
