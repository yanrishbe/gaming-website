package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/yanrishbe/gaming-website/db"

	"github.com/sirupsen/logrus"

	"github.com/yanrishbe/gaming-website/entity"

	"github.com/stretchr/testify/require"
)

var api *API

func TestMain(m *testing.M) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	var err error
	connStr := "user=postgres password=docker2147 dbname=gaming_website host=localhost port=5432 sslmode=disable"
	dbT, err := db.New(connStr)
	if err != nil {
		logrus.Fatal(err)
	}
	api, err = New(dbT)
	if err != nil {
		logrus.Fatal(err)
	}
	code := m.Run()
	err = api.DB.Close()
	if err != nil {
		logrus.Fatal(err)
	}
	os.Exit(code)
}

func marshal(t *testing.T, input interface{}) []byte {
	data, err := json.Marshal(input)
	require.NoError(t, err)
	return data
}

func unmarshal(t *testing.T, data []byte, output interface{}) {
	errResponse := json.Unmarshal(data, &output)
	require.NoError(t, error(errResponse))
}

func doRequest(t *testing.T, method, url string, body interface{}) (entity.User, error) {
	userOneByte := marshal(t, body)
	req := httptest.NewRequest(method, url, bytes.NewBuffer(userOneByte))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	if resp.Code == 204 {
		return entity.User{}, nil
	} else if resp.Result().StatusCode >= 200 && resp.Result().StatusCode < 300 {
		u := entity.User{}
		unmarshal(t, resp.Body.Bytes(), &u)
		return u, nil
	} else {
		err := entity.Error{}
		unmarshal(t, resp.Body.Bytes(), &err)
		err.Message = ""
		return entity.User{}, err
	}
}

func TestAPI_RegisterNewUser(t *testing.T) {
	r := require.New(t)
	_, err := doRequest(t, http.MethodPost, "/user", entity.User{Name: "", Balance: 400})
	r.Exactly(entity.Error{Type: entity.ErrCannotReg, Code: 422}, err)
	_, err = doRequest(t, http.MethodPost, "/user", "string")
	r.Exactly(entity.Error{Type: entity.ErrDecode, Code: 422}, err)
	u, err := doRequest(t, http.MethodPost, "/user", entity.User{Name: "userOne", Balance: 400})
	r.NoError(err)
	r.NotEmpty(u.ID)
	r.Equal(100, u.Balance)
}

func TestAPI_GetUser(t *testing.T) {
	r := require.New(t)
	u, err := doRequest(t, http.MethodPost, "/user", entity.User{Name: "getUser", Balance: 400})
	r.NoError(err)
	id := strconv.Itoa(u.ID)
	u, err = doRequest(t, http.MethodGet, "/user/"+id, nil)
	r.Exactly(entity.User{ID: u.ID, Name: "getUser", Balance: 100}, u)
	_, err = doRequest(t, http.MethodGet, "/user/str", nil)
	r.Exactly(entity.Error{Type: entity.ErrID, Code: 400}, err)
	_, err = doRequest(t, http.MethodGet, "/user/777", nil)
	r.Exactly(entity.Error{Type: entity.ErrUserNotFound, Code: 404}, err)

}

func TestAPI_DeleteUser(t *testing.T) {
	r := require.New(t)
	u, err := doRequest(t, http.MethodPost, "/user", entity.User{Name: "deleteUser", Balance: 400})
	r.NoError(err)
	id := strconv.Itoa(u.ID)
	_, err = doRequest(t, http.MethodDelete, "/user/"+id, nil)
	r.Nil(err)
	_, err = doRequest(t, http.MethodDelete, "/user/str", nil)
	r.Exactly(entity.Error{Type: entity.ErrID, Code: 400}, err)
	_, err = doRequest(t, http.MethodDelete, "/user/777", nil)
	r.Exactly(entity.Error{Type: entity.ErrUserNotFound, Code: 404}, err)
}

func TestAPI_TakeUserPoints(t *testing.T) {
	r := require.New(t)
	u, err := doRequest(t, http.MethodPost, "/user", entity.User{Name: "takePoints", Balance: 800})
	r.NoError(err)
	id := strconv.Itoa(u.ID)
	u, err = doRequest(t, http.MethodPost, "/user/"+id+"/take", ReqPoints{Points: 200})
	r.Exactly(300, u.Balance)
	_, err = doRequest(t, http.MethodPost, "/user/str/take", ReqPoints{Points: 200})
	r.Exactly(entity.Error{Type: entity.ErrID, Code: 400}, err)
	_, err = doRequest(t, http.MethodPost, "/user/"+id+"/take", ReqPoints{Points: 400})
	r.Exactly(entity.Error{Type: entity.ErrDB, Code: 503}, err)
}

func TestAPI_FundUserPoints(t *testing.T) {
	r := require.New(t)
	u, err := doRequest(t, http.MethodPost, "/user", entity.User{Name: "fundPoints", Balance: 800})
	r.NoError(err)
	id := strconv.Itoa(u.ID)
	u, err = doRequest(t, http.MethodPost, "/user/"+id+"/fund", ReqPoints{Points: 200})
	r.Exactly(700, u.Balance)
	_, err = doRequest(t, http.MethodPost, "/user/str/fund", ReqPoints{Points: 200})
	r.Exactly(entity.Error{Type: entity.ErrID, Code: 400}, err)
}
