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
	logrus.SetLevel(logrus.DebugLevel)
	var err error
	connStr := "user=postgres password=docker2147 dbname=gaming_website host=localhost port=5432 sslmode=disable"
	dbT, err = New(connStr)
	if err != nil {
		logrus.Fatal(err)
	}
	code := m.Run()
	os.Exit(code)
}

/*
func TestMain(m *testing.M) {
	code, err := runTests(m)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}
func runTests(m *testing.M) (int, error) {
	var err error
	dbT, err = New()
	if err != nil {
		return 0, err
	}
	defer func() {
		err = dbT.Close()
		if err != nil {
			logrus.Warn(err)
		}
	}()
	return m.Run(), nil
}
*/

func TestCanRegister(t *testing.T) {
	r := require.New(t)
	tt := []struct {
		name string
		user entity.User
		err  string
	}{
		{
			name: "empty user",
			user: entity.User{},
			err:  "user's name is empty",
		},
		{
			name: "low balance",
			user: entity.User{
				Name:    "Artem",
				Balance: 10,
			},
			err: "user has got not enough points",
		},
		{
			name: "all is ok",
			user: entity.User{
				Name:    "Artem",
				Balance: 300,
			},
			err: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			errStr := ""
			err := tc.user.CanRegister()
			if err != nil {
				errStr = err.Error()
			}
			r.Equal(tc.err, errStr)
		})
	}
}

func TestDB_UserFund(t *testing.T) {
	r := require.New(t)
	u := entity.User{
		Name:    "John",
		Balance: 300,
	}
	us, err := dbT.SaveUser(u)
	r.NoError(err)

	us, err = dbT.UserFund(us.ID, 1)
	r.NoError(err)
	r.Equal(1, us.Balance)
}

func TestDB_UserTake(t *testing.T) {
	r := require.New(t)

	tt := []struct {
		name string
		user entity.User
		err  entity.Error
	}{
		{
			name: "low balance",
			user: entity.User{
				Name:    "Artem",
				Balance: 600,
			},
			err: entity.Error{Type: entity.ErrDB},
		},
		{
			name: "all is ok",
			user: entity.User{
				Name:    "Artem",
				Balance: 1000,
			},
			err: entity.Error{Type: ""},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			us, err := dbT.SaveUser(tc.user)
			r.NoError(err)
			us, err = dbT.UserTake(us.ID, 400)
			errT, ok := err.(entity.Error)
			if !ok {
				errT.Type = ""
			}
			r.Equal(tc.err.Type, errT.Type)
			r.Equal(300, us.Balance)
		})
	}
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
	r.Error(dbT.DeleteUser(777))
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
