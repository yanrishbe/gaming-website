package postgres

import (
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
		return t, entity.DBErr(err)
	}
	return t, nil
}

func (db DB) GetTourn(id int) (entity.Tournament, error) {
	if id <= 0 {
		return entity.Tournament{}, entity.InvIDErr(errors.New("expected id > 0"))
	}
	tx, err := db.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return entity.Tournament{}, entity.DBErr(err)
	}
	var t entity.Tournament
	err = tx.QueryRow(`
		SELECT id, name, deposit, prize 
		FROM tournaments 
		WHERE id = $1 AND finished = FALSE`,
		id).Scan(&t.ID, &t.Name, &t.Deposit, &t.Prize)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't get data: %v", err))
	}
	uR, err := tx.Query(`
		SELECT users.id, users.name
		FROM tournament_req
		INNER JOIN users ON tournament_req.user_id = users.id
		WHERE tournament_req.tournament_id = $1`, id)
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("can't get data: %v", err))
	}
	//close rows or not????????
	defer uR.Close()
	var u entity.User
	for uR.Next() {
		err := uR.Scan(&u.ID, &u.Name)
		if err != nil {
			return t, entity.DBErr(fmt.Errorf("can't get data: %v", err))
		}
		t.Users = append(t.Users, u)
	}
	err = uR.Err()
	if err != nil {
		return t, entity.DBErr(fmt.Errorf("rows error: %v", err))
	}
	err = tx.Commit()
	if err != nil {
		return t, entity.DBErr(err)
	}
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
