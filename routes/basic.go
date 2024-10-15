package routes

import (
	"github.com/gin-gonic/gin"
)

// BasicHandler contains just a response, that is used to show the templating
func BasicHandler(c *gin.Context) {
	c.String(200, "hello there")
}
