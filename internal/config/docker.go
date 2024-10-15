//go:build docker

// This file contains a default configuration for local development.
// It automatically disables the requirement for authenticated calls and sets
// the default logging level to Debug to make debugging easier.
// Furthermore, the default listen address is set to localhost only as the
// authentication has been turned off to minimize the risk of data leaks

package config

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wisdom-oss/common-go/v2/middleware"
)
import "github.com/gin-contrib/logger"
import "github.com/gin-contrib/requestid"

const ListenAddress = "0.0.0.0:8000"

func init() {
	// set gin to the production mode (aka release mode)
	gin.SetMode(gin.ReleaseMode)
}

// Middlewares configures and outputs the middlewares used in the configuration.
// The contained middlewares are the following:
//   - gin.Logger
func Middlewares() []gin.HandlerFunc {
	var middlewares []gin.HandlerFunc

	middlewares = append(middlewares,
		logger.SetLogger(
			logger.WithDefaultLevel(zerolog.DebugLevel),
			logger.WithUTC(false),
		))

	middlewares = append(middlewares, requestid.New())
	middlewares = append(middlewares, middleware.ErrorHandler{}.Gin)
	middlewares = append(middlewares, gin.CustomRecovery(middleware.RecoveryHandler))

	// read the OpenID Connect issuer from the environment
	oidcIssuer := os.Getenv("OIDC_AUTHORITY")
	jwtValidator := middleware.JWTValidator{}
	err := jwtValidator.DiscoverAndConfigure(oidcIssuer)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to discover and configure JWT Validation")
	}

	middlewares = append(middlewares, jwtValidator.GinHandler)
	return middlewares
}
