package swagger

import "github.com/otaviokr/mongodb-proxy-ms/db"

// swagger:route GET /health health
// Health displays the available databases.
// responses:
//   200: HealthResponse list of available databases

// This text will appear as description of your response body.
// swagger:response health
type healthResponseWrapper struct {
	// in:body
	Body db.HealthResponse
}
