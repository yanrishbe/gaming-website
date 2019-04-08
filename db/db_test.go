package db

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yanrishbe/gaming-website/entities"
)

func TestCanRegister(t *testing.T) {
	r := require.New(t)
	var user = []entities.User{
		{Name: ""},
		{Name: "Y", Balance: 300},
		{Name: "N", Balance: 0},
	}
	r.False(canRegister(user[0]))
	r.True(canRegister(user[1]))
	r.False(canRegister(user[2]))
}

func TestNew(t *testing.T) {
	require.NotEmpty(t, New())
}

func TestDB_UserFund(t *testing.T) {
	r := require.New(t)
	db := New()
	r.Error(db.UserFund(100, 1))
	u := entities.User{
		Name:    "Jana",
		Balance: 300,
	}
	_, errSave := db.SaveUser(&u)
	r.NoError(errSave)
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(wg *sync.WaitGroup) {
			r.NoError(db.UserFund(1, 1))
			wg.Done()
		}(&wg)
	}
	wg.Wait()
	user, errGet := db.GetUser(1)
	r.NoError(errGet)
	r.Equal(100, user.Balance)
}

func TestDB_UserTake(t *testing.T) {
	r := require.New(t)
	db := New()
	r.Error(db.UserTake(100, 1))

	u, y := entities.User{
		Name:    "Jana",
		Balance: 600,
	}, entities.User{
		Name:    "M",
		Balance: 400,
	}
	_, errSave := db.SaveUser(&u)
	r.NoError(errSave)
	_, errSave = db.SaveUser(&y)
	r.NoError(errSave)
	var wg sync.WaitGroup
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(wg *sync.WaitGroup) {
			r.Error(db.UserTake(2, 400))
			wg.Done()
		}(&wg)
	}

	for i := 0; i < 100; i++ {
		go func(wg *sync.WaitGroup) {
			r.NoError(db.UserTake(1, 1))
			wg.Done()
		}(&wg)

	}
	wg.Wait()
	user, errGet := db.GetUser(1)
	r.NoError(errGet)
	r.Equal(200, user.Balance)
	user, errGet = db.GetUser(2)
	r.NoError(errGet)
	r.Equal(100, user.Balance)
}

func TestDB_DeleteUser(t *testing.T) {
	r := require.New(t)
	db := New()
	r.Error(db.DeleteUser(100))
	u := entities.User{
		Name:    "Jana",
		Balance: 300,
	}
	_, errSave := db.SaveUser(&u)
	r.NoError(errSave)
	r.NoError(db.DeleteUser(1))
}

func TestDB_SaveUser(t *testing.T) {
	r := require.New(t)
	db := New()
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(wg *sync.WaitGroup) {
			u := entities.User{
				Name:    "Jana",
				Balance: 600,
			}
			_, errSave := db.SaveUser(&u)
			r.NoError(errSave)
			wg.Done()
		}(&wg)
	}
	wg.Wait()
	r.Equal(100, db.CountUsers())

	for i := 1; i < 101; i++ {
		user, errGet := db.GetUser(i)
		r.NoError(errGet)
		r.Equal(300, user.Balance)
	}
	u := entities.User{
		Name:    "",
		Balance: 600,
	}
	_, errSave := db.SaveUser(&u)
	r.Error(errSave)
}

func TestDB_DataRace(t *testing.T) {
	r := require.New(t)
	db := New()
	for i := 0; i < 100; i++ {
		go func() {
			u := entities.User{
				Name:    "Jana",
				Balance: 600,
			}
			_, errSave := db.SaveUser(&u)
			r.NoError(errSave)
		}()
	}
	for i := 1; i < 101; i++ {
		go func() {
			r.NoError(db.UserTake(1, 1))
		}()
	}

	for i := 1; i < 101; i++ {
		go func() {
			r.NoError(db.UserFund(1, 2))
		}()
	}
	r.Equal(100, db.CountUsers())
	user, errGet := db.GetUser(1)
	r.NoError(errGet)
	r.Equal(400, user.Balance)
}
