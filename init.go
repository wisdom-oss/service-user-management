package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"

	"github.com/qustavo/dotsql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"microservice/config"
	"microservice/globals"

	_ "github.com/wisdom-oss/go-healthcheck/client"
)

// init is executed at every startup of the microservice and is always executed
// before main
func init() {
	configureLogger()
	connectDatabase()
	loadPreparedQueries()
	log.Info().Msg("initialization process finished")
}

// configureLogger handles the configuration of the logger used in the
// microservice. it reads the logging level from the `LOG_LEVEL` environment
// variable and sets it according to the parsed logging level. if an invalid
// value is supplied or no level is supplied, the service defaults to the
// `INFO` level
func configureLogger() {
	// set the time format to unix timestamps to allow easier machine handling
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// allow the logger to create an error stack for the logs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// now use the environment variable `LOG_LEVEL` to determine the logging
	// level for the microservice.
	rawLoggingLevel, isSet := os.LookupEnv("LOG_LEVEL")

	// if the value is not set, use the info level as default.
	var loggingLevel zerolog.Level
	if !isSet {
		loggingLevel = zerolog.InfoLevel
	} else {
		// now try to parse the value of the raw logging level to a logging
		// level for the zerolog package
		var err error
		loggingLevel, err = zerolog.ParseLevel(rawLoggingLevel)
		if err != nil {
			// since an error occurred while parsing the logging level, use info
			loggingLevel = zerolog.InfoLevel
			log.Warn().Msg("unable to parse value from environment. using info")
		}
	}
	// since now a logging level is set, configure the logger
	zerolog.SetGlobalLevel(loggingLevel)
}

// connectDatabase uses the previously read environment variables to connect the
// microservice to the PostgreSQL database used as the backend for all WISdoM
// services
func connectDatabase() {
	log.Info().Msg("connecting to the database")

	var err error
	globals.Db, err = pgxpool.New(context.Background(), "")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create database connection pool")
	}
	err = globals.Db.Ping(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("unable to verify the connection to the database")
	}
	log.Info().Msg("database connection established")
}

// loadPreparedQueries loads the prepared SQL queries from a file specified by
// the QUERY_FILE_LOCATION environment variable or the config.QueryFilePath
// constant.
// It initializes the globals.SqlQueries variable with the loaded queries.
// If there is an error loading the queries, it logs a fatal error and the
// program terminates.
// This function is typically called during the startup of the microservice.
func loadPreparedQueries() {
	log.Info().Msg("loading prepared sql queries")
	var err error
	location, locationChanged := os.LookupEnv("QUERY_FILE_LOCATION")
	if !locationChanged {
		// since the location has not changed, set the default value
		location = config.QueryFilePath
		log.Debug().Str("default", config.QueryFilePath).Msg("location for sql query file has not been changed")
	}
	globals.SqlQueries, err = dotsql.LoadFromFile(location)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load prepared queries")
	}
}
