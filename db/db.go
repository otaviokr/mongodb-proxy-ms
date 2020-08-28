package db

// Proxy is the abstraction of what you can do with the database.
type Proxy interface {
	Insert(dbName, collName string, JSONString []byte) ([]byte, error)
	Find(dbName, collName string, filter []byte) ([]byte, error)
	Update(dbName, collName string, request []byte) ([]byte, error)
	HealthCheck() ([]string, error)
}
