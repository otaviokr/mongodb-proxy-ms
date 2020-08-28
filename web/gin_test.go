package web_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/otaviokr/mongodb-proxy-ms/web"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
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
	assert.Equal(t, `{"Hello":"World"}`, recorder.Body.String(), "unexpected response")
	assert.Equal(t,
		"application/json; charset=utf-8",
		recorder.Result().Header.Get("Content-Type"),
		"unexpected Content-Type in Header")
	//assert.True(t, recorder.Flushed, "response was not flushed?")
}

func TestFind(t *testing.T) {
	testCases := []TestCase{
		{
			body: `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusOK,
			expectedMessage: `{"id":1}`,
			hasError:        false,
		},
		{
			body: "",
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusOK,
			expectedMessage: `{"id":1}`,
			hasError:        false,
		},
		{
			body: `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: ""},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "Missing database name",
			hasError:        false,
		},
		{
			body: `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: ""},
			},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "Missing collection name",
			hasError:        false,
		},
		{
			body: `{"id":1}`,
			params: []gin.Param{
				{Key: "ax", Value: "cool_ax"},
			},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "Missing database name",
			hasError:        false,
		},
	}

	ws := web.NewCustom(nil, &MockDBProxy{})

	for _, tc := range testCases {
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
		{expectedCode: http.StatusOK, expectedMessage: `{"databases":["a","b","c"]}`, hasError: false},
		{expectedCode: http.StatusInternalServerError, expectedMessage: `""`, hasError: true},
	}

	for _, tc := range testCases {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		ws := web.NewCustom(nil, &MockDBProxy{HasError: tc.hasError})
		ws.Health(ctx)

		assert.Equal(t, tc.expectedCode, recorder.Code, "unexpected status code")
		assert.Equal(t, tc.expectedMessage, recorder.Body.String(), "unexpected response")
	}
}

type MockDBProxy struct {
	HasError bool
}

func (m *MockDBProxy) Insert(dbName, collName string, JSONString []byte) ([]byte, error) {
	// TODO Mock insert to create test case
	return []byte{}, nil
}
func (m *MockDBProxy) Find(dbName, collName string, filter []byte) ([]byte, error) {
	return []byte(`{"id":1}`), nil
}
func (m *MockDBProxy) Update(dbName, collName string, request []byte) ([]byte, error) {
	// TODO Mock update to create test case
	return []byte{}, nil
}
func (m *MockDBProxy) HealthCheck() ([]string, error) {
	if m.HasError {
		return []string{}, fmt.Errorf("Failed as expected in test")
	}
	return []string{"a", "b", "c"}, nil
}