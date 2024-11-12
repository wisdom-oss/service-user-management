package users

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/wisdom-oss/common-go/v2/middleware"
	commonTypes "github.com/wisdom-oss/common-go/v2/types"

	"microservice/types"
	"microservice/utils"
)

func Information(c *gin.Context) {
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
		handler := middleware.RequireScope{}.Gin("user-management", commonTypes.ScopeRead)
		handler(c)
	}

	user, err := utils.GetUser(types.InternalIdentifier(userID))
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.JSON(200, types.ExtendedUser{
		User:        *user,
		Permissions: user.Permissions(),
	})

}
