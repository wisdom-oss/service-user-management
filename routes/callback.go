package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccessCodeQuery struct {
	Code  string `form:"code"`
	State string `form:"state"`
}

func Callback(c *gin.Context) {
	var query AccessCodeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}
	tokenGrantUrl := fmt.Sprintf("/token?grant_type=authorization_code&code=%s&state=%s", query.Code, query.State)
	c.Redirect(http.StatusFound, tokenGrantUrl)
}
