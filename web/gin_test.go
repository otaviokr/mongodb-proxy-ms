package web_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/otaviokr/mongodb-proxy-ms/db"
	"github.com/otaviokr/mongodb-proxy-ms/web"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestCase struct {
	testCaseID      string
	body            string
	params          []gin.Param
	expectedCode    int
	expectedMessage string
	hasError        bool
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHome(t *testing.T) {
	ws := web.Server{}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ws.Home(ctx)
	assert.Equal(t, http.StatusOK, recorder.Code, "unexpected status code")
	assert.Equal(t, `{"hello":"World"}`, recorder.Body.String(), "unexpected response")
	assert.Equal(t,
		"application/json; charset=utf-8",
		recorder.Result().Header.Get("Content-Type"),
		"unexpected Content-Type in Header")
	//assert.True(t, recorder.Flushed, "response was not flushed?")
}

func TestFind(t *testing.T) {
	testCases := []TestCase{
		{
			testCaseID: "findOK",
			body:       `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusOK,
			expectedMessage: `{"errors":"1"}`,
			hasError:        false,
		},
		// {
		// 	body: "",
		// 	params: []gin.Param{
		// 		{Key: "db", Value: "cool_db"},
		// 		{Key: "collection", Value: "cool_collection"},
		// 	},
		// 	expectedCode:    http.StatusOK,
		// 	expectedMessage: `{"id":1}`,
		// 	hasError:        false,
		// },
		// {
		// 	body: `{"id":1}`,
		// 	params: []gin.Param{
		// 		{Key: "db", Value: ""},
		// 		{Key: "collection", Value: "cool_collection"},
		// 	},
		// 	expectedCode:    http.StatusBadRequest,
		// 	expectedMessage: "Missing database name",
		// 	hasError:        false,
		// },
		// {
		// 	body: `{"id":1}`,
		// 	params: []gin.Param{
		// 		{Key: "db", Value: "cool_db"},
		// 		{Key: "collection", Value: ""},
		// 	},
		// 	expectedCode:    http.StatusBadRequest,
		// 	expectedMessage: "Missing collection name",
		// 	hasError:        false,
		// },
		// {
		// 	body: `{"id":1}`,
		// 	params: []gin.Param{
		// 		{Key: "ax", Value: "cool_ax"},
		// 	},
		// 	expectedCode:    http.StatusBadRequest,
		// 	expectedMessage: "Missing database name",
		// 	hasError:        false,
		// },
	}

	for _, tc := range testCases {
		ws := web.NewCustom(nil, &MockDBProxy{TestCaseID: tc.testCaseID})

		request, err := http.NewRequest("POST", "x", strings.NewReader(tc.body))
		if err != nil {
			t.FailNow()
		}

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Params = append(ctx.Params, tc.params...)
		ctx.Request = request

		ws.Find(ctx)

		assert.Equal(t, tc.expectedCode, recorder.Code, "unexpected status code")
		assert.Equal(t, tc.expectedMessage, recorder.Body.String(), "unexpected response")
	}
}

func TestHealth(t *testing.T) {
	testCases := []TestCase{
		{testCaseID: "healthUp", expectedCode: http.StatusOK, expectedMessage: `{"databases":["a","b","c"]}`, hasError: false},
		{testCaseID: "healthDown", expectedCode: http.StatusInternalServerError, expectedMessage: `{"databases":null}`, hasError: true},
		{testCaseID: "healthNoResponse", expectedCode: http.StatusNoContent, expectedMessage: "", hasError: true},
	}

	for _, tc := range testCases {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		ws := web.NewCustom(nil, &MockDBProxy{TestCaseID: tc.testCaseID, HasError: tc.hasError})
		ws.Health(ctx)

		assert.Equal(t, tc.expectedCode, recorder.Code, "unexpected status code")
		assert.Equal(t, tc.expectedMessage, recorder.Body.String(), "unexpected response")
	}
}

type MockDBProxy struct {
	TestCaseID string
	HasError   bool
}

func (m *MockDBProxy) DBWrapperFunc(db, clt string, req []byte,
	f func(ctx context.Context, c *mongo.Client, db, clt string, req []byte) ([]byte, error)) ([]byte, error) {
	switch m.TestCaseID {
	case "findOK":
		return []byte(`{"errors":"1"}`), nil
	case "healthUp":
		return []byte(`{"databases":["a","b","c"]}`), nil
	case "healthDown":
		return []byte{}, fmt.Errorf("error as expected in the test case")
	case "healthNoResponse":
		return []byte{}, nil
	}
	return []byte{}, nil
}

func (m *MockDBProxy) Find(database, collection string, filter interface{}) (*db.FindResponse, error) {
	return nil, nil
}

func (m *MockDBProxy) HealthCheck() (*db.HealthResponse, error) {
	return nil, nil
}

func (m *MockDBProxy) Insert(database, collection string, entry db.Quote) (*db.InsertResponse, error) {
	return &db.InsertResponse{
		InsertID: nil,
	}, nil
}

func (m *MockDBProxy) Update(database, collection string, filter, entry interface{}) (*db.UpdateResponse, error) {
	return nil, nil
}
