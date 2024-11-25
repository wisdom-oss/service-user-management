package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func OpenIDConfiguration(c *gin.Context) {
	c.JSON(200, gin.H{
		"jwks_uri": fmt.Sprintf("%s/.well-known/jwks.json", c.Request.Host),
	})
}
