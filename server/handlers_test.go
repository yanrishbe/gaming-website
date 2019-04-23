package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/yanrishbe/gaming-website/entity"

	"github.com/stretchr/testify/require"
)

var api *API

func TestMain(m *testing.M) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	var err error
	api, err = New()
	if err != nil {
		logrus.Fatal(err)
	}
	// err = api.DB.CreateTables()
	//if err != nil {
	//	logrus.Fatal(err)
	//}
	// tables are already created in New().
	// Please make sure you have a clear division of abstraction between modules/objects. DB is responsible for
	// initializing itself. Main func or tests had nothing to do with that. //fixme //////////////////

	//api.InitRouter() // you can call InitRouter when you create API.
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

// I will just put one comment for all tests in *Handler* : they are not readable.
// Technically they are good, they are validating the business-logic and they are doing their job.
// But it's very hard to read them!
// It's the same HTTP code boilerplate that should be encapsulated in separate func
// You need to move code that looks the same to separate functions, so that I could see small amount of statements here in test
// like
// I'd advice you to create a "mini-client" that does http requests for you and provides convenient methods to unmarshal and validate response
// u := entity.User{Balance:400}
// mustDo("POST", "/user", u).Code(201).Unmarshal(&u)
// assert.Equal(t,u.Balance,100)
// mustDo("GET", fmt.Sprint("/user/",u.ID)).Code(200).Unmarshal(&u)
// assert.Equal(t,u.Balance,100)
// in that case I'll be able to track what is happening

// we could do this this testing stuff together if you find it difficult

func TestAPI_RegisterNewUser(t *testing.T) {
	r := require.New(t)
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
	r.EqualValues(http.StatusOK, resp.Code)
}

func TestAPI_GetUser(t *testing.T) {
	r := require.New(t)
	userTwo := entity.User{Name: "userTwo", Balance: 400}
	userTwoByte := marshal(t, userTwo)
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userTwoByte))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusOK, resp.Code)

	response := entity.User{}
	unmarshal(t, resp.Body.Bytes(), &response)
	id := strconv.Itoa(response.ID)
	req = httptest.NewRequest(http.MethodGet, "/user/"+id, nil)
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	userTwoResponse := entity.User{}
	unmarshal(t, resp.Body.Bytes(), &userTwoResponse)
	r.EqualValues(http.StatusOK, resp.Code)

	userTwoExpected := entity.User{ID: response.ID, Name: "userTwo", Balance: 100}
	r.Exactly(userTwoExpected, userTwoResponse)

	req = httptest.NewRequest(http.MethodGet, "/user/str", nil)
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusBadRequest, resp.Code)
}

func TestAPI_DeleteUser(t *testing.T) {
	r := require.New(t)
	userThree := entity.User{Name: "userThree", Balance: 500}

	userThreeByte := marshal(t, userThree)
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userThreeByte))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusOK, resp.Code)

	response := entity.User{}
	unmarshal(t, resp.Body.Bytes(), &response)
	id := strconv.Itoa(response.ID)

	req = httptest.NewRequest(http.MethodDelete, "/user/"+id, nil)
	resp = httptest.NewRecorder()
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
	userFour := entity.User{Name: "userFour", Balance: 800}

	userFourByte := marshal(t, userFour)
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userFourByte))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusOK, resp.Code)

	response := entity.User{}
	unmarshal(t, resp.Body.Bytes(), &response)
	id := strconv.Itoa(response.ID)

	userFourPoints := ReqPoints{Points: 200}
	userFourPointsBytes := marshal(t, &userFourPoints)
	req = httptest.NewRequest(http.MethodPost, "/user/"+id+"/take", bytes.NewBuffer(userFourPointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	userFourResponse := entity.User{}
	unmarshal(t, resp.Body.Bytes(), &userFourResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	userFourExpected := entity.User{ID: response.ID, Name: "userFour", Balance: 300}
	r.Exactly(userFourExpected, userFourResponse)

	req = httptest.NewRequest(http.MethodPost, "/user/str/take", bytes.NewBuffer(userFourPointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusBadRequest, resp.Code)

	notPoints := "Not Points"
	req = httptest.NewRequest(http.MethodPost, "/user/"+id+"/take", bytes.NewBuffer([]byte(notPoints)))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)

	userFourPoints = ReqPoints{Points: 700}
	userFourPointsBytes = marshal(t, userFourPoints)
	req = httptest.NewRequest(http.MethodPost, "/user/"+id+"/take", bytes.NewBuffer(userFourPointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	fourResp := entity.Error{}
	unmarshal(t, resp.Body.Bytes(), &fourResp)
	r.EqualValues(http.StatusServiceUnavailable, resp.Code)
	fourExpected := entity.Error{Type: "database error"}
	r.Exactly(fourExpected.Type, fourResp.Type)
}

func TestAPI_FundUserPoints(t *testing.T) {
	r := require.New(t)
	userFive := entity.User{Name: "userFive", Balance: 300}

	userFiveByte := marshal(t, userFive)
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userFiveByte))
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusOK, resp.Code)

	response := entity.User{}
	unmarshal(t, resp.Body.Bytes(), &response)
	id := strconv.Itoa(response.ID)

	userFivePoints := ReqPoints{Points: 400}
	userFivePointsBytes := marshal(t, userFivePoints)
	req = httptest.NewRequest(http.MethodPost, "/user/"+id+"/fund", bytes.NewBuffer(userFivePointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	userFiveResponse := entity.User{}
	unmarshal(t, resp.Body.Bytes(), &userFiveResponse)
	r.EqualValues(http.StatusOK, resp.Code)
	userFiveExpected := entity.User{ID: response.ID, Name: "userFive", Balance: 400}
	r.Exactly(userFiveExpected, userFiveResponse)

	req = httptest.NewRequest(http.MethodPost, "/user/str/fund", bytes.NewBuffer(userFivePointsBytes))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusBadRequest, resp.Code)

	var notPoints = "Not Points"
	req = httptest.NewRequest(http.MethodPost, "/user/"+id+"/fund", bytes.NewBuffer([]byte(notPoints)))
	resp = httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusUnprocessableEntity, resp.Code)
}

// You already testing DataRace in DB, so you don't need to do it here :)
func TestAPI_DataRace(t *testing.T) {
	r := require.New(t)
	userDR := entity.User{Name: "User", Balance: 500}
	userDRByte := marshal(t, &userDR)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userDRByte))
	api.Router.ServeHTTP(resp, req)
	r.EqualValues(http.StatusOK, resp.Code)

	response := entity.User{}
	unmarshal(t, resp.Body.Bytes(), &response)
	id := strconv.Itoa(response.ID)

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
			req := httptest.NewRequest(http.MethodPost, "/user/"+id+"/take", bytes.NewBuffer(pointTakeBytes))
			api.Router.ServeHTTP(resp, req)
			if resp.Result().StatusCode != 200 {
				panic(resp.Body.String())
			}
		}(i)
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg2.Done()
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/user/"+id+"/fund", bytes.NewBuffer(pointFundBytes))
			api.Router.ServeHTTP(resp, req)
			if resp.Result().StatusCode != 200 {
				panic(resp.Body.String())
			}
		}(i)
	}

	wg2.Wait()
	u, errGet := api.DB.GetUser(response.ID)
	r.NoError(errGet)
	r.Equal(300, u.Balance)
}
