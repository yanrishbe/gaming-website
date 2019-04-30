package game

import (
	"errors"

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
	if points < 1 {
		return entity.User{}, entity.PointsErr(errors.New("points must be > 0"))
	}
	return c.db.TakePoints(id, points)
}

func (c Controller) FundPoints(id, points int) (entity.User, error) {
	if points < 1 {
		return entity.User{}, entity.PointsErr(errors.New("points must be > 0"))
	}
	return c.db.FundPoints(id, points)
}
