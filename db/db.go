package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Quote represents the central collection of the solution, where the quotes used by the Twitter bot is used.
type Quote struct {
	Publications    int    `json:"publications"`
	LastPublished   int64  `json:"last_published"`
	OriginalQuote   string `json:"original_quote"`
	TranslatedQuote string `json:"translated_quote"`
	Author          string `json:"author"`
}

// HealthResponse shows the databases available.
type HealthResponse struct {
	Databases []string `json:"databases"`
}

// HomeResponse is just a dummy JSON to indicate the server is up.
type HomeResponse struct {
	Hello string `json:"hello"`
}

// InsertResponse gives the ObjectID of the inserted data.
type InsertResponse struct {
	InsertedID interface{} `json:"InsertedID"`
}

// FindResponse returns the data found in database.
type FindResponse struct {
	Results []bson.M `json:"results,omitempty"`
	Errors  string   `json:"errors,omitempty"`
}

// UpdateResponse gives the result of the update.
type UpdateResponse struct {
	Results *mongo.UpdateResult `json:"results"`
}

// Proxy is the abstraction of what you can do with the database.
type Proxy interface {
	GetURI() string
	HealthCheck() (*HealthResponse, error)
	Insert(database, collection string, entry Quote) (*InsertResponse, error)
	Find(database, collection string, filter interface{}) (*FindResponse, error)
	Update(database, collection string, filter, entry interface{}) (*UpdateResponse, error)
}
