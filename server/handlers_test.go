package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/yanrishbe/gaming-website/logger"

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
	require.NoError(t, errResponse)
}

func TestAPI_RegisterNewUser(t *testing.T) {
	r := require.New(t)
	log := logger.New("debug")
	api := New(log)
	api.InitRouter()
	var userOne = entity.User{Name: "", Balance: 400}
	var userOneByte = marshal(t, userOne)
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userOneByte))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userOneResponse = &UserResp{}
	unmarshal(t, resp.Body.Bytes(), userOneResponse)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)
	var userOneExpected = &UserResp{ID: 0, User: entity.User{Name: "", Balance: 400}, Error: "user's data is not valid"}
	r.Exactly(userOneExpected, userOneResponse)

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
	userOneExpected = &UserResp{ID: 1, User: entity.User{Name: "userOne", Balance: 100}, Error: ""}
	r.EqualValues(http.StatusCreated, resp.Code)
	r.Exactly(userOneExpected, userOneResponse)
}

func TestAPI_GetUser(t *testing.T) {
	r := require.New(t)
	log := logger.New("debug")
	api := New(log)
	api.InitRouter()
	var userTwo = &entity.User{Name: "userTwo", Balance: 100}
	api.DB.UsersMap[1] = userTwo
	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userTwoResponse = &UserResp{}
	unmarshal(t, resp.Body.Bytes(), userTwoResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	var userTwoExpected = &UserResp{ID: 1, User: entity.User{Name: "userTwo", Balance: 100}, Error: ""}
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
	log := logger.New("debug")
	api := New(log)
	api.InitRouter()
	var userThree = &entity.User{Name: "userThree", Balance: 500}
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
	log := logger.New("debug")
	api := New(log)
	api.InitRouter()
	var userFour = &entity.User{Name: "userFour", Balance: 800}
	api.DB.UsersMap[2] = userFour
	var userFourPoints = ReqPoints{Points: 200}
	var userFourPointsBytes = marshal(t, userFourPoints)
	req := httptest.NewRequest(http.MethodPost, "/user/2/take", bytes.NewBuffer(userFourPointsBytes))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userFourResponse = &UserResp{}
	unmarshal(t, resp.Body.Bytes(), userFourResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	var userFourExpected = &UserResp{ID: 2, User: entity.User{Name: "userFour", Balance: 600}, Error: ""}
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
	userFourExpected = &UserResp{ID: 0, User: entity.User{Name: "", Balance: 0},
		Error: "not enough balance to execute the request"}
	r.Exactly(userFourExpected, userFourResponse)

}

func TestAPI_FundUserPoints(t *testing.T) {
	r := require.New(t)
	log := logger.New("debug")
	api := New(log)
	api.InitRouter()
	var userFive = &entity.User{Name: "userFive", Balance: 200}
	api.DB.UsersMap[2] = userFive
	var userFivePoints = ReqPoints{Points: 400}
	var userFivePointsBytes = marshal(t, userFivePoints)
	req := httptest.NewRequest(http.MethodPost, "/user/2/fund", bytes.NewBuffer(userFivePointsBytes))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	var userFiveResponse = &UserResp{}
	unmarshal(t, resp.Body.Bytes(), userFiveResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	var userFiveExpected = &UserResp{ID: 2, User: entity.User{Name: "userFive", Balance: 600}, Error: ""}
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
	log := logger.New("debug")
	api := New(log)
	api.InitRouter()
	userDR := &entity.User{Name: "User", Balance: 500}
	userDRByte := marshal(t, userDR)
	// req := make ([]*http.Request,100)
	//req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userDRByte))

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer func() {
				log.Println("worker:", i, "total users", api.DB.CountUsers())
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
		go func(i int) {
			defer wg2.Done()
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/user/1/take", bytes.NewBuffer(userDRBytes))
			api.Router.ServeHTTP(resp, req)
		}(i)
	}
	wg2.Wait()
	userDRResult, errGet := api.DB.GetUser(1)
	r.NoError(errGet)
	r.Equal(100, userDRResult.Balance)
}
