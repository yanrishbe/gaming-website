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
	return c.db.RegUser(u)
}

func (c Controller) GetUser(id int) (entity.User, error) {
	return c.db.GetUser(id)
}

//func (c Controller) Take
