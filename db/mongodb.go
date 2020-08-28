package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"encoding/json"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// MongoDBProxy manages everything related to MongoDB connection, queries etc.
type MongoDBProxy struct {
	hostname string
	port     int
	URI      string
}

// NewConnection instantiates the MongoDB proxy connector (client, context etc.)
func NewConnection(hostname string, port int, username, password string) (*MongoDBProxy, error) {
	// port 27017
	URI := fmt.Sprintf("mongodb://%s%s%s",
		getUserCredentialForConnectionString(username, password),
		hostname,
		validatePortForConnectionString(port))

	return &MongoDBProxy{
		hostname: hostname,
		port:     port,
		URI:      URI,
	}, nil
}

// Insert will create a new document in collection collName in database dbName.
// Insert("okr", "okr_coll", []byte(`{"id": 1,"name": "A green door","price": 12.50,"tags": ["home", "green"]}`), *client, ctx)
func (m *MongoDBProxy) Insert(dbName, collName string, JSONString []byte) ([]byte, error) {
	var bdoc interface{}
	bson.UnmarshalJSON(JSONString, &bdoc)

	client, ctx, cancelFunc, err := m.getConnection()
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to connect to database")
		return []byte{}, err
	}
	defer cancelFunc()
	defer client.Disconnect(ctx)

	r, err := client.Database(dbName).Collection(collName).InsertOne(ctx, &bdoc)
	if err != nil {
		log.Error().
			Str("database", dbName).
			Str("collection", collName).
			Msgf("failed to insert into database")
		return []byte{}, err
	}

	// TODO fmt.Printf("%+v\n", r)
	log.Info().
		Msgf("created new document: %s", fmt.Sprintf("%+v", r.InsertedID))
	return []byte(fmt.Sprintf("%v", r.InsertedID)), nil
}

// Find will fetch all documents that match filter.
// Find("okr", "okr_coll", []byte(`{ "id": 1 }`), *client, ctx)
func (m *MongoDBProxy) Find(dbName, collName string, filter []byte) ([]byte, error) {
	var bdoc interface{}
	bson.UnmarshalJSON(filter, &bdoc)

	client, ctx, cancelFunc, err := m.getConnection()
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to connect to database")
		return []byte{}, err
	}
	defer cancelFunc()
	defer client.Disconnect(ctx)

	cursor, err := client.Database(dbName).Collection(collName).Find(ctx, &bdoc)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to search in database")
		return []byte{}, err
	}

	var items []string
	for cursor.Next(ctx) {
		item := cursor.Current
		items = append(items, fmt.Sprintf("%+v\n", item))
	}

	results, err := json.Marshal(items)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to parse documents found")
		return []byte{}, err
	}

	return results, nil
}

// Update will modify the fields defined in update in all documents that match filter.
func (m *MongoDBProxy) Update(dbName, collName string, request []byte) ([]byte, error) {
	var parsed struct {
		Filter interface{} `json:"filter"`
		Update interface{} `json:"update"`
	}

	err := json.Unmarshal(request, &parsed)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to parse request")
		return []byte{}, err
	}

	client, ctx, cancelFunc, err := m.getConnection()
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to connect to database")
		return []byte{}, err
	}
	defer cancelFunc()
	defer client.Disconnect(ctx)

	result, err := client.Database(dbName).Collection(collName).UpdateMany(ctx, parsed.Filter, parsed.Update)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to perform update in database")
		return []byte{}, err
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Error().
			Msgf("failed to parse the result as a JSON string")
		return []byte{}, err
	}

	return resultJSON, nil
}

// HealthCheck will return the existing databases if connection is OK.
func (m *MongoDBProxy) HealthCheck() ([]string, error) {
	client, ctx, cancelFunc, err := m.getConnection()
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to connect to database")
		return []string{}, err
	}
	defer cancelFunc()
	defer client.Disconnect(ctx)

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to get database names")
		return []string{}, err
	}
	return databases, nil
}

func (m *MongoDBProxy) getConnection() (*mongo.Client, context.Context, context.CancelFunc, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(m.URI))
	if err != nil {
		log.Error().
			Msgf("error connecting to database")
		return nil, nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Error().
			Msgf("failed to get context")
		cancel()
		return nil, nil, nil, err
	}

	return client, ctx, cancel, nil
}

func getUserCredentialForConnectionString(username, password string) string {
	if len(strings.TrimSpace(username)) == 0 || len(strings.TrimSpace(password)) == 0 {
		log.Warn().
			Str("username", username).
			Str("password", password).
			Msg("ignoring username/password because at least one is empty")
		return ""
	}

	return fmt.Sprintf("%s:%s@", getSanitizedString(username), getSanitizedString(password))
}

func validatePortForConnectionString(p int) string {
	if p > 0 {
		return fmt.Sprintf(":%d", p)
	}
	return ""
}

func getSanitizedString(s string) string {
	result := strings.ReplaceAll(s, "%", "%25")
	result = strings.ReplaceAll(result, "@", "%40")
	result = strings.ReplaceAll(result, ":", "%3A")
	result = strings.ReplaceAll(result, "/", "%2F")
	return result
}
