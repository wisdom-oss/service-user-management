//go:build docker

// This file contains a default configuration for local development.
// It automatically disables the requirement for authenticated calls and sets
// the default logging level to Debug to make debugging easier.
// Furthermore, the default listen address is set to localhost only as the
// authentication has been turned off to minimize the risk of data leaks

package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/wisdom-oss/common-go/v2/middleware"

	apiErrors "microservice/internal/errors"

	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
)

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

	middlewares = append(middlewares, logger.SetLogger(
		logger.WithDefaultLevel(zerolog.DebugLevel),
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		}),
	))

	middlewares = append(middlewares, requestid.New())
	middlewares = append(middlewares, middleware.ErrorHandler{}.Gin)
	middlewares = append(middlewares, gin.CustomRecovery(middleware.RecoveryHandler))

	return middlewares
}

func PrepareRouter() *gin.Engine {
	router := gin.New()
	router.HandleMethodNotAllowed = true
	router.ForwardedByClientIP = true
	_ = router.SetTrustedProxies([]string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})
	router.Use(Middlewares()...)

	router.NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, apiErrors.MethodNotAllowed)
	})
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, apiErrors.NotFound)

	})

	return router
}
