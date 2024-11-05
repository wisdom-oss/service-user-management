package routes

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/wisdom-oss/common-go/v2/middleware"
	"github.com/wisdom-oss/common-go/v2/types"

	"microservice/utils"
)

func UserInformation(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "me" {
		_userID, set := c.Get("subject")
		if !set {
			c.Abort()
			_ = c.Error(errors.New("no subject found in request context"))
			return
		}
		userID, _ = _userID.(string)
	} else {
		// let the request pass through the scope requirer, to protect from reading
		// the user properties of other users
		handler := middleware.RequireScope{}.Gin("user-management", types.ScopeRead)
		handler(c)
	}

	user, err := utils.GetUser(utils.InternalID(userID))
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.JSON(200, user)

}
