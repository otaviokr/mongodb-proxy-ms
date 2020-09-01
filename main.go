package main

import (
	"os"
	"strconv"

	"github.com/otaviokr/mongodb-proxy-ms/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/otaviokr/mongodb-proxy-ms/swagger" // only to triggger compilation; not used in the code here
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Info().Msg("Starting log...")

	dbUsername := os.Getenv("MONGODB_USER")
	dbPassword := os.Getenv("MONGODB_PASS")
	dbHostname := os.Getenv("MONGODB_HOST")
	portAsString := os.Getenv("MONGODB_PORT")
	var err error
	dbPort, err := strconv.Atoi(portAsString)
	if err != nil {
		log.Warn().
			Str("MONGODB_PORT", os.Getenv(("MONGODB_PORT"))).
			Err(err).
			Msg("Assuming default value: 27017")
		dbPort = 27017
	}

	router := web.New(dbHostname, dbPort, dbUsername, dbPassword)
	router.Run(":8080")
}
