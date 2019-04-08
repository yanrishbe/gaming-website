package db

import (
	"testing"
	"time"

	"github.com/yanrishbe/gaming-website/entities"

	"github.com/stretchr/testify/require"
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
	r.NoError(db.SaveUser(&u))

	for i := 0; i < 100; i++ {
		go func() {
			r.NoError(db.UserFund(1, 1))
		}()
		time.Sleep(10 * time.Millisecond)
	}
	r.Equal(100, db.UsersMap[1].Balance)
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
	r.NoError(db.SaveUser(&u))
	r.NoError(db.SaveUser(&y))
	for i := 0; i < 100; i++ {
		go func() {
			r.Error(db.UserTake(2, 400))
		}()
		time.Sleep(10 * time.Millisecond)
	}

	for i := 0; i < 100; i++ {
		go func() {
			r.NoError(db.UserTake(1, 1))
		}()
		time.Sleep(10 * time.Millisecond)
	}
	r.Equal(200, db.UsersMap[1].Balance)
}
func TestDB_DeleteUser(t *testing.T) {
	r := require.New(t)
	db := New()
	r.Error(db.DeleteUser(100))
	u := entities.User{
		Name:    "Jana",
		Balance: 300,
	}
	r.NoError(db.SaveUser(&u))
	r.NoError(db.DeleteUser(1))
}

func TestDB_SaveUser(t *testing.T) {
	r := require.New(t)
	db := New()

	for i := 0; i < 100; i++ {
		go func() {
			u := entities.User{
				Name:    "Jana",
				Balance: 600,
			}
			r.NoError(db.SaveUser(&u))
		}()
		time.Sleep(10 * time.Millisecond)
	}
	r.Equal(100, len(db.UsersMap))

	for i := 1; i < 101; i++ {
		r.Equal(300, db.UsersMap[i].Balance)
	}
	u := entities.User{
		Name:    "",
		Balance: 600,
	}
	r.Error(db.SaveUser(&u))
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
			r.NoError(db.SaveUser(&u))
		}()
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(10 * time.Millisecond)
	r.Equal(100, len(db.UsersMap))

	for i := 1; i < 101; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				r.NoError(db.UserTake(i, 1))
			}
		}(i)
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(10 * time.Millisecond)
	for i := 1; i < 101; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				r.NoError(db.UserFund(i, 1))
			}
		}(i)
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(10 * time.Millisecond)
	for i := 1; i < 101; i++ {
		r.Equal(300, db.UsersMap[i].Balance)
	}
}
