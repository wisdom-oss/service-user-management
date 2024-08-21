package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
	wisdomMiddleware "github.com/wisdom-oss/common-go/middleware"
	"github.com/wisdom-oss/common-go/types"
	healthcheckServer "github.com/wisdom-oss/go-healthcheck/server"

	"microservice/config"
	"microservice/routes"

	"microservice/globals"
)

// the main function bootstraps the http server and handlers used for this
// microservice
func main() {
	// create a new logger for the main function
	l := log.With().Str("step", "main").Logger()
	l.Info().Msgf("starting %s service", globals.ServiceName)

	// create the healthcheck server
	hcServer := healthcheckServer.HealthcheckServer{}
	hcServer.InitWithFunc(func() error {
		// test if the database is reachable
		return globals.Db.Ping(context.Background())
	})
	err := hcServer.Start()
	if err != nil {
		l.Fatal().Err(err).Msg("unable to start healthcheck server")
	}
	go hcServer.Run()

	// create a new router
	router := chi.NewRouter()
	// add some middlewares to the router to allow identifying requests
	router.Use(httplog.Handler(l))
	router.Use(config.Middlewares()...)
	router.NotFound(wisdomMiddleware.NotFoundError)
	// now mount the routes as some examples
	router.HandleFunc("/", routes.BasicHandler)
	router.HandleFunc("/internal-error", routes.BasicWithErrorHandling)
	router.With(wisdomMiddleware.RequireScope(globals.ServiceName, types.ScopeRead)).HandleFunc("/read", routes.BasicHandler)
	router.With(wisdomMiddleware.RequireScope(globals.ServiceName, types.ScopeWrite)).HandleFunc("/write", routes.BasicHandler)
	router.With(wisdomMiddleware.RequireScope(globals.ServiceName, types.ScopeDelete)).HandleFunc("/delete", routes.BasicHandler)
	router.With(wisdomMiddleware.RequireScope(globals.ServiceName, types.ScopeAdmin)).HandleFunc("/admin", routes.BasicHandler)

	// now boot up the service
	// Configure the HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", "8000"),
		WriteTimeout: time.Second * 600,
		ReadTimeout:  time.Second * 600,
		IdleTimeout:  time.Second * 600,
		Handler:      router,
	}

	// Start the server and log errors that happen while running it
	go func() {
		if err := server.ListenAndServe(); err != nil {
			l.Fatal().Err(err).Msg("An error occurred while starting the http server")
		}
	}()

	// Set up the signal handling to allow the server to shut down gracefully

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, os.Interrupt)

	// Block further code execution until the shutdown signal was received
	l.Info().Msg("server ready to accept connections")
	<-cancelSignal

}
