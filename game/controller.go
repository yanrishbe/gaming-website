package game

import (
	"errors"
	"math/rand"
	"time"

	"github.com/yanrishbe/gaming-website/entity"
	"github.com/yanrishbe/gaming-website/postgres"
)

type Controller struct {
	db postgres.DB
}

func New(db postgres.DB) Controller {
	return Controller{db: db}
}

func (c Controller) RegUser(u entity.User) (entity.User, error) {
	err := u.IsValid()
	if err != nil {
		return u, err
	}
	if u.Balance < 300 {
		return u, entity.RegErr(errors.New("low balance"))
	}
	u.Balance -= 300
	return c.db.CreateUser(u)
}

func (c Controller) GetUser(id int) (entity.User, error) {
	return c.db.GetUser(id)
}

func (c Controller) DelUser(id int) error {
	return c.db.DelUser(id)
}

func (c Controller) TakePoints(id, points int) (entity.User, error) {
	if points <= 0 {
		return entity.User{}, entity.PointsErr(errors.New("points must be greater than 0"))
	}
	return c.db.TakePoints(id, points)
}

func (c Controller) FundPoints(id, points int) (entity.User, error) {
	if points <= 0 {
		return entity.User{}, entity.PointsErr(errors.New("points must be greater than 0"))
	}
	return c.db.FundPoints(id, points)
}

func (c Controller) RegTourn(t entity.Tournament) (entity.Tournament, error) {
	err := t.IsValid()
	if err != nil {
		return t, err
	}
	if t.Deposit <= 0 {
		return t, entity.RegErr(errors.New("deposit must be greater than 0"))
	}
	return c.db.CreateTourn(t)
}

func (c Controller) GetTourn(id int) (entity.Tournament, error) {
	t, err := c.db.GetTourn(id)
	return t, err
}

func (c Controller) JoinTourn(tID, uID int) (entity.Tournament, error) {
	t, err := c.db.JoinTourn(tID, uID, func(a int, b int) error {
		if a < b {
			return entity.RegErr(errors.New("balance is lower than deposit"))
		}
		return nil
	})
	if err != nil {
		return t, err
	}
	return c.db.GetTourn(t.ID)
}

func (c Controller) FinishTourn(id int) (entity.Tournament, error) {
	err := c.db.FinishTourn(id, func(users []int) int {
		rand.Seed(time.Now().UTC().UnixNano())
		r := rand.Intn(len(users))
		return users[r]
	})
	if err != nil {
		return entity.Tournament{}, err
	}
	return c.db.GetTourn(id)
}

func (c Controller) DelTourn(id int) error {
	t, err := c.GetTourn(id)
	if err != nil {
		return err
	}
	if t.Status != entity.Finished {
		_, err := c.FinishTourn(id)
		if err != nil {
			return err
		}
	}
	err = c.db.DelTourn(id)
	return err
}
