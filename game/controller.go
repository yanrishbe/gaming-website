package game

import (
	"github.com/yanrishbe/gaming-website/entity"
	"github.com/yanrishbe/gaming-website/postgres"
)

type Controller struct {
	db postgres.DB
}

func New(db postgres.DB) Controller {
	return Controller{db: db}
}

func (c Controller) Register(u entity.User) (entity.User, error) {
	id, err := c.db.RegUser(u)
	if err != nil {
		return entity.User{}, err
	}
	u, err = c.db.GetUser(id)
	return u, err
}

func (c Controller) GetUser(id int) (entity.User, error) {

}
