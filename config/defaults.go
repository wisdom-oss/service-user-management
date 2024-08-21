//go:build !docker

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
// The go toolchain automatically uses this file if no build tags have been
// set.
// This should be the case in local development and testing.
// To achieve the behavior in docker containers, supply the build tag "docker"
// and the Docker defaults are used

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
	return middlewares
}

// QueryFilePath contains the default file path under which the
// sql queries are stored
const QueryFilePath = "./resources/queries.sql"

// ListenAddress sets the host on which the microservice listens to incoming
// requests
const ListenAddress = "127.0.0.1"

// httpLogger generates a new logger which is configured to use a plain text
// format while logging
func httpLogger() func(next http.Handler) http.Handler {
	l := httplog.NewLogger(globals.ServiceName)
	httplog.Configure(httplog.Options{
		JSON:            false,
		Concise:         true,
		TimeFieldFormat: time.RFC3339,
	})
	return httplog.Handler(l)

}
