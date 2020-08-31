package swagger

import "github.com/otaviokr/mongodb-proxy-ms/db"

// swagger:route POST /find/{Database}/{Collection} find
// Find returns entries that match the filter defined..
// responses:
//   200: FindResponse shows the result of the find

// This text will appear as description of the response body.
// swagger:response find
type findResponseWrapper struct {
	// in:body
	Body db.FindResponse
}

// swagger:parameters find
type findParamsWrapper struct {
	// This text will appear as description of the request body.

	// in:path
	Database string
	// in:path
	Collection string

	// in:body
	Body db.Quote
}
