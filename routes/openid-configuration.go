package routes

import (
	"context"
	"fmt"
	"microservice/internal/db"
	"microservice/types"
	"path"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
)

func OpenIDConfiguration(c *gin.Context) {
	query, err := db.Queries.Raw("get-services")
	if err != nil {
		panic(err)
	}
	var services []types.Service
	err = pgxscan.Select(context.Background(), db.Pool, &services, query)
	if err != nil {
		panic(err)
	}

	scopes := make([]string, 0)
	for _, service := range services {
		for _, scope := range service.SupportedScopes {
			scopes = append(scopes, fmt.Sprintf("%s:%s", service.Name, scope))
		}
	}

	scopes = append(scopes, "*:*")

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if proxiedScheme := c.Request.Header.Get("X-Forwarded-Proto"); proxiedScheme != "" {
		scheme = proxiedScheme
	}

	pathPrefix := c.Request.Header.Get("X-Forwarded-Prefix")

	c.JSON(200, gin.H{
		"issuer":                                "user-management",
		"authorization_endpoint":                fmt.Sprintf("%s://%s", scheme, path.Clean(fmt.Sprintf("%s/%s/login", c.Request.Host, pathPrefix))),
		"userinfo_endpoint":                     fmt.Sprintf("%s://%s", scheme, path.Clean(fmt.Sprintf("%s/%s/users/me", c.Request.Host, pathPrefix))),
		"token_endpoint":                        fmt.Sprintf("%s://%s", scheme, path.Clean(fmt.Sprintf("%s/%s/token", c.Request.Host, pathPrefix))),
		"jwks_uri":                              fmt.Sprintf("%s://%s", scheme, path.Clean(fmt.Sprintf("%s/%s/.well-known/jwks.json", c.Request.Host, pathPrefix))),
		"scopes_supported":                      scopes,
		"id_token_signing_alg_values_supported": []string{"none"},
		"response_types_supported":              []string{"token"},
		"grant_types_supported":                 []string{"authorization_code", "refresh_token", "client_credentials"},
	})
}
