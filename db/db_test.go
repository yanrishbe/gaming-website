package db

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

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
	var user = []entity.User{
		{Name: ""},
		{Name: "Y", Balance: 300},
		{Name: "N", Balance: 0},
	}
	r.Error(user[0].CanRegister())
	r.NoError(user[1].CanRegister())
	r.Error(user[2].CanRegister())

	// Ideally in tests entities and expectations should be as near as possible, so it's easy
	// to track logic and to identify which exact test is failing.
	// you also want to check that error you had is the error you expected
	// see example below:
	// (try to break something and notice nice test error message. It will show where exactly test has failed using test name)

	// this is called a "table test"
	tt := []struct { // first we create struct with test cases
		name string
		user entity.User
		err  string
	}{
		{
			name: "empty user", // then we specify test cases with names, explaining what is the case we are testing
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

	for _, tc := range tt { // then we run theese test in this boilerplate loop to make sure error messages will be beautiful)
		t.Run(tc.name, func(t *testing.T) {
			errStr := "" // because your errors are wrapped - it's hard to check equality for them, so we simply check equality of error message stings
			err := tc.user.CanRegister()
			if err != nil {
				errStr = err.Error()
			}
			assert.Equal(t, tc.err, errStr)
		})
	}
}

// Very nice test, nothing to remove or add from here :)
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
	// here you could simply used  "table test" approach i've shown above to avoid these double users
	// and i'm not sure why you need 2 users that both pass the test successfully.
	// usually you don't want to repeat tests, you want each one to be individual from others.

	u, y := entity.User{ // never use this scary multiple var initialization. It looks bad. U should initialize each variable on it's own.
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

// here you check only successful deletion, what about deletion of user that doesn't exist?
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

// This looks like a duplicate from previous test, you were also testing data race there
// So you want to make those test simpler (without go func()...), as data race is already tested here.
// Test itself is nicely written :)
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
