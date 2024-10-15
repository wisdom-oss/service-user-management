//go:build !docker

// This file contains a default configuration for local development.
// It automatically disables the requirement for authenticated calls and sets
// the default logging level to Debug to make debugging easier.
// Furthermore, the default listen address is set to localhost only as the
// authentication has been turned off to minimize the risk of data leaks

package config

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/wisdom-oss/common-go/v2/middleware"
)
import "github.com/gin-contrib/logger"
import "github.com/gin-contrib/requestid"

const ListenAddress = "127.0.0.1:8000"

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
	return middlewares
}
