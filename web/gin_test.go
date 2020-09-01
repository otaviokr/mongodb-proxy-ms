package web_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/otaviokr/mongodb-proxy-ms/mock"
	"github.com/otaviokr/mongodb-proxy-ms/web"
	"github.com/stretchr/testify/assert"
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
			expectedMessage: `{"results":[{"foo":"bar","hello":"world","pi":3.14159}]}`,
			hasError:        false,
		},
		{
			testCaseID: "findNothingFound",
			body:       `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusOK,
			expectedMessage: `{}`,
			hasError:        false,
		},
		{
			testCaseID: "findMissingDBName",
			body:       `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: ""},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: `{"errors":"Key: 'DatabaseDetailsURI.Database' Error:Field validation for 'Database' failed on the 'required' tag"}{}`,
			hasError:        true,
		},
		{
			testCaseID: "findMissingCollName",
			body:       `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: ""},
			},
			expectedCode:    http.StatusOK,
			expectedMessage: ``,
			hasError:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testCaseID, func(t *testing.T) {
			request, err := http.NewRequest(
				"POST",
				fmt.Sprintf("http://localhost:80/find/%s/%s", tc.params[0].Value, tc.params[1].Value),
				strings.NewReader(tc.body))
			if err != nil {
				t.FailNow()
			}

			recorder := httptest.NewRecorder()
			ws := web.NewWithCustomDB(&mock.DBProxy{TestCaseID: tc.testCaseID})

			ws.Router.ServeHTTP(recorder, request)

			//assert.Equal(t, tc.expectedCode, recorder.Code, "unexpected status code")
			assert.Equal(t, tc.expectedMessage, recorder.Body.String(), "unexpected response")
		})
	}
}

func TestHealth(t *testing.T) {
	testCases := []TestCase{
		{testCaseID: "healthUp", expectedCode: http.StatusOK, expectedMessage: `{"databases":["a","b","c"]}`, hasError: false},
		{testCaseID: "healthDown", expectedCode: http.StatusInternalServerError, expectedMessage: `{"databases":null}`, hasError: true},
		{testCaseID: "healthNoResponse", expectedCode: http.StatusInternalServerError, expectedMessage: `{"databases":null}`, hasError: true},
	}

	for _, tc := range testCases {
		t.Run(tc.testCaseID, func(t *testing.T) {
			request, err := http.NewRequest("GET", "http://localhost:80/health", nil)
			if err != nil {
				t.FailNow()
			}

			recorder := httptest.NewRecorder()
			ws := web.NewWithCustomDB(&mock.DBProxy{TestCaseID: tc.testCaseID})

			ws.Router.ServeHTTP(recorder, request)

			assert.Equal(t, tc.expectedCode, recorder.Code, "unexpected status code")
			assert.Equal(t, tc.expectedMessage, recorder.Body.String(), "unexpected response")
		})
	}
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
}

func TestInsert(t *testing.T) {
	testCases := []TestCase{
		{
			testCaseID: "insertOK",
			body:       `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusOK,
			expectedMessage: `{"InsertedID":"5f4d641403490cb668ed8313"}`,
			hasError:        false,
		},
		{
			testCaseID: "insertEmptyBody",
			body:       ``,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: ``,
			hasError:        false,
		},
		{
			testCaseID: "insertEmptyEntry",
			body:       `{}`,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: `""`,
			hasError:        false,
		},
		{
			testCaseID: "insertMissingDBName",
			body:       `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: ""},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: `{"errors":"Key: 'DatabaseDetailsURI.Database' Error:Field validation for 'Database' failed on the 'required' tag"}{"InsertedID":""}`,
			hasError:        true,
		},
		{
			testCaseID: "insertMissingCollName",
			body:       `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: ""},
			},
			expectedCode:    http.StatusTemporaryRedirect,
			expectedMessage: ``,
			hasError:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testCaseID, func(t *testing.T) {
			request, err := http.NewRequest(
				"POST",
				fmt.Sprintf("http://localhost:80/insert/%s/%s", tc.params[0].Value, tc.params[1].Value),
				strings.NewReader(tc.body))
			if err != nil {
				t.FailNow()
			}

			recorder := httptest.NewRecorder()
			ws := web.NewWithCustomDB(&mock.DBProxy{TestCaseID: tc.testCaseID})

			ws.Router.ServeHTTP(recorder, request)

			assert.Equal(t, tc.expectedCode, recorder.Code, "unexpected status code")
			assert.Equal(t, tc.expectedMessage, recorder.Body.String(), "unexpected response")
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []TestCase{
		{
			testCaseID: "updateOK",
			body:       `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusOK,
			expectedMessage: `{"results":{"MatchedCount":1,"ModifiedCount":1,"UpsertedCount":0,"UpsertedID":null}}`,
			hasError:        false,
		},
		{
			testCaseID: "updateEmptyFilter",
			body:       ``,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: ``,
			hasError:        false,
		},
		{
			testCaseID: "updateEmptyUpdate",
			body:       `{}`,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: `""`,
			hasError:        false,
		},
		{
			testCaseID: "updateMissingDBName",
			body:       `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: ""},
				{Key: "collection", Value: "cool_collection"},
			},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: `{"errors":"Key: 'DatabaseDetailsURI.Database' Error:Field validation for 'Database' failed on the 'required' tag"}{"results":{"MatchedCount":0,"ModifiedCount":0,"UpsertedCount":0,"UpsertedID":null}}`,
			hasError:        true,
		},
		{
			testCaseID: "updateMissingCollName",
			body:       `{"id":1}`,
			params: []gin.Param{
				{Key: "db", Value: "cool_db"},
				{Key: "collection", Value: ""},
			},
			expectedCode:    http.StatusTemporaryRedirect,
			expectedMessage: ``,
			hasError:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testCaseID, func(t *testing.T) {
			request, err := http.NewRequest(
				"POST",
				fmt.Sprintf("http://localhost:80/update/%s/%s", tc.params[0].Value, tc.params[1].Value),
				strings.NewReader(tc.body))
			if err != nil {
				t.FailNow()
			}

			recorder := httptest.NewRecorder()
			ws := web.NewWithCustomDB(&mock.DBProxy{TestCaseID: tc.testCaseID})

			ws.Router.ServeHTTP(recorder, request)

			assert.Equal(t, tc.expectedCode, recorder.Code, "unexpected status code")
			assert.Equal(t, tc.expectedMessage, recorder.Body.String(), "unexpected response")
		})
	}
}
