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
	userOne := entity.User{Name: "", Balance: 400}
	userOneByte := marshal(t, userOne)
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userOneByte))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	response := entity.Error{}
	unmarshal(t, resp.Body.Bytes(), &response)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)
	expected := entity.Error{Type: "user's data is not valid"}
	r.Exactly(expected.Type, response.Type)

	notUser := "Not a user"
	req = httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer([]byte(notUser)))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)

	userOne = entity.User{Name: "userOne", Balance: 400}
	userOneByte = marshal(t, userOne)
	req = httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userOneByte))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	userOneResponse := entity.User{}
	unmarshal(t, resp.Body.Bytes(), &userOneResponse)
	userOneExpected := entity.User{ID: 1, Name: "userOne", Balance: 100}
	r.EqualValues(http.StatusOK, resp.Code)
	r.Exactly(userOneExpected, userOneResponse)
}

func TestAPI_GetUser(t *testing.T) {
	r := require.New(t)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	api := New()
	api.InitRouter()
	userTwo := entity.User{ID: 1, Name: "userTwo", Balance: 100}
	api.DB.UsersMap[1] = userTwo
	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	userTwoResponse := entity.User{}
	unmarshal(t, resp.Body.Bytes(), &userTwoResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	userTwoExpected := entity.User{ID: 1, Name: "userTwo", Balance: 100}
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
	userThree := entity.User{ID: 3, Name: "userThree", Balance: 500}
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
	userFour := entity.User{ID: 2, Name: "userFour", Balance: 800}
	api.DB.UsersMap[2] = userFour
	userFourPoints := ReqPoints{Points: 200}
	userFourPointsBytes := marshal(t, &userFourPoints)
	req := httptest.NewRequest(http.MethodPost, "/user/2/take", bytes.NewBuffer(userFourPointsBytes))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	userFourResponse := entity.User{}
	unmarshal(t, resp.Body.Bytes(), &userFourResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	userFourExpected := entity.User{ID: 2, Name: "userFour", Balance: 600}
	r.Exactly(userFourExpected, userFourResponse)

	req = httptest.NewRequest(http.MethodPost, "/user/str/take", bytes.NewBuffer(userFourPointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusBadRequest, resp.Code)

	req = httptest.NewRequest(http.MethodPost, "/user/3/take", bytes.NewBuffer(userFourPointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusNotFound, resp.Code)

	notPoints := "Not Points"
	req = httptest.NewRequest(http.MethodPost, "/user/2/take", bytes.NewBuffer([]byte(notPoints)))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)

	userFourPoints = ReqPoints{Points: 700}
	userFourPointsBytes = marshal(t, userFourPoints)
	req = httptest.NewRequest(http.MethodPost, "/user/2/take", bytes.NewBuffer(userFourPointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	fourResp := entity.Error{}
	unmarshal(t, resp.Body.Bytes(), &fourResp)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)
	fourExpected := entity.Error{Type: "balance is not enough"}
	r.Exactly(fourExpected.Type, fourResp.Type)
}

func TestAPI_FundUserPoints(t *testing.T) {
	r := require.New(t)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	api := New()
	api.InitRouter()
	userFive := entity.User{ID: 2, Name: "userFive", Balance: 200}
	api.DB.UsersMap[2] = userFive
	userFivePoints := ReqPoints{Points: 400}
	userFivePointsBytes := marshal(t, userFivePoints)
	req := httptest.NewRequest(http.MethodPost, "/user/2/fund", bytes.NewBuffer(userFivePointsBytes))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	userFiveResponse := entity.User{}
	unmarshal(t, resp.Body.Bytes(), &userFiveResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	userFiveExpected := entity.User{ID: 2, Name: "userFive", Balance: 600}
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
	userDRByte := marshal(t, &userDR)
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userDRByte))
			api.Router.ServeHTTP(resp, req)
		}(i)
	}
	wg.Wait()
	r.Equal(100, api.DB.CountUsers())
	userDRPointTake := ReqPoints{Points: 1}
	pointTakeBytes := marshal(t, &userDRPointTake)
	userDRPointFund := ReqPoints{Points: 2}
	pointFundBytes := marshal(t, &userDRPointFund)

	var wg2 sync.WaitGroup
	wg2.Add(200)

	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg2.Done()
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/user/1/take", bytes.NewBuffer(pointTakeBytes))
			api.Router.ServeHTTP(resp, req)
		}(i)
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg2.Done()
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/user/1/fund", bytes.NewBuffer(pointFundBytes))
			api.Router.ServeHTTP(resp, req)
		}(i)
	}

	wg2.Wait()
	u, errGet := api.DB.GetUser(1)
	r.NoError(errGet)
	r.Equal(300, u.Balance)
}
