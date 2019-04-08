package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yanrishbe/gaming-website/entities"
)

func marshal(t *testing.T, input interface{}) []byte {
	data, err := json.Marshal(input)
	require.NoError(t, err)
	return data
}

func unmarshal(t *testing.T, data []byte, output interface{}) {
	errResponse := json.Unmarshal(data, &output)
	require.NoError(t, errResponse)
}

func TestAPI_RegisterNewUser(t *testing.T) {
	r := require.New(t)
	api := New()
	api.InitRouter()
	var userOne = entities.User{Name: "", Balance: 400}
	var userOneByte = marshal(t, userOne)
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userOneByte))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userOneResponse = &UserResponse{}
	unmarshal(t, resp.Body.Bytes(), userOneResponse)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)
	var userOneExpected = &UserResponse{User: entities.User{ID: 0, Name: "", Balance: 400}, Error: "user's data is not valid"}
	r.Exactly(userOneExpected, userOneResponse)

	var notUser = "Not a user"
	req = httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer([]byte(notUser)))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)

	userOne = entities.User{Name: "userOne", Balance: 400}
	userOneByte = marshal(t, userOne)
	req = httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userOneByte))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	unmarshal(t, resp.Body.Bytes(), userOneResponse)
	userOneExpected = &UserResponse{User: entities.User{ID: 1, Name: "userOne", Balance: 100}, Error: ""}
	r.EqualValues(http.StatusCreated, resp.Code)
	r.Exactly(userOneExpected, userOneResponse)
}

func TestAPI_GetUser(t *testing.T) {
	r := require.New(t)
	api := New()
	api.InitRouter()
	var userTwo = &entities.User{ID: 1, Name: "userTwo", Balance: 100}
	api.DB.UsersMap[1] = userTwo
	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userTwoResponse = &UserResponse{}
	unmarshal(t, resp.Body.Bytes(), userTwoResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	var userTwoExpected = &UserResponse{User: entities.User{ID: 1, Name: "userTwo", Balance: 100}, Error: ""}
	r.Exactly(userTwoExpected, userTwoResponse)

	req = httptest.NewRequest(http.MethodGet, "/user/str", nil)
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusBadRequest, resp.Code)

	req = httptest.NewRequest(http.MethodGet, "/user/2", nil)
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusNotFound, resp.Code)
}

func TestAPI_DeleteUser(t *testing.T) {
	r := require.New(t)
	api := New()
	api.InitRouter()
	var userThree = &entities.User{ID: 2, Name: "userThree", Balance: 500}
	api.DB.UsersMap[2] = userThree
	req := httptest.NewRequest(http.MethodDelete, "/user/2", nil)
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userThreeResponseDeleted = resp.Body.String()
	r.EqualValues(http.StatusOK, resp.Code)
	r.Equal("successfully deleted the user", userThreeResponseDeleted)

	req = httptest.NewRequest(http.MethodDelete, "/user/str", nil)
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusBadRequest, resp.Code)

	req = httptest.NewRequest(http.MethodDelete, "/user/3", nil)
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusNotFound, resp.Code)
}

func TestAPI_TakeUserPoints(t *testing.T) {
	r := require.New(t)
	api := New()
	api.InitRouter()
	var userFour = &entities.User{ID: 2, Name: "userFour", Balance: 800}
	api.DB.UsersMap[2] = userFour
	var userFourPoints = RequestPoints{Points: 200}
	var userFourPointsBytes = marshal(t, userFourPoints)
	req := httptest.NewRequest(http.MethodPost, "/user/2/take", bytes.NewBuffer(userFourPointsBytes))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userFourResponse = &UserResponse{}
	unmarshal(t, resp.Body.Bytes(), userFourResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	var userFourExpected = &UserResponse{User: entities.User{ID: 2, Name: "userFour", Balance: 600}, Error: ""}
	r.Exactly(userFourExpected, userFourResponse)

	req = httptest.NewRequest(http.MethodPost, "/user/str/take", bytes.NewBuffer(userFourPointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusBadRequest, resp.Code)

	req = httptest.NewRequest(http.MethodPost, "/user/3/take", bytes.NewBuffer(userFourPointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusNotFound, resp.Code)

	var notPoints = "Not Points"
	req = httptest.NewRequest(http.MethodPost, "/user/2/take", bytes.NewBuffer([]byte(notPoints)))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)

	userFourPoints = RequestPoints{Points: 700}
	userFourPointsBytes = marshal(t, userFourPoints)
	req = httptest.NewRequest(http.MethodPost, "/user/2/take", bytes.NewBuffer(userFourPointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	unmarshal(t, resp.Body.Bytes(), userFourResponse)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)
	userFourExpected = &UserResponse{User: entities.User{ID: 0, Name: "", Balance: 0},
		Error: "not enough balance to execute the request"}
	r.Exactly(userFourExpected, userFourResponse)

}

func TestAPI_FundUserPoints(t *testing.T) {
	r := require.New(t)
	api := New()
	api.InitRouter()
	var userFive = &entities.User{ID: 2, Name: "userFive", Balance: 200}
	api.DB.UsersMap[2] = userFive
	var userFivePoints = RequestPoints{Points: 400}
	var userFivePointsBytes = marshal(t, userFivePoints)
	req := httptest.NewRequest(http.MethodPost, "/user/2/fund", bytes.NewBuffer(userFivePointsBytes))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userFiveResponse = &UserResponse{}
	unmarshal(t, resp.Body.Bytes(), userFiveResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	var userFiveExpected = &UserResponse{User: entities.User{ID: 2, Name: "userFive", Balance: 600}, Error: ""}
	r.Exactly(userFiveExpected, userFiveResponse)

	req = httptest.NewRequest(http.MethodPost, "/user/str/fund", bytes.NewBuffer(userFivePointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusBadRequest, resp.Code)

	req = httptest.NewRequest(http.MethodPost, "/user/3/fund", bytes.NewBuffer(userFivePointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusNotFound, resp.Code)

	var notPoints = "Not Points"
	req = httptest.NewRequest(http.MethodPost, "/user/2/fund", bytes.NewBuffer([]byte(notPoints)))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)

}
