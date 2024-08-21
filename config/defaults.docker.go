//go:build docker

package config

import (
	"net/http"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/wisdom-oss/common-go/middleware"

	"microservice/globals"
)

// This file contains default paths that are used inside the service to load
// resource files.
// The go toolchain automatically uses this file if the build tag "docker" has
// been set.
// This should be the case in docker images and containers only.
// To achieve the behavior for local development and testing, please remove the
// "docker" build tag

// Middlewares contains the middlewares used per default in this service.
// To disable single middlewares, please remove the line in which this array
// is used and add the middlewares that shall be used manually to the router
var middlewares = []func(next http.Handler) http.Handler{
	httpLogger(),
	chiMiddleware.RequestID,
	chiMiddleware.RealIP,
	middleware.ErrorHandler,
}

func Middlewares() []func(http.Handler) http.Handler {
	validator := middleware.JWTValidator{}
	_ = validator.Configure("api-gateway")
	middlewares := append(middlewares, validator.Handler)
	return middlewares
}

// QueryFilePath contains the default file path under which the
// sql queries are stored
const QueryFilePath = "./queries.sql"

// ListenAddress sets the host on which the microservice listens to incoming
// requests
const ListenAddress = ""

func httpLogger() func(next http.Handler) http.Handler {
	l := httplog.NewLogger(globals.ServiceName)
	httplog.Configure(httplog.Options{
		JSON:            true,
		Concise:         true,
		TimeFieldFormat: time.RFC3339,
	})
	return httplog.Handler(l)

}
