package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/yanrishbe/gaming-website/entity"

	"github.com/stretchr/testify/require"
)

func marshal(t *testing.T, input interface{}) []byte {
	data, err := json.Marshal(input)
	require.NoError(t, err)
	return data
}

func unmarshal(t *testing.T, data []byte, output interface{}) {
	errResponse := json.Unmarshal(data, &output)
	require.NoError(t, error(errResponse))
}

func TestAPI_RegisterNewUser(t *testing.T) {
	r := require.New(t)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	api := New()
	api.InitRouter()
	var userOne = entity.User{Name: "", Balance: 400}
	var userOneByte = marshal(t, userOne)
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userOneByte))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userOneResponse = &UserResp{}
	unmarshal(t, resp.Body.Bytes(), userOneResponse)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)
	var userOneExpected = &UserResp{User: entity.User{Name: "", Balance: 400}, Error: entity.Error{Type: "user's data is not valid"}}
	r.Exactly(userOneExpected.Error.Type, userOneResponse.Error.Type)

	var notUser = "Not a user"
	req = httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer([]byte(notUser)))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)

	userOne = entity.User{Name: "userOne", Balance: 400}
	userOneByte = marshal(t, userOne)
	req = httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userOneByte))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	unmarshal(t, resp.Body.Bytes(), userOneResponse)
	userOneExpected = &UserResp{User: entity.User{ID: 1, Name: "userOne", Balance: 100}, Error: entity.Error{}}
	r.EqualValues(http.StatusOK, resp.Code)
	r.Exactly(userOneExpected, userOneResponse)
}

func TestAPI_GetUser(t *testing.T) {
	r := require.New(t)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	api := New()
	api.InitRouter()
	var userTwo = entity.User{ID: 1, Name: "userTwo", Balance: 100}
	api.DB.UsersMap[1] = userTwo
	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userTwoResponse = &UserResp{}
	unmarshal(t, resp.Body.Bytes(), userTwoResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	var userTwoExpected = &UserResp{User: entity.User{ID: 1, Name: "userTwo", Balance: 100}, Error: entity.Error{}}
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
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	api := New()
	api.InitRouter()
	var userThree = entity.User{ID: 3, Name: "userThree", Balance: 500}
	api.DB.UsersMap[2] = userThree
	req := httptest.NewRequest(http.MethodDelete, "/user/2", nil)
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusNoContent, resp.Code)

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
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	api := New()
	api.InitRouter()
	var userFour = entity.User{ID: 2, Name: "userFour", Balance: 800}
	api.DB.UsersMap[2] = userFour
	var userFourPoints = ReqPoints{Points: 200}
	var userFourPointsBytes = marshal(t, userFourPoints)
	req := httptest.NewRequest(http.MethodPost, "/user/2/take", bytes.NewBuffer(userFourPointsBytes))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userFourResponse = &UserResp{}
	unmarshal(t, resp.Body.Bytes(), userFourResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	var userFourExpected = &UserResp{User: entity.User{ID: 2, Name: "userFour", Balance: 600}, Error: entity.Error{}}
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

	userFourPoints = ReqPoints{Points: 700}
	userFourPointsBytes = marshal(t, userFourPoints)
	req = httptest.NewRequest(http.MethodPost, "/user/2/take", bytes.NewBuffer(userFourPointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	unmarshal(t, resp.Body.Bytes(), userFourResponse)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)
	userFourExpected = &UserResp{Error: entity.Error{Type: "balance is not enough"}}
	r.Exactly(userFourExpected.Error.Type, userFourResponse.Error.Type)

}

func TestAPI_FundUserPoints(t *testing.T) {
	r := require.New(t)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	api := New()
	api.InitRouter()
	var userFive = entity.User{ID: 2, Name: "userFive", Balance: 200}
	api.DB.UsersMap[2] = userFive
	var userFivePoints = ReqPoints{Points: 400}
	var userFivePointsBytes = marshal(t, userFivePoints)
	req := httptest.NewRequest(http.MethodPost, "/user/2/fund", bytes.NewBuffer(userFivePointsBytes))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userFiveResponse = &UserResp{}
	unmarshal(t, resp.Body.Bytes(), userFiveResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	var userFiveExpected = &UserResp{User: entity.User{ID: 2, Name: "userFive", Balance: 600}, Error: entity.Error{}}
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

func TestAPI_DataRace(t *testing.T) {
	r := require.New(t)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	api := New()
	api.InitRouter()
	userDR := entity.User{Name: "User", Balance: 500}
	userDRByte := marshal(t, userDR)
	// req := make ([]*http.Request,100)
	//req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userDRByte))
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer func() {
				//log.Println("worker:", i, "total users", api.DB.CountUsers())
				wg.Done()
			}()
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userDRByte))
			api.Router.ServeHTTP(resp, req)
		}(i)
	}
	wg.Wait()
	r.Equal(100, api.DB.CountUsers())
	var userDRPoints = ReqPoints{Points: 1}
	var userDRBytes = marshal(t, userDRPoints)
	var wg2 sync.WaitGroup
	wg2.Add(100)

	for i := 0; i < 100; i++ {
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/user/1/take", bytes.NewBuffer(userDRBytes))
			//us, _ := api.DB.GetUser(1)
			//log.Println("worker:", i, "points", us.Balance)
			api.Router.ServeHTTP(resp, req)
		}(i, &wg2)
	}
	wg2.Wait()
	_, errGet := api.DB.GetUser(1)
	r.NoError(errGet)
	bal, _ := api.DB.GetBalance(1)
	r.Equal(100, bal)
}
