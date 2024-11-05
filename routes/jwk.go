package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"microservice/resources"
)

func JWK(c *gin.Context) {
	set := jwk.NewSet()
	set.AddKey(resources.PublicJWK)
	c.JSON(http.StatusOK, set)
}
