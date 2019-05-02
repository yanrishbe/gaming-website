package postgres

import "github.com/yanrishbe/gaming-website/entity"

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

}
