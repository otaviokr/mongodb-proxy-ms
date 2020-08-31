package swagger

import (
	"github.com/otaviokr/mongodb-proxy-ms/db"
)

// swagger:route GET / home
// Home is just a testing endpoit to see if server is running.
// responses:
//   200: HomeResponse just a dummy JSON

// This text will appear as description of your response body.
// swagger:response home
type homeResponseWrapper struct {
	// in:body
	Body db.HomeResponse
}
