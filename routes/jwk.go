package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"microservice/resources"
)

func JWK(c *gin.Context) {
	c.JSON(http.StatusOK, resources.KeySet)
}
