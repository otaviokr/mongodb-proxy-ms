package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/otaviokr/mongodb-proxy-ms/db"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

// Server wraps everything related to web server we provide.
type Server struct {
	Router *gin.Engine
	mongo  db.Proxy
}

// DatabaseDetailsURI holds the database information passed in URI.
type DatabaseDetailsURI struct {
	Database   string `json:"Database" uri:"Database" binding:"required"`
	Collection string `json:"Collection" uri:"Collection" binding:"required"`
}

// UpdateRequest contains both the filter and the updates to be sent in a request.
type UpdateRequest struct {
	Filter  interface{} `json:"filter"`
	Updates db.Quote    `json:"updates"`
}

// NewCustom creates a new instance of Server.
func NewCustom(router *gin.Engine, mongo db.Proxy) *Server {
	return &Server{
		Router: router,
		mongo:  mongo,
	}
}

// New creates a new instance of a WebServer.
func New(dbHost string, dbPort int, dbUser, dbPass string) *Server {

	mongo, err := db.NewConnection(dbHost, dbPort, dbUser, dbPass)
	if err != nil {
		// If connection failed, certainly health() should fail too
		log.Error().
			Err(err).
			Msgf("failed to connect to mongodb at %s:%d", dbHost, dbPort)

		return nil
	}

	return NewWithCustomDB(mongo)
}

// NewWithCustomDB creates a new instance of a WebServer, with a custom DB handler.
func NewWithCustomDB(mongo db.Proxy) *Server {
	router := gin.Default()
	router.Use(cors.Default())

	ws := &Server{
		Router: router,
		mongo:  mongo,
	}

	router.GET("/", ws.Home)
	router.GET("/health", ws.Health)
	router.POST("/insert/:Database/:Collection", ws.Insert)
	router.POST("/find/:Database/:Collection", ws.Find)
	router.POST("/update/:Database/:Collection", ws.Update)

	return ws
}

// Run starts the webserver on address.
func (w *Server) Run(address string) {
	err := w.Router.Run(address)
	log.Error().
		Err(err).
		Msg("error while running the webserver")
}

// Home serves requests for home (index).
func (w *Server) Home(c *gin.Context) {
	response := db.HomeResponse{
		Hello: "World",
	}
	c.JSON(http.StatusOK, response)
}

// Health return the names of available databases or err is DB is down.
func (w *Server) Health(c *gin.Context) {
	result, err := w.mongo.HealthCheck()
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to connect to mongodb")
		c.JSON(http.StatusInternalServerError, db.HealthResponse{})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Insert creeates a new entry in the database.
func (w *Server) Insert(c *gin.Context) {
	var databaseDetails DatabaseDetailsURI
	err := c.ShouldBindUri(&databaseDetails)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to parse URI")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
	}

	request, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error reading request body")
		c.JSON(http.StatusBadRequest, "")
		return
	}

	var quote db.Quote
	err = json.Unmarshal(request, &quote)
	if err != nil {
		panic(err)
	}

	result, err := w.mongo.Insert(databaseDetails.Database, databaseDetails.Collection, quote)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error inserting data into database")
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	parsedJSON := db.InsertResponse{
		InsertedID: result.InsertedID,
	}

	c.JSON(http.StatusOK, parsedJSON)
}

// Find serves requests for fetching data in database.
func (w *Server) Find(c *gin.Context) {
	var databaseDetails DatabaseDetailsURI
	err := c.ShouldBindUri(&databaseDetails)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to parse URI")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
	}

	filter, err := ioutil.ReadAll(c.Request.Body)
	// This error usually is caused by Buffer Overflow.
	// Should I simulate it just for test coverage sake?
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error reading request body")
		c.Writer.Header().Add("Content-Type", "application/json;charset=utf-8")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.WriteString("")
		return
	}

	var filterParsed interface{}

	if len(filter) == 0 {
		log.Debug().Msg("Filter is empty. Adapting it to get all documents")
		filterParsed = bson.M{}
	} else {
		err = bson.UnmarshalExtJSON(filter, true, &filterParsed)
		if err != nil {
			panic(err)
		}
	}

	result, err := w.mongo.Find(databaseDetails.Database, databaseDetails.Collection, filterParsed)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error while searching data in database")
		c.Writer.Header().Add("Content-Type", "application/json;charset=utf-8")
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.WriteString("")
		return
	}

	c.JSON(http.StatusOK, result)
}

// Update changes values in an existing entry in the database.
func (w *Server) Update(c *gin.Context) {
	var databaseDetails DatabaseDetailsURI
	err := c.ShouldBindUri(&databaseDetails)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to parse URI")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
	}

	request, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error reading request body")
		c.JSON(http.StatusBadRequest, "")
		return
	}

	var parsed UpdateRequest
	err = json.Unmarshal(request, &parsed)
	if err != nil {
		panic(err)
	}

	update := bson.D{{"$set", parsed.Updates}}
	result, err := w.mongo.Update(databaseDetails.Database, databaseDetails.Collection, parsed.Filter, update)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error while updating data...")
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, result)
}

// func validateParams(database, collection string) string {
// 	if len(strings.TrimSpace(database)) == 0 {
// 		return "Missing database name"
// 	} else if len(strings.TrimSpace(collection)) == 0 {
// 		return "Missing collection name"
// 	} else {
// 		return ""
// 	}
// }
