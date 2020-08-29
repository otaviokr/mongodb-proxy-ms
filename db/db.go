package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// Proxy is the abstraction of what you can do with the database.
type Proxy interface {
	DBWrapperFunc(db, clt string, req []byte,
		f func(ctx context.Context, c *mongo.Client, db, clt string, req []byte) ([]byte, error)) ([]byte, error)
}
