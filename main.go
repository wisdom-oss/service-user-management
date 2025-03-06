package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/wisdom-oss/common-go/v2/middleware"
	"github.com/wisdom-oss/common-go/v2/types"
	healthcheckServer "github.com/wisdom-oss/go-healthcheck/server"
	"golang.org/x/sync/errgroup"

	"microservice/internal"
	"microservice/internal/config"
	"microservice/internal/db"
	"microservice/routes"
	"microservice/routes/clients"
	"microservice/routes/permissions"
	"microservice/routes/users"
)

// the main function bootstraps the http server and handlers used for this
// microservice
func main() {
	// create a new logger for the main function
	l := log.Logger
	l.Info().Msgf("configuring %s service", internal.ServiceName)

	// create the healthcheck server
	hcServer := healthcheckServer.HealthcheckServer{}
	hcServer.InitWithFunc(func() error {
		// test if the database is reachable
		return db.Pool.Ping(context.Background())
	})
	err := hcServer.Start()
	if err != nil {
		l.Fatal().Err(err).Msg("unable to start healthcheck server")
	}
	go hcServer.Run()

	// create jwt validator using localhost to get data
	jwtValidator := middleware.JWTValidator{}
	protect := middleware.RequireScope{}
	requireRead := protect.Gin("user-management", types.ScopeRead)
	requireWrite := protect.Gin("user-management", types.ScopeWrite)
	requireDelete := protect.Gin("user-management", types.ScopeDelete)

	service := config.PrepareRouter()

	service.GET("/login", routes.InitiateLogin)
	service.GET("/callback", routes.Callback)
	service.POST("/token", routes.Token)
	service.POST("/revoke", jwtValidator.GinHandler, routes.RevokeToken)

	wellKnown := service.Group("/.well-known")
	{
		wellKnown.GET("/jwks.json", routes.JWK)
		wellKnown.GET("/openid-configuration", routes.OpenIDConfiguration)
	}

	userManagement := service.Group("/users", jwtValidator.GinHandler)
	{
		userManagement.GET("/:userID", users.Information)
		userManagement.GET("/", requireRead, users.List)
		// userManagement.PATCH("/:userID", protect.Gin("user-management", types.ScopeWrite))   // todo: update user
		userManagement.DELETE("/:userID", requireDelete, users.Delete) // todo: delete user
	}

	permissionManagement := service.Group("/permissions", jwtValidator.GinHandler)
	{
		permissionManagement.PATCH("/assign", requireWrite, permissions.Assign)
		permissionManagement.PATCH("/delete", requireDelete, permissions.Delete)
	}

	clientManagement := service.Group("/clients", jwtValidator.GinHandler)
	{
		clientManagement.POST("/", requireWrite, clients.Create)
		clientManagement.DELETE("/:clientID", requireDelete, clients.Delete)
	}

	externalServer := &http.Server{
		Addr:    config.ListenAddress,
		Handler: service,
	}

	var g errgroup.Group

	// Start the server and log errors that happen while running it
	g.Go(func() error {
		err := externalServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("unable to run http server")
		}
		return err
	})

	// Set up the signal handling to allow the server to shut down gracefully
	cancelSignal := make(chan os.Signal, 1)
	cleanupSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM)

	// start the refresh token cleanup
	go cleanupRefreshTokens(cleanupSignal)

	// Block further code execution until the shutdown signal was received
	l.Info().Msg("server ready to accept connections")

	// configure the JWT validator here to allow it to fetch the JWKS from
	// itself
	err = jwtValidator.Configure("user-management", "http://localhost:8000/.well-known/jwks.json", false)
	if err != nil {
		l.Fatal().Err(err).Msg("An error occurred while configuring the JWT Validator")
	}

	<-cancelSignal
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = externalServer.Shutdown(ctx)
	if err != nil {
		l.Fatal().Err(err).Msg("An error occurred while shutting down http server")
	}

	if err := g.Wait(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			l.Fatal().Err(err).Msg("An error occurred while executing servers")
		}
	}

}

func cleanupRefreshTokens(sig chan os.Signal) {
	query, err := db.Queries.Raw("cleanup-expired-tokens")
	if err != nil {
		log.Warn().Err(err).Msg("unable to cleanup expired refresh tokens")
		return
	}

	ticker := time.Tick(15 * time.Second)
	var leaveLoop = false
	for {
		if leaveLoop {
			break
		}
		select {
		case <-ticker:
			_, err = db.Pool.Exec(context.Background(), query)
			if err != nil {
				log.Warn().Err(err).Msg("unable to cleanup expired refresh tokens")
				continue
			}
			break
		case <-sig:
			leaveLoop = true
			continue
		}
	}
}
