package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Callback(c *gin.Context) {
	var query struct {
		Code  string `form:"code"`
		State string `form:"state"`
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}
	c.String(http.StatusSeeOther, `Please send a POST request to '/token?grant_type=authorization_code&code=%s&state=%s' to generate a token set`, query.Code, query.State)
}
