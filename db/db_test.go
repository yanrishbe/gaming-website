package db

import (
	"os"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/require"
	"github.com/yanrishbe/gaming-website/entity"
)

var dbT DB

func TestMain(m *testing.M) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	var err error
	dbT, err = New()
	dbT.createTables() //fixme require???
	if err != nil {
		logrus.Fatal(err)
	}
	code := m.Run()
	dbT.Close()
	os.Exit(code)
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

func TestDB_UserFund(t *testing.T) {
	r := require.New(t)
	u := entity.User{
		Name:    "John",
		Balance: 300,
	}
	us, err := dbT.SaveUser(u)
	r.NoError(err)
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			_, err := dbT.UserFund(us.ID, 1)
			r.NoError(err)
			wg.Done()
		}()
	}
	wg.Wait()
	usg, err := dbT.GetUser(us.ID)
	r.NoError(err)
	r.Equal(100, usg.Balance)
}

func TestDB_UserTake(t *testing.T) {
	r := require.New(t)

	u, y := entity.User{
		Name:    "Jana",
		Balance: 600,
	}, entity.User{
		Name:    "M",
		Balance: 400,
	}
	us, err := dbT.SaveUser(u)
	r.NoError(err)
	usE, err := dbT.SaveUser(y)
	r.NoError(err)
	var wg sync.WaitGroup
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func() {
			_, err := dbT.UserTake(usE.ID, 400)
			r.Error(err)
			wg.Done()
		}()
	}

	for i := 0; i < 100; i++ {
		go func() {
			_, err := dbT.UserTake(us.ID, 1)
			r.NoError(err)
			wg.Done()
		}()

	}
	wg.Wait()
	usg, err := dbT.GetUser(us.ID)
	r.NoError(err)
	r.Equal(200, usg.Balance)
	usy, err := dbT.GetUser(usE.ID)
	r.NoError(err)
	r.Equal(100, usy.Balance)
}

func TestDB_DeleteUser(t *testing.T) {
	r := require.New(t)
	u := entity.User{
		Name:    "Jana",
		Balance: 300,
	}
	us, err := dbT.SaveUser(u)
	r.NoError(err)
	r.NoError(dbT.DeleteUser(us.ID))
}

func TestDB_DataRace(t *testing.T) {
	r := require.New(t)

	u := entity.User{
		Name:    "Jana",
		Balance: 600,
	}
	us, err := dbT.SaveUser(u)
	r.NoError(err)

	var wg2 sync.WaitGroup
	wg2.Add(200)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg2.Done()
			_, err := dbT.UserTake(us.ID, 1)
			r.NoError(err)
		}()

	}
	for i := 0; i < 100; i++ {
		go func() {
			defer wg2.Done()
			_, err := dbT.UserFund(us.ID, 2)
			r.NoError(err)
		}()
	}
	wg2.Wait()
	usg, err := dbT.GetUser(us.ID)
	r.NoError(err)
	r.Equal(400, usg.Balance)
}
