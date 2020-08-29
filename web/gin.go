package web

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/otaviokr/mongodb-proxy-ms/db"
	"github.com/rs/zerolog/log"
)

// Server wraps everything related to web server we provide.
type Server struct {
	router *gin.Engine
	mongo  db.Proxy
}

// NewCustom creates a new instance of Server.
func NewCustom(router *gin.Engine, mongo db.Proxy) *Server {
	return &Server{
		router: router,
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

	router := gin.New()

	ws := &Server{
		router: router,
		mongo:  mongo,
	}

	router.GET("/", ws.Home)
	router.GET("/health", ws.Health)
	router.POST("/insert/:db/:collection", ws.Insert)
	router.POST("/find/:db/:collection", ws.Find)
	router.POST("/update/:db/:collection", ws.Update)

	return ws
}

// Run starts the webserver on address.
func (w *Server) Run(address string) {
	err := w.router.Run(address)
	log.Error().
		Err(err).
		Msg("error while running the webserver")
}

// Home serves requests for home (index).
func (w *Server) Home(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"Hello": "World"})
}

// Health return the names of available databases or err is DB is down.
func (w *Server) Health(c *gin.Context) {
	//result, err := w.mongo.HealthCheck()
	result, err := w.mongo.DBWrapperFunc("", "", nil, db.HealthCheck)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to connect to mongodb")
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.Writer.Header().Add("Content-Type", "application/json;charset=utf-8")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(result)
}

// Insert creeates a new entry in the database.
func (w *Server) Insert(c *gin.Context) {
	database := c.Params.ByName("db")
	collection := c.Params.ByName("collection")

	if msg := validateParams(database, collection); len(msg) > 0 {
		c.Writer.Header().Add("Content-Type", "application/json;charset=utf-8")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.WriteString(msg)
		return
	}

	request, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error reading request body")
		c.JSON(http.StatusBadRequest, "")
		return
	}

	//result, err := w.mongo.Insert(database, collection, request)
	result, err := w.mongo.DBWrapperFunc(database, collection, request, db.Insert)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error inserting data into database")
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, result)
}

// Find serves requests for fetching data in database.
func (w *Server) Find(c *gin.Context) {
	database := c.Params.ByName("db")
	collection := c.Params.ByName("collection")

	if msg := validateParams(database, collection); len(msg) > 0 {
		c.Writer.Header().Add("Content-Type", "application/json;charset=utf-8")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.WriteString(msg)
		return
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

	result, err := w.mongo.DBWrapperFunc(database, collection, filter, db.Find)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error while searching data in database")
		c.Writer.Header().Add("Content-Type", "application/json;charset=utf-8")
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.WriteString("")
		return
	}

	c.Writer.Header().Add("Content-Type", "application/json;charset=utf-8")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(result)
}

// Update changes values in an existing entry in the database.
func (w *Server) Update(c *gin.Context) {
	database := c.Params.ByName("db")
	collection := c.Params.ByName("collection")

	if msg := validateParams(database, collection); len(msg) > 0 {
		c.Writer.Header().Add("Content-Type", "application/json;charset=utf-8")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.WriteString(msg)
		return
	}

	request, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error reading request body")
		c.JSON(http.StatusBadRequest, "")
		return
	}

	result, err := w.mongo.DBWrapperFunc(database, collection, request, db.Update)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error while updating data...")
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, result)
}

func validateParams(database, collection string) string {
	if len(strings.TrimSpace(database)) == 0 {
		return "Missing database name"
	} else if len(strings.TrimSpace(collection)) == 0 {
		return "Missing collection name"
	} else {
		return ""
	}
}
