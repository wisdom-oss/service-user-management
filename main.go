package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/wisdom-oss/common-go/v2/middleware"
	"github.com/wisdom-oss/common-go/v2/types"
	healthcheckServer "github.com/wisdom-oss/go-healthcheck/server"

	"microservice/internal"
	"microservice/internal/config"
	"microservice/internal/db"
	"microservice/internal/errors"
	"microservice/routes"
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

	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.Use(config.Middlewares()...)
	r.NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, errors.MethodNotAllowed)
	})
	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, errors.NotFound)

	})

	r.GET("/login", routes.InitiateLogin)
	r.GET("/callback", routes.Callback)
	r.POST("/token", routes.Token)
	r.POST("/revoke", jwtValidator.GinHandler, routes.RevokeToken)

	wellKnown := r.Group("/.well-known")
	{
		wellKnown.GET("/jwks.json", routes.JWK)
	}

	userManagement := r.Group("/users", jwtValidator.GinHandler)
	{
		userManagement.GET("/:userID", users.Information)
		userManagement.GET("/", protect.Gin("user-management", types.ScopeRead), users.List)
		// userManagement.PUT("/", protect.Gin("user-management", types.ScopeWrite)) // todo: write route to create new user
		// userManagement.PATCH("/:userID", protect.Gin("user-management", types.ScopeWrite))   // todo: update user
		// userManagement.DELETE("/:userID", protect.Gin("user-management", types.ScopeDelete)) // todo: delete user
	}

	l.Info().Msg("finished service configuration")
	l.Info().Msg("starting http server")

	// Start the server and log errors that happen while running it
	go func() {
		if err := r.Run(config.ListenAddress); err != nil {
			l.Fatal().Err(err).Msg("An error occurred while starting the http server")
		}
	}()

	// Set up the signal handling to allow the server to shut down gracefully

	cancelSignal := make(chan os.Signal, 1)
	cleanupSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, os.Interrupt)
	signal.Notify(cleanupSignal, os.Interrupt)

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
