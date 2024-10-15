package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	_ "microservice/internal/db" // side effect import to connect to the database and parse the sql queries from it's embed

	_ "github.com/wisdom-oss/go-healthcheck/client"
)

// init is executed at every startup of the microservice and is always executed
// before main
func init() {
	configureLogger()
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
