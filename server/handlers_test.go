package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yanrishbe/gaming-website/entities"
)

//func TestDevices(t *testing.T) {
//	neoHandler, err := neo4j.New(user, pass, host, port)
//	require.NoError(t, err)
//	logger := logrus.New()
//	logger.SetLevel(logrus.DebugLevel)
//	dash, err := New(logger, neoHandler)
//	require.NoError(t, err)
//
//	req := httptest.NewRequest(http.MethodGet, "/count/devices", nil)
//	resp := httptest.NewRecorder()
//	dash.Router.ServeHTTP(resp, req)
//	var devices domain.CountResponse
//	err = json.Unmarshal(resp.Body.Bytes(), &devices)
//	require.NoError(t, err)
//	require.True(t, devices.Count > 0)
//}

func TestCanRegister(t *testing.T) {
	require := require.New(t)
	var user = []entities.User{
		{Name: ""},
		{Name: "Y", Balance: 300},
		{Name: "N", Balance: 0},
	}
	require.False(canRegister(user[0]))
	require.True(canRegister(user[1]))
	require.False(canRegister(user[2]))
}

func TestAPI_RegisterNewUser(t *testing.T) {
	require := require.New(t)
	registerAUser := []byte(`
{
"name": "",
"balance": 400
}`)
	api := New()
	api.InitRouter()
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(registerAUser))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)

	var userResponse = &UserResponse{}
	data := resp.Body.Bytes()
	log.Println(string(data))
	errResponse := json.Unmarshal(resp.Body.Bytes(), userResponse)
	require.NoError(errResponse)
	//require.Equal(100, userResponse.Balance)
	//require.Equ(userResponse.Error)
}

func marshall(t *testing.T, input interface{}) []byte {
	data, err := json.Marshal(input)
	require.NoError(t, err)
	return data
}
