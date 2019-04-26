package entity

import "errors"

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

func (u User) IsValid() error {
	if u.Name == "" {
		return errors.New("empty name")
	} else if u.Balance < 300 {
		return errors.New("low balance")
	}
	return nil
}
