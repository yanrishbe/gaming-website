package db

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yanrishbe/gaming-website/entity"
)

func TestMain(m *testing.M) {

}

func TestCanRegister(t *testing.T) {
	r := require.New(t)
	var user = []entity.User{
		{Name: ""},
		{Name: "Y", Balance: 300},
		{Name: "N", Balance: 0},
	}
	r.Error(user[0].CanRegister())
	r.NoError(user[1].CanRegister())
	r.Error(user[2].CanRegister())
}

func TestNew(t *testing.T) {
	require.NotEmpty(t, New())
}

func TestDB_UserFund(t *testing.T) {
	r := require.New(t)
	db := New()
	_, err := db.UserFund(100, 1)
	r.Error(err)
	u := entity.User{
		Name:    "Jana",
		Balance: 300,
	}
	_, err = db.SaveUser(u)
	r.NoError(err)
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			_, err := db.UserFund(1, 1)
			r.NoError(err)
			wg.Done()
		}()
	}
	wg.Wait()
	us, err := db.GetUser(1)
	r.NoError(err)
	r.Equal(100, us.Balance)
}

func TestDB_UserTake(t *testing.T) {
	r := require.New(t)
	db := New()
	_, err := db.UserTake(100, 1)
	r.Error(err)

	u, y := entity.User{
		Name:    "Jana",
		Balance: 600,
	}, entity.User{
		Name:    "M",
		Balance: 400,
	}
	_, err = db.SaveUser(u)
	r.NoError(err)
	_, err = db.SaveUser(y)
	r.NoError(err)
	var wg sync.WaitGroup
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func() {
			_, err := db.UserTake(2, 400)
			r.Error(err)
			wg.Done()
		}()
	}

	for i := 0; i < 100; i++ {
		go func() {
			_, err := db.UserTake(1, 1)
			r.NoError(err)
			wg.Done()
		}()

	}
	wg.Wait()
	us, err := db.GetUser(1)
	r.NoError(err)
	r.Equal(200, us.Balance)
	us, err = db.GetUser(2)
	r.NoError(err)
	r.Equal(100, us.Balance)
}

func TestDB_DeleteUser(t *testing.T) {
	r := require.New(t)
	db := New()
	r.Error(db.DeleteUser(100))
	u := entity.User{
		Name:    "Jana",
		Balance: 300,
	}
	_, err := db.SaveUser(u)
	r.NoError(err)
	r.NoError(db.DeleteUser(1))
}

func TestDB_SaveUser(t *testing.T) {
	r := require.New(t)
	db := New()
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			u := entity.User{
				Name:    "Jana",
				Balance: 600,
			}
			_, err := db.SaveUser(u)
			r.NoError(err)
			wg.Done()
		}()
	}
	wg.Wait()
	r.Equal(100, db.CountUsers())

	for i := 1; i < 101; i++ {
		us, err := db.GetUser(i)
		r.NoError(err)
		r.Equal(300, us.Balance)
	}
	u := entity.User{
		Name:    "",
		Balance: 600,
	}
	_, err := db.SaveUser(u)
	r.Error(err)
}

func TestDB_DataRace(t *testing.T) {
	r := require.New(t)
	db := New()
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {

		go func() {
			defer wg.Done()
			u := entity.User{
				Name:    "Jana",
				Balance: 600,
			}
			_, err := db.SaveUser(u)
			r.NoError(err)
		}()

	}
	wg.Wait()
	var wg2 sync.WaitGroup
	wg2.Add(200)
	for i := 1; i < 101; i++ {
		go func() {
			defer wg2.Done()
			_, err := db.UserTake(1, 1)
			r.NoError(err)
		}()

	}
	for i := 1; i < 101; i++ {
		go func() {
			defer wg2.Done()
			_, err := db.UserFund(1, 2)
			r.NoError(err)
		}()
	}
	wg2.Wait()
	r.Equal(100, db.CountUsers())
	us, err := db.GetUser(1)
	r.NoError(err)
	r.Equal(400, us.Balance)
}
