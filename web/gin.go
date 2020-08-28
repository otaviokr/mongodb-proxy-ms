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
	router.POST("/insert/:db/:collection", ws.insert)
	router.POST("/find/:db/:collection", ws.Find)
	router.POST("/update/:db/:collection", ws.update)

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
	result, err := w.mongo.HealthCheck()
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to connect to mongodb")
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, struct {
		Databases []string `json:"databases"`
	}{result})
}

func (w *Server) insert(c *gin.Context) {
	database := c.Params.ByName("db")
	collection := c.Params.ByName("collection")

	request, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error reading request body")
		c.JSON(http.StatusBadRequest, "")
		return
	}

	result, err := w.mongo.Insert(database, collection, request)
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

	if len(strings.TrimSpace(database)) == 0 {
		c.Writer.Header().Add("Content-Type", "application/json;charset=utf-8")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.WriteString("Missing database name")
		return
	}

	if len(strings.TrimSpace(collection)) == 0 {
		c.Writer.Header().Add("Content-Type", "application/json;charset=utf-8")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.WriteString("Missing collection name")
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

	result, err := w.mongo.Find(database, collection, filter)
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

func (w *Server) update(c *gin.Context) {
	database := c.Params.ByName("db")
	collection := c.Params.ByName("collection")

	request, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error reading request body")
		c.JSON(http.StatusBadRequest, "")
		return
	}

	result, err := w.mongo.Update(database, collection, request)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error while updating data...")
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, result)
}
