package routes

import (
	"context"
	"fmt"
	"microservice/internal/db"
	"microservice/types"
	"path"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-contrib/location"
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

	url := location.Get(c)
	pathPrefix := c.Request.Header.Get("X-Forwarded-Prefix")

	c.JSON(200, gin.H{
		"issuer":                                "user-management",
		"authorization_endpoint":                fmt.Sprintf("%s://%s", url.Scheme, path.Clean(fmt.Sprintf("%s/%s/login", url.Host, pathPrefix))),
		"token_endpoint":                        fmt.Sprintf("%s://%s", url.Scheme, path.Clean(fmt.Sprintf("%s/%s/token", url.Host, pathPrefix))),
		"userinfo_endpoint":                     fmt.Sprintf("%s://%s", url.Scheme, path.Clean(fmt.Sprintf("%s/%s/users/me", url.Host, pathPrefix))),
		"jwks_uri":                              fmt.Sprintf("%s://%s", url.Scheme, path.Clean(fmt.Sprintf("%s/%s/.well-known/jwks.json", url.Host, pathPrefix))),
		"scopes_supported":                      scopes,
		"id_token_signing_alg_values_supported": []string{"none"},
		"response_types_supported":              []string{"token"},
		"grant_types_supported":                 []string{"authorization_code", "refresh_token", "client_credentials"},
	})
}
