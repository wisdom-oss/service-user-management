package permissions

import (
	"github.com/gin-gonic/gin"
	commonTypes "github.com/wisdom-oss/common-go/v2/types"

	"microservice/internal/db"
	"microservice/internal/errors"
	"microservice/types"
	"microservice/utils"
)

// Assign only creates new permission assignments and doesn't modify existing
// ones
func Assign(c *gin.Context) {
	var parameters struct {
		UserID      string `json:"user" binding:"required"`
		Assignments []struct {
			Service string `json:"service" binding:"required"`
			Scope   string `json:"scope" binding:"required"`
		} `json:"assignments" binding:"required"`
	}
	err := c.BindJSON(&parameters)
	if err != nil {
		res := errors.ErrMissingParameter
		res.Errors = []error{err}
		res.Emit(c)
		return
	}

	user, err := utils.GetUser(types.InternalIdentifier(parameters.UserID))
	if err != nil {
		c.Abort()
		if err == utils.ErrNoUser {
			errors.ErrUnknownUser.Emit(c)
		} else {
			_ = c.Error(err)
		}
		return
	}

	query, err := db.Queries.Raw("assign-permission")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}
	tx, err := db.Pool.Begin(c)

	for _, assignment := range parameters.Assignments {
		var service types.Service
		err = service.LoadFromDB(types.ExternalIdentifier(assignment.Service))
		if err != nil {
			c.Abort()
			tx.Rollback(c)
			_ = c.Error(err)
			return
		}

		var scope commonTypes.Scope
		err = scope.Parse(assignment.Scope)
		if err != nil {
			c.Abort()
			tx.Rollback(c)
			errors.ErrInvalidScope.Emit(c)
			return
		}

		_, err = tx.Exec(c, query, user.ID, service.ID, scope)
		if err != nil {
			c.Abort()
			tx.Rollback(c)
			_ = c.Error(err)
			return
		}
	}
	response := struct {
		types.User
		Permissions map[string][]string `json:"permissions"`
	}{
		User:        *user,
		Permissions: user.Permissions(),
	}
	c.JSON(200, response)
}
