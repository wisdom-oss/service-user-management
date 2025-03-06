package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"microservice/internal/db"
	"microservice/internal/errors"
)

func Delete(c *gin.Context) {
	userID := c.Param("userID")
	err := uuid.Validate(userID)
	if err != nil {
		c.Abort()
		errors.ErrUnknownUser.Emit(c)
		return
	}

	query, err := db.Queries.Raw("delete-user")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	_, err = db.Pool.Exec(c, query, userID)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
