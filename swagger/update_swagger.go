package swagger

import (
	"github.com/otaviokr/mongodb-proxy-ms/db"
	"github.com/otaviokr/mongodb-proxy-ms/web"
)

// swagger:route POST /update/{Database}/{Collection} update
// Update changes the values in one or more entries in the collection.
// responses:
//   200: UpdateResponse shows the result of the update

// This text will appear as description of the response body.
// swagger:response update
type updateResponseWrapper struct {
	// in:body
	Body db.UpdateResponse
}

// swagger:parameters update
type updateParamsWrapper struct {
	// This text will appear as description of the request body.

	// in:path
	Database string
	// in:path
	Collection string

	// in:body
	Body web.UpdateRequest
}
