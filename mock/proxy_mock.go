package mock

import (
	"context"
	"fmt"

	"github.com/otaviokr/mongodb-proxy-ms/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DBProxy is an implementation of Proxy, to be used in unit tests.
type DBProxy struct {
	TestCaseID string
	HasError   bool
}

// DBWrapperFunc simulates the output from database.
func (m *DBProxy) DBWrapperFunc(db, clt string, req []byte,
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

// Aggregate simulates an aggregation in database.
func (m *DBProxy) Aggregate(db, clt string, req interface{}) (*db.AggregateResponse, error) {
	// TODO
	return nil, nil
}

// Find simulates the output of MongoDB.Find().
func (m *DBProxy) Find(database, collection string, filter interface{}) (*db.FindResponse, error) {

	var results []bson.M
	var errors string

	switch m.TestCaseID {
	case "findOK":
		results = []bson.M{{"foo": "bar", "hello": "world", "pi": 3.14159}}
		errors = ""
	case "findNothingFound":
		results = []bson.M{}
		errors = ""
	case "findMissingDBName":
		// Not reached.
	case "findMissingCollName":
		// Not reached.
	default:
		results = []bson.M{}
		errors = "Testcase not defined - " + m.TestCaseID
	}

	return &db.FindResponse{
		Results: results,
		Errors:  errors,
	}, nil
}

// GetURI returns the value or URI..
func (m *DBProxy) GetURI() string {
	return "TES_URI"
}

// HealthCheck simulates the output of MongoDB.HealthCheck().
func (m *DBProxy) HealthCheck() (*db.HealthResponse, error) {
	var databases []string
	var err error

	switch m.TestCaseID {
	case "healthUp":
		databases = []string{"a", "b", "c"}
	case "healthDown":
		err = fmt.Errorf("healthdown")
	case "healthNoResponse":
		err = fmt.Errorf("noresponse")
	default:
		err = fmt.Errorf("Unexpected test case: %s", m.TestCaseID)
	}

	return &db.HealthResponse{
		Databases: databases,
	}, err
}

// Insert simulates the output of MongoDB.Insert().
func (m *DBProxy) Insert(database, collection string, entry db.Quote) (*db.InsertResponse, error) {

	var insertedID string
	var err error
	switch m.TestCaseID {
	case "insertOK":
		insertedID = "5f4d641403490cb668ed8313"
	case "insertEmptyBody":
		err = fmt.Errorf("Request is empty")
	case "insertEmptyEntry":
		err = fmt.Errorf("Request has no data")
	case "insertMissingDBName":
		// Not reached.
	case "insertMissingCollName":
		// Not reached.
	default:
		err = fmt.Errorf("Unexpected test case: %s", m.TestCaseID)
	}
	return &db.InsertResponse{
		InsertedID: insertedID,
	}, err
}

// Update simulates the output of MongoDB.Update().
func (m *DBProxy) Update(database, collection string, filter, entry interface{}) (*db.UpdateResponse, error) {

	var updateResult mongo.UpdateResult
	var errors error

	switch m.TestCaseID {
	case "updateOK":
		updateResult = mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1, UpsertedCount: 0, UpsertedID: nil}
	case "updateEmptyFilter":
		updateResult = mongo.UpdateResult{MatchedCount: 100, ModifiedCount: 100, UpsertedCount: 0, UpsertedID: nil}
	case "updateEmptyUpdate":
		errors = fmt.Errorf("No updates given")
	case "updateMissingDBName":
		// Not reached.
	case "updateMissingCollName":
		// Not reached.
	default:
		errors = fmt.Errorf("Unexpected test case name: %s", m.TestCaseID)
	}

	return &db.UpdateResponse{
		Results: &updateResult,
	}, errors
}
