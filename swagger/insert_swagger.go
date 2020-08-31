package swagger

import (
	"github.com/otaviokr/mongodb-proxy-ms/db"
)

// swagger:route POST /insert/{Database}/{Collection} insert
// Insert adds a new entry in the collection.
// responses:
//   200: InsertResponse shows the result of the insert

// This text will appear as description of the response body.
// swagger:response insert
type insertResponseWrapper struct {
	// in:body
	Body db.InsertResponse
}

// swagger:parameters insert
type insertParamsWrapper struct {
	// This text will appear as description of the request body.

	// in:path
	Database string
	// in:path
	Collection string

	// in:body
	Body db.Quote
}
